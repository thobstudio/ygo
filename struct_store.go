package ynotgo

import (
	"cmp"
	"errors"
	"fmt"
	"math"
	"reflect"
	"slices"

	"github.com/thobstudio/ynotgo/lib0"
)

type StateVector map[uint32]uint32

type PendingStructs struct {
	missing map[uint32]uint32
	update  []byte
}

type StructStore struct {
	clients        map[uint32][]SharedStruct
	pendingStructs *PendingStructs
	pendingDs      []byte
}

func newStructStore() *StructStore {
	return &StructStore{
		clients: make(map[uint32][]SharedStruct, 0),
	}
}

func (store *StructStore) GetStructItem(client uint32, index int) SharedStruct {
	if structs, ok := store.clients[client]; ok {
		return structs[index]
	}
	return nil
}

func (store *StructStore) GetStructs(client uint32) []SharedStruct {
	if structs, ok := store.clients[client]; ok {
		return structs
	}
	return nil
}

func (store *StructStore) SetStructs(client uint32, structs []SharedStruct) {
	store.clients[client] = structs
}

func (store *StructStore) ClientsCount() int {
	return len(store.clients)
}

func (store *StructStore) State(client uint32) uint32 {
	if structs, ok := store.clients[client]; ok {
		lastStruct := structs[len(structs)-1]
		return lastStruct.State()
	}
	return 0
}

func (store *StructStore) StateVector() StateVector {
	sm := make(StateVector, len(store.clients))
	for client, structs := range store.clients {
		structItem := structs[len(structs)-1]
		sm[client] = structItem.State()
	}

	return sm
}

func (store *StructStore) AddStructItem(item SharedStruct) error {
	client := item.Id().client
	clock := item.Id().clock
	structs, ok := store.clients[client]
	if !ok {
		structs = make([]SharedStruct, 0)
	} else {
		lastStruct := structs[len(structs)-1]
		if lastStruct.State() != clock {
			return errors.New("clock is not equal to state of the last struct item")
		}
	}
	store.clients[client] = append(structs, item)
	return nil
}

func findIndexSS(structs []SharedStruct, clock uint32) (uint32, error) {
	left := uint32(0)
	right := uint32(len(structs) - 1)
	mid := structs[right]
	if mid.Id().clock == clock {
		return right, nil
	}

	midindex := uint32(math.Floor(float64((clock / (mid.State() - 1)) * right)))
	for left <= right {
		mid = structs[midindex]
		if mid.Id().clock <= clock {
			if clock < mid.State() {
				return midindex, nil
			}
			left = midindex + 1
		} else {
			right = midindex - 1
		}
		midindex = uint32(math.Floor(float64((left + right) / 2)))
	}

	return 0, errors.New(fmt.Sprintf("struct item with clock %v not found \n", clock))
}

func (store *StructStore) FindStructIndex(client uint32, clock uint32) (uint32, error) {
	if structs, ok := store.clients[client]; ok {
		return findIndexSS(structs, clock)
	}
	return 0, errors.New(fmt.Sprintf("no structs for client %v \n", client))
}

func (store *StructStore) GetItem(id *ID) (SharedStruct, error) {
	if structs, ok := store.clients[id.client]; ok {
		index, err := findIndexSS(structs, id.clock)
		if err != nil {
			return nil, err
		}
		return structs[index], nil
	}
	return nil, errors.New(fmt.Sprintf("GetItem no structs for client : %v \n", id.client))
}

func (store *StructStore) ReplaceStruct(prev SharedStruct, next SharedStruct) error {
	client := prev.Id().client
	if structs, ok := store.clients[client]; ok {
		index, err := findIndexSS(structs, prev.Id().clock)
		if err != nil {
			return err
		}
		structs[index] = next
		store.clients[client] = structs
	}
	return nil
}

func (store *StructStore) GetItemCleanEnd(tx *Transaction, id *ID) (SharedStruct, error) {
	structs, ok := store.clients[id.client]
	if !ok {
		return nil, errors.New(fmt.Sprintf("not structs for client %v \n", id.client))
	}

	index, err := findIndexSS(structs, id.clock)
	if err != nil {
		return nil, err
	}

	str := structs[index]

	if id.clock != str.State()-1 && reflect.TypeOf(str) != reflect.TypeOf(&Gc{}) {
		store.InsertStruct(id.client, index+1, str.(*Item).SplitItem(tx, id.clock-str.Id().clock+1))
	}

	return str, nil
}

