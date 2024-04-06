package ynotgo

type YMap struct {
	doc           *Doc
	length        uint32
	item          *Item
	start         *Item
	itemMap       map[string]*Item
	prelimContent map[string]interface{}
}

func newYMap() *YMap {
	return &YMap{}
}

func NewYMap() *YMap {
	return newYMap()
}

func (m *YMap) integrate(doc *Doc, item *Item) {
	m.doc = doc
	m.item = item
}
