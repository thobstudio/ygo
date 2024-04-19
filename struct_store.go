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
	return store.clients[client]
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
		return lastStruct.Id().clock + lastStruct.Length()
	}
	return 0
}

func (store *StructStore) StateVector() StateVector {
	sm := make(StateVector, len(store.clients))
	for client, structs := range store.clients {
		structItem := structs[len(structs)-1]
		sm[client] = structItem.Id().clock + structItem.Length()
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
		if lastStruct.Id().clock+lastStruct.Length() != clock {
			return errors.New("AddStructItem unexpected case")
		}
	}
	structs = append(structs, item)
	store.clients[client] = structs
	return nil
}

func findIndexSS(structs []SharedStruct, clock uint32) (uint32, error) {
	left := uint32(0)
	right := uint32(len(structs) - 1)
	mid := structs[right]
	if mid.Id().clock == clock {
		return right, nil
	}

	midindex := uint32(math.Floor(float64((clock / (mid.Id().clock + mid.Length() - 1)) * right)))
	for left <= right {
		mid = structs[midindex]
		if mid.Id().clock <= clock {
			if clock < mid.Id().clock+mid.Length() {
				return midindex, nil
			}
			left = midindex + 1
		} else {
			right = midindex - 1
		}
		midindex = uint32(math.Floor(float64((left + right) / 2)))
	}

	return 0, errors.New("findIndexSS unexpected case")
}

func (store *StructStore) FindStructIndex(client uint32, clock uint32) (uint32, error) {
	if structs, ok := store.clients[client]; ok {
		return findIndexSS(structs, clock)
	}
	return 0, errors.New(fmt.Sprintf("FindStructIndex invalid %v client", client))
}

func (store *StructStore) GetItem(id *ID) (SharedStruct, error) {
	if structs, ok := store.clients[id.client]; ok {
		index, err := findIndexSS(structs, id.clock)
		if err != nil {
			return nil, err
		}
		return structs[index], nil
	}
	return nil, errors.New(fmt.Sprintf("GetItem no items for client : %v \n", id.client))
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
		return nil, errors.New("StructStore client not found")
	}

	index, err := findIndexSS(structs, id.clock)
	if err != nil {
		return nil, err
	}

	str := structs[index]

	if id.clock != str.Id().clock+str.Length()-1 && reflect.TypeOf(str) != reflect.TypeOf(&Gc{}) {
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

func (store *StructStore) WriteStructs(encoder *UpdateEncoderV1, client uint32, clock uint32) error {
	structs := store.GetStructs(client)
	clock = max(clock, structs[0].Id().clock)
	startNewStuct, err := findIndexSS(structs, clock)
	if err != nil {
		return nil
	}

	if err = lib0.WriteVarUint(encoder.writter, uint32(len(structs)-int(startNewStuct))); err != nil {
		return err
	}
	if err = encoder.WriteClient(client); err != nil {
		return err
	}

	if err = lib0.WriteVarUint(encoder.writter, clock); err != nil {
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

func (store *StructStore) WriteClientStructs(encoder *UpdateEncoderV1, structs map[uint32]uint32) error {
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

	if err := lib0.WriteVarUint(encoder.writter, uint32(len(filteredStructs))); err != nil {
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