func (store *StructStore) InsertStruct(client uint32, index uint32, structItem SharedStruct) {
	if structs, ok := store.clients[client]; ok {
		structs = append(structs[:index+1], structs[index:]...)
		structs[index] = structItem
		store.clients[client] = structs
	}
}

func (store *StructStore) MergeWithLefts(client uint32, pos int) int {
	if structs, ok := store.clients[client]; ok {
		right := structs[pos]
		left := structs[pos-1]
		i := pos
		for i > 0 {
			if left.Deleted() == right.Deleted() && reflect.TypeOf(left) == reflect.TypeOf(right) {
				if ok, err := left.MergeWith(right); ok && err != nil {
					if rightItem, ok := right.(*Item); ok && rightItem.parentSub != "" {
						if parent, ok := rightItem.parent.(SharedType); ok {
							if item := parent.GetItem(rightItem.parentSub); item == rightItem {
								parent.SetItem(rightItem.parentSub, left.(*Item))
							}
						}
					}

					i -= 1
					right = left
					left = structs[i]
					continue
				}
			}
			break
		}

		merged := pos - i
		if merged > 0 {
			structs = structs[pos+1-merged : merged]
		}
		store.clients[client] = structs
	}
	return 0
}

func (store *StructStore) WriteStructs(encoder UpdateEncoder, client uint32, clock uint32) error {
	structs := store.GetStructs(client)
	clock = max(clock, structs[0].Id().clock)
	startNewStuct, err := findIndexSS(structs, clock)
	if err != nil {
		return nil
	}

	if err = lib0.WriteVarUint(encoder.Writer(), uint32(len(structs)-int(startNewStuct))); err != nil {
		return err
	}
	if err = encoder.WriteClient(client); err != nil {
		return err
	}

	if err = lib0.WriteVarUint(encoder.Writer(), uint32(clock)); err != nil {
		return err
	}

	firstStruct := structs[startNewStuct]
	if err = firstStruct.Write(encoder, clock-firstStruct.Id().clock); err != nil {
		return err
	}

	for i := startNewStuct + 1; i < uint32(len(structs)); i++ {
		err = structs[i].Write(encoder, 0)
		if err != nil {
			return err
		}
	}

	return nil
}

func (store *StructStore) WriteClientStructs(encoder UpdateEncoder, structs map[uint32]uint32) error {
	filteredStructs := make(map[uint32]uint32, len(structs))
	for client, clock := range structs {
		if store.State(client) > clock {
			filteredStructs[client] = clock
		}
	}

	for client := range store.StateVector() {
		if _, ok := structs[client]; !ok {
			filteredStructs[client] = 0
		}
	}

	if err := lib0.WriteVarUint(encoder.Writer(), uint32(len(filteredStructs))); err != nil {
		return err
	}

	ss := make([]uint32, len(filteredStructs))
	index := 0
	for client := range filteredStructs {
		ss[index] = client
		index++
	}

	slices.SortFunc(ss, func(a, b uint32) int {
		return cmp.Compare(b, a)
	})

	for _, client := range ss {
		if err := store.WriteStructs(encoder, client, filteredStructs[client]); err != nil {
			return err
		}
	}
	return nil
}

func (store *StructStore) WriteDeleteSet(encoder UpdateEncoder) error {
	ds, err := store.CreateDeleteSet()
	if err != nil {
		return err
	}

	if err := lib0.WriteVarUint(encoder.Writer(), uint32(ds.ClientsCount())); err != nil {
		return err
	}
	ds.ForEach(func(client uint32, dsItems []*DeleteItem) {
		fmt.Printf("client : %v dsitems : %v ", client, dsItems)
		encoder.ResetDsCurVal()
		lib0.WriteVarUint(encoder.Writer(), client)
		length := len(dsItems)
		lib0.WriteVarUint(encoder.Writer(), uint32(length))
		for i := 0; i < length; i++ {
			item := dsItems[i]
			encoder.WriteDsClock(item.clock)
			encoder.WriteDsLen(item.length)
		}
	})

	return nil
}

