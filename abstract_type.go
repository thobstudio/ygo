package ynotgo

type AbstractType struct {
	doc *Doc
}

func (at *AbstractType) integrate(doc *Doc) {
	at.doc = doc
}
