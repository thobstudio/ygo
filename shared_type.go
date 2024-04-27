package ynotgo

type SharedType interface {
	Doc() *Doc
	Item() *Item
	Start() *Item
	SetStart(item *Item)
	Length() uint32
	SetLength(length uint32)
	GetItem(key string) *Item
	SetItem(key string, item *Item)
	ForEachItem(func(item *Item, key string))
	ClearItemMap()
	Integrate(doc *Doc, item *Item)
	Copy() SharedType
	Clone() SharedType
	ToJSON() any
	Write(encoder UpdateEncoder) error
}
