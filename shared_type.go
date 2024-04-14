package ynotgo

type SharedType interface {
	Item() *Item
	Start() *Item
	SetStart(item *Item)
	Length() uint32
	SetLength(length uint32)
	GetItem(key string) *Item
	SetItem(key string, item *Item)
	Integrate(doc *Doc, item *Item)
	Copy() SharedType
	Clone() SharedType
}
