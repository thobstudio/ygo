package ynotgo

import "errors"

type YMap struct {
	AbstractType
	prelimContent map[string]any
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

func (m *YMap) Has(key string) bool {
	if val, ok := m.itemMap[key]; ok {
		if val != nil {
			return !val.Deleted()
		}
	}
	return false
}

func (m *YMap) Get(key string) any {
	if val, ok := m.itemMap[key]; ok {
		if val != nil && !val.Deleted() {
			return val.content.Content()[val.length-1]
		}
	}
	return nil
}

func (m *YMap) Set(key string, value any) {
	if m.doc != nil {
		m.doc.Transact(func(tx *Transaction) (any, error) {
			if err := typeMapSet(m, key, value, tx); err != nil {
				return nil, err
			}
			return nil, nil
		}, nil, true)
	} else {
		m.prelimContent[key] = value
	}
}

func (m *YMap) Delete(key string) {
	if m.doc != nil {
		m.doc.Transact(func(tx *Transaction) (any, error) {
			if err := typeMapDelete(m, tx, key); err != nil {
				return nil, err
			}
			return nil, nil
		}, nil, true)
	} else {
		delete(m.prelimContent, key)
	}
}

func (m *YMap) Clear() {
	if m.doc != nil {
		m.doc.Transact(func(tx *Transaction) (any, error) {
			for key, item := range m.itemMap {
				if !item.Deleted() {
					if err := typeMapDelete(m, tx, key); err != nil {
						return nil, err
					}
				}
			}
			return nil, nil
		}, nil, true)
	} else {
		m.prelimContent = make(map[string]interface{})
	}
}

func (m *YMap) Clone() SharedType {
	nm := newYMap()
	for key, item := range m.itemMap {
		if !item.Deleted() {
			value := item.content.Content()[item.length-1]
			m.Set(key, value)
		}
	}
	return nm
}

func (m *YMap) Copy() SharedType {
	return newYMap()
}

func (m *YMap) Write(encoder UpdateEncoder) error {
	return encoder.WriteTypeRef(YMapRefID)
}

func (m *YMap) ToJSON() any {
	jsonmap := make(map[string]any)
	for key, item := range m.itemMap {
		if !item.Deleted() {
			val := item.content.Content()[item.length-1]
			if v, ok := val.(SharedType); ok {
				jsonmap[key] = v.ToJSON()
			} else {
				jsonmap[key] = val
			}
		}
	}
	return jsonmap
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