func (store *StructStore) CreateDeleteSet() (*DeleteSet, error) {
	ds := newDeleteSet()
	for client, structs := range store.clients {
		dsItems := make([]*DeleteItem, 0)
		for i := 0; i < len(structs); i++ {
			_struct := structs[i]
			if _struct.Deleted() {
				clock := _struct.Id().clock
				length := _struct.Length()
				if i+1 < len(structs) {
					for j := i + 1; j < len(structs); j++ {
						next := structs[j]
						length += next.Length()
					}
				}
				dsItems = append(dsItems, newDeleteItem(clock, length))
			}
		}

		if len(dsItems) > 0 {
			ds.clients[client] = dsItems
		}
	}
	return ds, nil
}

type ClientStructRef struct {
	i    uint32
	refs []SharedStruct
}

func (store *StructStore) integrateStructs(clientsStructsRefs map[uint32]*ClientStructRef, tx *Transaction) (*PendingStructs, error) {
	if len(clientsStructsRefs) == 0 {
		return nil, nil
	}
	clientRefIds := make([]uint32, len(clientsStructsRefs))
	i := 0
	for client := range clientsStructsRefs {
		clientRefIds[i] = client
		i++
	}
	// Sort descending
	slices.SortFunc(clientRefIds, func(a, b uint32) int {
		return cmp.Compare(b, a)
	})

	stack := make([]SharedStruct, 0)

	getNextStructTarget := func() *ClientStructRef {
		if len(clientRefIds) == 0 {
			return nil
		}
		nextStructsTarget := clientsStructsRefs[clientRefIds[len(clientRefIds)-1]]
		for uint32(len(nextStructsTarget.refs)) == nextStructsTarget.i {
			// pop id
			clientRefIds = clientRefIds[1:]
			if len(clientRefIds) > 0 {
				nextStructsTarget = clientsStructsRefs[clientRefIds[len(clientRefIds)-1]]
			} else {
				return nil
			}
		}

		return nextStructsTarget
	}

	curStructsTarget := getNextStructTarget()
	if curStructsTarget == nil {
		return nil, nil
	}

	restStructs := newStructStore()
	missingSv := make(map[uint32]uint32)
	//
	updateMissingSv := func(client, clock uint32) {
		mclock, ok := missingSv[client]
		if !ok || mclock > clock {
			missingSv[client] = clock
		}
	}

	addStackToRestSS := func() {
		for _, item := range stack {
			client := item.Id().client
			if unapplicableItems, ok := clientsStructsRefs[client]; ok {
				unapplicableItems.i--
				restStructs.clients[client] = unapplicableItems.refs[unapplicableItems.i:]
				delete(clientsStructsRefs, client)
				unapplicableItems.i = 0
				unapplicableItems.refs = []SharedStruct{}
			} else {
				restStructs.SetStructs(client, []SharedStruct{item})
			}

			filteredClientRefIds := make([]uint32, len(clientRefIds))
			i := 0
			for _, id := range clientRefIds {
				if id != client {
					filteredClientRefIds[i] = client
					i++
				}
			}
			clientRefIds = filteredClientRefIds[:i]
		}
		stack = []SharedStruct{}
	}

	stackHead := curStructsTarget.refs[curStructsTarget.i]
	curStructsTarget.i++
	state := make(map[uint32]uint32)

	for {
		if reflect.TypeOf(stackHead) != reflect.TypeOf(&Skip{}) {
			shclient := stackHead.Id().client
			localClock, ok := state[shclient]
			if !ok {
				localClock = store.State(shclient)
				state[shclient] = localClock
			}
			offset := localClock - stackHead.Id().clock
			if offset < 0 {
				stack = append(stack, stackHead)
				updateMissingSv(stackHead.Id().client, stackHead.Id().clock-1)
				addStackToRestSS()
				//
			} else {
				if missing, ok := stackHead.GetMissing(tx, store); ok {
					stack = append(stack, stackHead)

					structRefs, ok := clientsStructsRefs[missing]
					if ok {
						updateMissingSv(missing, store.State(missing))
						addStackToRestSS()
					} else {
						stackHead = structRefs.refs[structRefs.i]
						structRefs.i++
						continue
					}
					//
				} else if offset == 0 || offset < stackHead.Length() {
					stackHead.Integrate(tx, offset)
					// Cache the stack head
					state[stackHead.Id().client] = stackHead.State()
				}
				//
			}
		}

		// Iterate to next stack head
		if len(stack) > 0 {
			stackHead = stack[0]
			stack = stack[1:]
		} else if curStructsTarget != nil && curStructsTarget.i < uint32(len(curStructsTarget.refs)) {
			stackHead = curStructsTarget.refs[curStructsTarget.i]
			curStructsTarget.i++
		} else {
			curStructsTarget = getNextStructTarget()
			if curStructsTarget == nil {
				break
			} else {
				stackHead = curStructsTarget.refs[curStructsTarget.i]
				curStructsTarget.i++
			}
		}
	}

	if restStructs.ClientsCount() > 0 {
		encoder := newUpdateEncoderV1()
		err := restStructs.WriteClientStructs(encoder, make(map[uint32]uint32))
		if err != nil {
			return nil, err
		}
		update, err := encoder.ToUint8Array()
		if err != nil {
			return nil, err
		}
		return &PendingStructs{
			missing: missingSv,
			update:  update,
		}, nil
	}

	return nil, nil
}

