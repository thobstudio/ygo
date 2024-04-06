package ynotgo

type PendingStructs struct {
	missing map[uint32]uint32
	update  []byte
}

type StructStore struct {
	clients        map[uint32]interface{}
	pendingStructs *PendingStructs
	pendingDs      []byte
}

func newStructStore() *StructStore {
	return &StructStore{
		clients: make(map[uint32]interface{}, 0),
	}
}
