package ynotgo

type AbstractType struct {
	doc     *Doc
	length  uint32
	item    *Item
	start   *Item
	itemMap map[string]*Item
}

func (at *AbstractType) integrate(doc *Doc, item *Item) {
	at.doc = doc
	at.item = item
}
