package ynotgo

import (
	"cmp"
	"math"
	"slices"
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

func (ds *DeleteSet) AddDeleteItem(client, clock, length uint32) {
	if _, ok := ds.clients[client]; !ok {
		ds.clients[client] = make([]*DeleteItem, 0)
	}

	ds.clients[client] = append(ds.clients[client], newDeleteItem(clock, length))
}

func (ds *DeleteSet) SortAndMerge() {
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
