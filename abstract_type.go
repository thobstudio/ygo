package ynotgo

type AbstractType struct {
	item    *Item
	start   *Item
	doc     *Doc
	length  uint32
	itemMap map[string]*Item
}

func (at *AbstractType) Doc() *Doc {
	return at.doc
}

func (at *AbstractType) Item() *Item {
	return at.item
}

func (at *AbstractType) Start() *Item {
	return at.start
}

func (at *AbstractType) SetStart(item *Item) {
	at.start = item
}

func (at *AbstractType) Length() uint32 {
	return at.length
}

func (at *AbstractType) SetLength(length uint32) {
	at.length = length
}

func (at *AbstractType) GetItem(key string) *Item {
	return at.itemMap[key]
}

func (at *AbstractType) SetItem(key string, item *Item) {
	at.itemMap[key] = item
}

func (at *AbstractType) Integrate(doc *Doc, item *Item) {
	at.doc = doc
	at.item = item
}

func (at *AbstractType) Copy() SharedType {
	return &AbstractType{}
}

func (at *AbstractType) Clone() SharedType {
	return &AbstractType{}
}

func (at *AbstractType) ToJSON() any {
	return nil
}

func (at *AbstractType) Write(encoder UpdateEncoder) error {
	panic("to be implemented by parent structs")
}
