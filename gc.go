package ynotgo

type Gc struct {
	id     *ID
	length uint32
}

func newGc(id *ID, length uint32) *Gc {
	return &Gc{
		id:     id,
		length: length,
	}
}

func (item *Gc) Id() *ID {
	return item.id
}

func (item *Gc) Length() uint32 {
	return item.length
}

func (item *Gc) MergeWith(right SharedStruct) (bool, error) {
	return false, nil
}

func (item *Gc) Integrate(tx *Transaction, offset uint32) error {
	return nil
}

func (item *Gc) Deleted() bool {
	return false
}
