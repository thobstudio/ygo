package ynotgo

type ID struct {
	client uint32
	clock  uint32
}

func newId(client uint32, clock uint32) *ID {
	return &ID{
		client: client,
		clock:  clock,
	}
}

func NewID(client uint32, clock uint32) *ID {
	return newId(client, clock)
}
