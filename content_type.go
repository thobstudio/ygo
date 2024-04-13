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

// TODO: Need to implement this
func (content *ContentType) Delete(tx *Transaction) {}

// TODO: Need to implement this
func (content *ContentType) Gc(store *StructStore) {}

// TODO: Implement this once we have taken care of encoder
func (content *ContentType) Write() {}

func (content *ContentType) Ref() uint8 {
	return 7
}
