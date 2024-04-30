package ynotgo

import (
	"cmp"
	"errors"
	"math"
	"slices"

	"github.com/thobstudio/ynotgo/lib0"
)

type DeleteItem struct {
	clock  uint32
	length uint32
}

func newDeleteItem(clock uint32, length uint32) *DeleteItem {
	return &DeleteItem{
		clock:  clock,
		length: length,
	}
}

type DeleteSet struct {
	clients map[uint32][]*DeleteItem
}

func newDeleteSet() *DeleteSet {
	return &DeleteSet{
		clients: make(map[uint32][]*DeleteItem, 0),
	}
}

func NewDeleteSet() *DeleteSet {
	return newDeleteSet()
}

func (ds *DeleteSet) ClientsCount() int {
	return len(ds.clients)
}

func (ds *DeleteSet) SetClientDeleteItems(client uint32, items []*DeleteItem) {
	ds.clients[client] = items
}

func (ds *DeleteSet) AddDeleteItem(client, clock, length uint32) {
	if _, ok := ds.clients[client]; !ok {
		ds.clients[client] = make([]*DeleteItem, 0)
	}

	ds.clients[client] = append(ds.clients[client], newDeleteItem(clock, length))
}

func (ds *DeleteSet) SortAndMergeDeleteSet() {
	for client, items := range ds.clients {
		slices.SortFunc(items, func(a *DeleteItem, b *DeleteItem) int {
			return cmp.Compare(a.clock, b.clock)
		})

		i, j := 1, 1
		for ; i < len(items); i++ {
			left := items[j-1]
			right := items[i]
			if left.clock+left.length >= right.clock {
				left.length = max(left.length, right.clock+right.length-left.clock)
			} else {
				if j < i {
					items[j] = right
				}
				j = j + 1
			}
		}

		ds.clients[client] = items[:j]
	}
}

func findIndexDeleteSet(dis []*DeleteItem, clock uint32) (uint32, bool) {
	left := uint32(0)
	right := uint32(len(dis) - 1)
	for left <= right {
		midindex := uint32(math.Floor(float64((left + right) / 2)))
		mid := dis[midindex]
		if mid.clock <= clock {
			if clock < mid.clock+mid.length {
				return midindex, true
			}
			left = midindex + 1
		} else {
			right = midindex - 1
		}
	}
	return 0, false
}

func (ds *DeleteSet) TryGcDeleteSet(store *StructStore, gcFilter func(item *Item) bool) error {
	if gcFilter == nil {
		gcFilter = func(item *Item) bool { return true }
	}

	for client, deleteItems := range ds.clients {
		structs := store.GetStructs(client)

		for di := len(deleteItems) - 1; di > 0; di++ {
			deleteItem := deleteItems[di]
			endDeleteItemClock := deleteItem.clock + deleteItem.length
			si, err := store.FindStructIndex(client, deleteItem.clock)
			if err != nil {
				return err
			}
			for ; int(si) < len(structs); si++ {
				str := structs[si]
				if str.Id().clock >= endDeleteItemClock {
					break
				}

				if item, ok := str.(*Item); ok && item.Deleted() && !item.Keep() && gcFilter(item) {
					if err := item.Gc(store, false); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func (ds *DeleteSet) TryMergeDeleteSet(store *StructStore) error {
	for client, deleteItems := range ds.clients {
		for di := len(deleteItems) - 1; di >= 0; di-- {
			deleteItem := deleteItems[di]
			deleteItemIndexInStore, err := store.FindStructIndex(client, deleteItem.clock+deleteItem.length-1)
			if err != nil {
				return err
			}
			mostRightIndexToCheck := min(store.ClientsCount()-1, 1+int(deleteItemIndexInStore))
			si := mostRightIndexToCheck
			str := store.GetStructItem(client, si)
			for si > 0 && str.Id().clock >= deleteItem.clock {
				si -= 1 + store.MergeWithLefts(client, si)
				str = store.GetStructItem(client, si)
			}
		}
	}
	return nil
}

func (ds *DeleteSet) Write(encoder UpdateEncoder) error {
	if err := lib0.WriteVarUint(encoder.Writer(), uint32(ds.ClientsCount())); err != nil {
		return err
	}

	clients := make([]uint32, len(ds.clients))
	i := 0
	for client := range ds.clients {
		clients[i] = client
		i++
	}
	slices.Sort(clients)
	for _, client := range clients {
		dsItems := ds.clients[client]
		encoder.ResetDsCurVal()
		if err := lib0.WriteVarUint(encoder.Writer(), client); err != nil {
			return err
		}
		length := len(dsItems)
		if err := lib0.WriteVarUint(encoder.Writer(), uint32(length)); err != nil {
			return err
		}
		for i := 0; i < length; i++ {
			item := dsItems[i]
			if err := encoder.WriteDsClock(item.clock); err != nil {
				return err
			}
			if err := encoder.WriteDsLen(item.length); err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadDeleteSet(decoder DsDecoder) (*DeleteSet, error) {
	ds := newDeleteSet()
	reader := decoder.Reader()
	numClients, err := lib0.ReadVarUint(reader)
	if err != nil {
		return nil, err
	}
	for i := 0; i < int(numClients); i++ {
		decoder.ResetDsCurVal()
		client, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, errors.New("failed to read client")
		}
		numOfDeletes, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, errors.New("failed to read number of deletes")
		}
		if numOfDeletes > 0 {
			dsField, ok := ds.clients[client]
			if !ok {
				dsField = make([]*DeleteItem, 0)
			}
			for j := 0; j < int(numOfDeletes); j++ {
				dsClock, err := lib0.ReadVarUint(reader)
				if err != nil {
					return nil, errors.New("failed to delete set clock")
				}
				dsLength, err := lib0.ReadVarUint(reader)
				if err != nil {
					return nil, errors.New("failed to read delete set length")
				}
				// FIXME: instead of append we could probably grow the slice and then just assign value
				// that way it becomes more performant
				dsField = append(dsField, newDeleteItem(dsClock, dsLength))
			}
			ds.clients[client] = dsField
		}
	}

	return ds, nil
}
