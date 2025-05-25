package ygo

type ID struct {
	client uint
	clock  uint
}

func NewID(client, clock uint) *ID {
	return &ID{
		client: client,
		clock:  clock,
	}
}
