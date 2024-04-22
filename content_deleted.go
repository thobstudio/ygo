package ynotgo

type ContentDeleted struct {
	length uint32
}

func newContentDeleted(length uint32) *ContentDeleted {
	return &ContentDeleted{
		length: length,
	}
}

func (content *ContentDeleted) Length() int {
	return int(content.length)
}

func (content *ContentDeleted) Content() []any {
	return []any{}
}

func (content *ContentDeleted) Countable() bool {
	return false
}

func (content *ContentDeleted) Copy() Content {
	return &ContentDeleted{
		length: content.length,
	}
}

func (content *ContentDeleted) Splice(offset uint32) (Content, error) {
	right := &ContentDeleted{
		length: content.length - offset,
	}
	content.length = offset
	return right, nil
}

func (content *ContentDeleted) MergeWith(right Content) (bool, error) {
	content.length += uint32(right.Length())
	return true, nil
}

func (content *ContentDeleted) Integrate(tx *Transaction, item *Item) {
	tx.deleteSet.AddDeleteItem(item.id.client, item.id.clock, content.length)
	item.markDeleted()
}

func (content *ContentDeleted) Delete(tx *Transaction) {}

func (content *ContentDeleted) Gc(store *StructStore) {}

func (content *ContentDeleted) Write(encoder *UpdateEncoderV1, offset uint32) error {
	return encoder.WriteLen(content.length - offset)
}

func (content *ContentDeleted) Ref() uint8 {
	return 1
}
