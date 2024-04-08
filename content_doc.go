package ynotgo

type ContentDoc struct {
	doc *Doc
}

func (content *ContentDoc) Length() int {
	return 1
}

func (content *ContentDoc) Content() []any {
	return []any{content.doc}
}

func (content *ContentDoc) Countable() bool {
	return true
}

func (content *ContentDoc) Copy() Content {
	return &ContentDoc{
		doc: newDoc(
			WithGuid(content.doc.guid),
		),
	}
}

func (content *ContentDoc) Splice(offset uint32) (Content, error) {
	return nil, nil
}

func (content *ContentDoc) MergeWith(right Content) (bool, error) {
	return false, nil
}

func (content *ContentDoc) Integrate(tx *Transaction, item *Item) {
	content.doc.item = item
	tx.subdocsAdded[content.doc] = true
	if content.doc.shouldLoad {
		tx.subdocsLoaded[content.doc] = true
	}
}

func (content *ContentDoc) Delete(tx *Transaction) {
	if _, ok := tx.subdocsAdded[content.doc]; ok {
		delete(tx.subdocsAdded, content.doc)
	} else {
		tx.subdocsRemoved[content.doc] = true
	}
}

func (content *ContentDoc) Gc(store *StructStore) {}

// TODO: Implement this once we have taken care of encoder
func (content *ContentDoc) Write() {}

func (content *ContentDoc) Ref() uint8 {
	return 9
}
