package ynotgo

import "reflect"

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
	if reflect.TypeOf(item) != reflect.TypeOf(right) {
		return false, nil
	}
	item.length += right.Length()
	return false, nil
}

func (item *Gc) Integrate(tx *Transaction, offset uint32) error {
	if offset > 0 {
		item.id.clock += offset
		item.length -= offset
	}
	tx.doc.store.AddStructItem(item)
	return nil
}

func (item *Gc) Deleted() bool {
	return false
}

func (item *Gc) Write(encoder UpdateEncoder, offset uint32) error {
	return nil
}

func (item *Gc) GetMissing(tx *Transaction, store *StructStore) (uint32, bool) {
	return 0, false
}
