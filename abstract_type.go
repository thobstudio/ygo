package ynotgo

type AbstractType struct {
	item    *Item
	start   *Item
	doc     *Doc
	length  uint32
	itemMap map[string]*Item
}

func (at *AbstractType) Integrate(doc *Doc, item *Item) {
	at.doc = doc
	at.item = item
}
