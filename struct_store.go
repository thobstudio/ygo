package ynotgo

import (
	"errors"
	"fmt"
	"math"
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
	sm := make(StateVector, 0)
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

	midindex := uint32(math.Floor(float64((clock / (mid.Id().clock + mid.Length() - 1) / 2))))
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

