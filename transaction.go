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
	var beforeStateClock uint32 = 0
	if item != nil {
		if c, ok := tx.beforeState[item.id.client]; ok {
			beforeStateClock = c
		}
	}
	if item == nil || ((item.id.clock < beforeStateClock) && !item.Deleted()) {
		if _, ok := tx.changed[t]; !ok {
			tx.changed[t] = make(map[string]bool)
		}
		tx.changed[t][parentSub] = true
	}
}

func cleanupTransactions(transactionCleanups []*Transaction, index int) error {
	var err error
	if index < len(transactionCleanups) {
		tx := transactionCleanups[index]
		doc := tx.doc
		store := doc.store
		ds := tx.deleteSet

		ds.SortAndMergeDeleteSet()
		tx.afterState = store.StateVector()

		actions := make([]func(), 0)
		actions = append(actions, func() {
			//
		})

		// Run all the actions
		for _, action := range actions {
			action()
		}

		if doc.opts.gc {
			ds.TryGcDeleteSet(store, func(item *Item) bool { return true })
		}

		ds.TryMergeDeleteSet(store)

		for client, clock := range tx.afterState {
			beforeClock, ok := tx.beforeState[client]
			if !ok {
				beforeClock = 0
			}

			if beforeClock != clock {
				structs := store.GetStructs(client)
				beforeClockStructIndex, _ := findIndexSS(structs, beforeClock)
				firstChangePos := max(beforeClockStructIndex, 1)
				for i := len(structs) - 1; i >= int(firstChangePos); {
					i -= 1 + store.MergeWithLefts(client, i)
				}
			}
		}

		for i := len(tx.mergeStructs) - 1; i >= 0; i-- {
			id := tx.mergeStructs[i].Id()
			client := id.client
			clock := id.clock
			structs := store.GetStructs(client)
			replacedStructPos, _ := findIndexSS(structs, clock)
			if int(replacedStructPos+1) < len(structs) {
				if store.MergeWithLefts(client, int(replacedStructPos)+1) > 1 {
					continue
				}
			}

			if replacedStructPos > 0 {
				store.MergeWithLefts(client, int(replacedStructPos))
			}
		}

		if !tx.local {
			afterClock, afterOk := tx.afterState[doc.opts.clientId]
			beforeClock, beforeOk := tx.beforeState[doc.opts.clientId]
			if afterOk && beforeOk && afterClock != beforeClock {
				doc.opts.clientId = generateNewClientId()
			}
		}

		doc.Emit("afterTransactionCleanup", tx, doc)
		encoder := newUpdateEncoderV1()
		hasContent, err := tx.WriteMessage(encoder)
		if err != nil {
			return err
		}
		if hasContent {
			buf, err := encoder.ToUint8Array()
			if err != nil {
				return err
			}
			doc.Emit("update", buf, tx.origin, tx.doc, tx)
		}

		if len(transactionCleanups) <= index+1 {
			doc.transactionCleanups = []*Transaction{}
		} else {
			cleanupTransactions(transactionCleanups, index+1)
		}
	}
	return err
}

func (tx *Transaction) WriteMessage(encoder UpdateEncoder) (bool, error) {
	dslen := len(tx.deleteSet.clients)
	changedClocks := false
	for client, clock := range tx.afterState {
		if cd, ok := tx.beforeState[client]; !ok || cd != clock {
			changedClocks = true
			break
		}
	}
	if dslen == 0 && !changedClocks {
		return false, nil
	}

	tx.deleteSet.SortAndMergeDeleteSet()
	if err := tx.doc.store.WriteClientStructs(encoder, tx.beforeState); err != nil {
		return false, err
	}
	if err := tx.deleteSet.Write(encoder); err != nil {
		return false, err
	}

	return true, nil
}
