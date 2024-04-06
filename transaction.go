package ynotgo

type Transaction struct {
	doc            *Doc
	origin         interface{}
	local          bool
	deleteSet      *DeleteSet
	subdocsAdded   map[*Doc]bool
	subdocsRemoved map[*Doc]bool
	subdocsLoaded  map[*Doc]bool
	beforeState    map[uint32]uint32
	afterState     map[uint32]uint32
}

func newTrasaction(doc *Doc, origin interface{}, local bool) *Transaction {
	return &Transaction{
		doc:         doc,
		deleteSet:   newDeleteSet(),
		beforeState: doc.store.StateVector(),
		afterState:  make(map[uint32]uint32, 0),
		//
		origin:         origin,
		local:          local,
		subdocsAdded:   make(map[*Doc]bool, 0),
		subdocsRemoved: make(map[*Doc]bool, 0),
		subdocsLoaded:  make(map[*Doc]bool, 0),
	}
}

func NewTrasaction(doc *Doc, origin interface{}, local bool) *Transaction {
	return newTrasaction(doc, origin, local)
}

type TransactionHandler func(tx *Transaction) (interface{}, error)

