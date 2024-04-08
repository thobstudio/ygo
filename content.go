package ynotgo

type Content interface {
	Length() int
	Content() []any
	Countable() bool
	Copy() Content
	Splice(offset uint32) (Content, error)
	MergeWith(right Content) (bool, error)
	Integrate(tx *Transaction, item *Item)
	Delete(tx *Transaction)
	Gc(store *StructStore)
	Write()
	Ref() uint8
}
