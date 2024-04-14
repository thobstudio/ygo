package ynotgo

type Transaction struct {
	doc          *Doc
	origin       interface{}
	local        bool
	deleteSet    *DeleteSet
	changed      map[SharedType]map[string]bool // map[string]bool is like a set
	beforeState  map[uint32]uint32
	afterState   map[uint32]uint32
	mergeStructs []SharedStruct
}

func newTrasaction(doc *Doc, origin interface{}, local bool) *Transaction {
	return &Transaction{
		doc:          doc,
		deleteSet:    newDeleteSet(),
		beforeState:  doc.store.StateVector(),
		afterState:   make(map[uint32]uint32, 0),
		changed:      make(map[SharedType]map[string]bool),
		origin:       origin,
		local:        local,
		mergeStructs: make([]SharedStruct, 0),
	}
}

func NewTrasaction(doc *Doc, origin interface{}, local bool) *Transaction {
	return newTrasaction(doc, origin, local)
}

type TransactionHandler func(tx *Transaction) (interface{}, error)
func (tx *Transaction) AddChangedType(t SharedType, parentSub string) {
	item := t.Item()
	beforeStateClock, ok := tx.beforeState[item.id.client]
	if !ok {
		beforeStateClock = 0
	}
	if item == nil || ((item.id.clock < beforeStateClock) && !item.Deleted()) {
		if _, ok := tx.changed[t]; !ok {
			tx.changed[t] = make(map[string]bool)
		}
		tx.changed[t][parentSub] = true
	}
}