func (store *StructStore) readAndApplyDeleteSet(decoder UpdateDecoder, tx *Transaction) ([]byte, error) {
	unappliedDs := newDeleteSet()
	reader := decoder.Reader()
	numClients, err := lib0.ReadVarUint(reader)
	if err != nil {
		return nil, err
	}
	//
	for i := 0; i < int(numClients); i++ {
		decoder.ResetDsCurVal()
		client, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, errors.New("failed to read client from decoder")
		}
		numOfDeletes, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, errors.New("failed to read number of deletes from decoder")
		}
		structs := store.GetStructs(client)
		if structs == nil {
			structs = make([]SharedStruct, 0)
		}
		state := store.State(client)

		for i := 0; i < int(numOfDeletes); i++ {
			clock, err := decoder.ReadDsClock()
			if err != nil {
				return nil, err
			}
			dsLength, err := decoder.ReadDsLen()
			if err != nil {
				return nil, err
			}
			clockEnd := clock + dsLength

			if clock < state {
				if state < clockEnd {
					unappliedDs.AddDeleteItem(client, state, clockEnd-state)
				}

				index, err := findIndexSS(structs, clock)
				if err != nil {
					return nil, err
				}

				if structItem, ok := structs[index].(*Item); ok && !structItem.Deleted() && structItem.id.clock < clock {
					structs = append(structs[:index+1], structs[index:]...)
					structs[index] = structItem.SplitItem(tx, clock-structItem.id.clock)
					index++
				}

				for int(index) < len(structs) {
					if structItem, ok := structs[index].(*Item); ok {
						index++
						if structItem.id.clock < clockEnd {
							if !structItem.Deleted() {
								if clockEnd < structItem.State() {
									structs = append(structs[:index+1], structs[index:]...)
									structs[index] = structItem.SplitItem(tx, clock-structItem.id.clock)
								}
								structItem.Delete(tx)
							}
						} else {
							break
						}
					}
				}
			} else {
				unappliedDs.AddDeleteItem(client, clockEnd, clockEnd-clock)
			}
		}

		if unappliedDs.ClientsCount() > 0 {
			dsencoder := newUpdateEncoderV2()
			if err := lib0.WriteVarUint(dsencoder.writer, 0); err != nil {
				return nil, errors.New("failed to write structs length")
			}
			if err := unappliedDs.Write(dsencoder); err != nil {
				return nil, err
			}
			update, err := dsencoder.ToUint8Array()
			if err != nil {
				return nil, errors.New("failed to encode unapplied delete set")
			}

			return update, nil
		}
	}

	return nil, nil
}
