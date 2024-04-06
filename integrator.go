package ynotgo

type DocIntegrator interface {
	integrate(doc *Doc, item *Item)
}
