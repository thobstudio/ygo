package ynotgo

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
func (ds *DeleteSet) AddDeleteItem(client, clock, length uint32) {
	if _, ok := ds.clients[client]; !ok {
		ds.clients[client] = make([]*DeleteItem, 0)
	}

	ds.clients[client] = append(ds.clients[client], newDeleteItem(clock, length))
}

