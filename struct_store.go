package ynotgo

import (
	"errors"
)

type StateVector map[uint32]uint32

type PendingStructs struct {
	missing map[uint32]uint32
	update  []byte
}

type StructStore struct {
	clients        map[uint32][]interface{}
	pendingStructs *PendingStructs
	pendingDs      []byte
}

func newStructStore() *StructStore {
	return &StructStore{
		clients: make(map[uint32][]interface{}, 0),
	}
}

func getIdAndLength(s interface{}) (*ID, uint32, error) {
	switch t := s.(type) {
	case *Item:
		return t.id, t.length, nil
	case *Gc:
		return t.id, t.length, nil
	default:
		return nil, 0, errors.New("unsupported type")
	}
}

func (store *StructStore) State(client uint32) uint32 {
	if structs, ok := store.clients[client]; ok {
		sid, slen, err := getIdAndLength(structs[len(structs)-1])
		if err != nil {
			return 0
		}
		return sid.clock + slen
	}
	return 0
}

func (store *StructStore) StateVector() StateVector {
	sm := make(StateVector, 0)
	for client, items := range store.clients {
		s := items[len(items)-1]
		sId, sLen, err := getIdAndLength(s)
		if err != nil {
			sm[client] = sId.clock + sLen
		}
	}

	return sm
}
func (store *StructStore) AddItemOrGc(item any) error {
	sid, _, err := getIdAndLength(item)
	if err != nil {
		return err
	}
	structs, ok := store.clients[sid.client]
	if !ok {
		store.clients[sid.client] = make([]any, 0)
	} else {
		lsid, lslen, err := getIdAndLength(structs[len(structs)-1])
		if err != nil {
			return err
		}

		if lsid.clock+lslen != sid.clock {
			return errors.New("unexpected case")
		}
	}

	store.clients[sid.client] = append(store.clients[sid.client], item)

	return nil
}
