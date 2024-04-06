package ynotgo

type DeleteItem struct {
	clock  uint32
	length uint32
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
