package ynotgo

type Content interface {
	Content() []any
	Countable() bool
	Length() uint32
	Delete(tx *Transaction)
}
