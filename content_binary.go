package ynotgo

type ContentBinary struct {
	content []byte
}

func (content *ContentBinary) Length() int {
	return 1
}

func (content *ContentBinary) Content() []any {
	return []any{content.content}
}

func (content *ContentBinary) Countable() bool {
	return true
}

func (content *ContentBinary) Copy() Content {
	return &ContentBinary{
		content: content.content,
	}
}

func (content *ContentBinary) Splice(offset uint32) (Content, error) {
	return nil, nil
}

func (content *ContentBinary) MergeWith(right Content) (bool, error) {
	return false, nil
}

func (content *ContentBinary) Integrate(tx *Transaction, item *Item) {}

func (content *ContentBinary) Delete(tx *Transaction) {}

func (content *ContentBinary) Gc(store *StructStore) {}

// TODO: Implement this once the encoder has been taken care of
func (content *ContentBinary) Write() {}

func (content *ContentBinary) Ref() uint8 {
	return 3
}
