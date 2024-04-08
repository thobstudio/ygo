package ynotgo

type SharedType interface {
	Integrate(doc *Doc, item *Item)
	Copy() SharedType
	Clone() SharedType
}
