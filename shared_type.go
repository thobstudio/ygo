package ynotgo

type SharedType interface {
	Item() *Item
	Integrate(doc *Doc, item *Item)
	Copy() SharedType
	Clone() SharedType
}
