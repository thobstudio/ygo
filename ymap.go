package ynotgo

type YMap struct {
	AbstractType
}

func newYMap() *YMap {
	return &YMap{}
}

func NewYMap() *YMap {
	return newYMap()
}

func (m *YMap) integrate(doc *Doc) {
	m.AbstractType.integrate(doc)
}
