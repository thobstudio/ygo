package ynotgo

import "errors"

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

func compareIds(a *ID, b *ID) bool {
	return a == b || a != nil && b != nil && a.client == b.client && a.clock == b.clock
}

func findRootTypeKey(t SharedType) (string, error) {
	for key, value := range t.Doc().share {
		if value == t {
			return key, nil
		}
	}
	return "", errors.New("unexpected case")
}
