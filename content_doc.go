package ynotgo

type DocOpts struct {
	gc      bool
	autLoad bool
	meta    any
}

type ContentDoc struct {
	doc  *Doc
	opts *DocOpts
}

func newContentDoc(doc *Doc) *ContentDoc {
	opts := &DocOpts{
		gc:      true,
		autLoad: false,
	}
	if !doc.gc {
		opts.gc = true
	}
	if doc.autoLoad {
		opts.autLoad = true
	}
	if doc.meta != nil {
		opts.meta = doc.meta
	}
	return &ContentDoc{
		doc:  doc,
		opts: opts,
	}
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
	// TODO: to be implemented
}

func (content *ContentDoc) Delete(tx *Transaction) {
	// TODO: to be implemented
}

func (content *ContentDoc) Gc(store *StructStore) {}

// TODO: Implement this once we have taken care of encoder
func (content *ContentDoc) Write(encoder UpdateEncoder, offset uint32) error {
	if err := encoder.WriteString(content.doc.guid); err != nil {
		return err
	}
	if err := encoder.WriteAny(content.opts); err != nil {
		return err
	}
	return nil
}

func (content *ContentDoc) Ref() uint8 {
	return 9
}
