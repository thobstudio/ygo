package ynotgo

type YMap struct {
	AbstractType
	prelimContent map[string]interface{}
}

func newYMap() *YMap {
	return &YMap{
		AbstractType: AbstractType{
			itemMap: make(map[string]*Item),
		},
	}
}

func NewYMap() *YMap {
	return newYMap()
}

func (m *YMap) Item() *Item {
	return m.item
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
func typeMapSet(m *YMap, key string, value any, tx *Transaction) error {
	left := m.GetItem(key)
	doc := tx.doc
	ownClientId := doc.clientId
	var content Content

	if value == nil {
		content = newContentAny([]any{value})
	} else {
		switch val := value.(type) {
		case uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, string, bool:
			content = newContentAny([]any{val})
		case []byte:
			content = newContentBinary(val)
		case *Doc:
			content = newContentDoc(val)
		case SharedType:
			content = newContentType(val)
		default:
			return errors.New("unexpected content type")
		}
	}

	id := newId(ownClientId, doc.store.State(ownClientId))
	var leftId *ID
	if left != nil {
		leftId = left.LastId()
	}

	item := newItem(id, left, nil, leftId, nil, m, key, content)
	return item.Integrate(tx, 0)
}

func typeMapDelete(m *YMap, tx *Transaction, key string) error {
	item := m.GetItem(key)
	if item != nil {
		item.Delete(tx)
	}

	return nil
}
