package ynotgo

type Content interface {
	Content() []any
	Countable() bool
	Length() int
	Delete(tx *Transaction)
	integrate(tx *Transaction, item *Item) error
}
