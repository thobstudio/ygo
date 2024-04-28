package ynotgo

type ContentType struct {
	contentType SharedType
}

func newContentType(contentType SharedType) *ContentType {
	return &ContentType{
		contentType: contentType,
	}
}

func (content *ContentType) Length() int {
	return 1
}

func (content *ContentType) Content() []any {
	return []any{content.contentType}
}

func (content *ContentType) Countable() bool {
	return true
}

func (content *ContentType) Copy() Content {
	return &ContentType{
		contentType: content.contentType.Copy(),
	}
}

func (content *ContentType) Splice(offset uint32) (Content, error) {
	return nil, nil
}

func (content *ContentType) MergeWith(right Content) (bool, error) {
	return false, nil
}

func (content *ContentType) Integrate(tx *Transaction, item *Item) {
	content.contentType.Integrate(tx.doc, item)
}

func (content *ContentType) Delete(tx *Transaction) {
	panic("not implemented")
}

func (content *ContentType) Gc(store *StructStore) {
	item := content.contentType.Start()
	for item != nil {
		item.Gc(store, true)
		item = item.right
	}
	content.contentType.SetStart(nil)
	content.contentType.ForEachItem(func(item *Item, key string) {
		for item != nil {
			item.Gc(store, true)
			item = item.left
		}
	})
	content.contentType.ClearItemMap()
}

func (content *ContentType) Write(encoder UpdateEncoder, offset uint32) error {
	return content.contentType.Write(encoder)
}

func (content *ContentType) Ref() uint8 {
	return 7
}
