package ynotgo

type YMap struct {
	AbstractType
	prelimContent map[string]interface{}
}

func newYMap() *YMap {
	return &YMap{}
}

func NewYMap() *YMap {
	return newYMap()
}

func (m *YMap) Integrate(doc *Doc, item *Item) {
	m.doc = doc
	m.item = item
}

func (m *YMap) Has(key string) bool {
	if val, ok := m.itemMap[key]; ok {
		if val != nil {
			return !val.Deleted()
		}
	}
	return false
}

func (m *YMap) Get(key string) interface{} {
	if val, ok := m.itemMap[key]; ok {
		if val != nil && !val.Deleted() {
			return val.content.Content()[val.length-1]
		}
	}
	return nil
}
