package ynotgo

type Gc struct {
	id     *ID
	length uint32
}

func (item *Gc) Id() *ID {
	return item.id
}

func (item *Gc) Length() uint32 {
	return item.length
}

func (item *Gc) MergeWith() (bool, error) {
	return false, nil
}

func (item *Gc) Integrate(tx *Transaction, offset uint32) error {
	return nil
}
