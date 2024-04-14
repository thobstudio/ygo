package ynotgo

type Item struct {
	id          *ID
	length      uint32
	left        *Item
	right       *Item
	leftOrigin  *ID
	rightOrigin *ID
	parent      interface{}
	parentSub   string
	content     Content
	info        byte
}

func newItem(
	id *ID,
	left *Item,
	right *Item,
	leftOrigin *ID,
	rightOrigin *ID,
	parent interface{},
	parentSub string,
	content Content,
) *Item {
	var info byte = 0
	if content.Countable() {
		info = BIT2
	}
	return &Item{
		id:          id,
		left:        left,
		right:       right,
		leftOrigin:  leftOrigin,
		rightOrigin: rightOrigin,
		parent:      parent,
		parentSub:   parentSub,
		content:     content,
		info:        info,
		length:      uint32(content.Length()),
	}
}

func (item *Item) Id() *ID {
	return item.id
}

func (item *Item) Length() uint32 {
	return item.length
}

func (item *Item) MergeWith(right SharedStruct) (bool, error) {
	return false, nil
}

func (item *Item) Integrate(tx *Transaction, offset uint32) error {
	return nil
}

func (item *Item) Delete(tx *Transaction) {
	if !item.Deleted() {
		parent := item.parent.(SharedType)
		if item.Countable() && item.parentSub != "" {
			parent.SetLength(parent.Length() - item.length)
		}
		item.markDeleted()
		tx.deleteSet.AddDeleteItem(item.id.client, item.id.clock, item.length)
		tx.AddChangedType(parent, item.parentSub)
		item.content.Delete(tx)
	}
}

func (item *Item) Deleted() bool {
	// BIT3 is bitmask for deleted
	return item.info&BIT3 > 0
}

func (item *Item) SetDeleted(doDelete bool) {
	if item.Deleted() != doDelete {
		item.info ^= BIT3
	}
}

func (item *Item) markDeleted() {
	item.info |= BIT3
}

func (item *Item) Countable() bool {
	return item.info&BIT2 > 0
}

func (item *Item) Keep() bool {
	return item.info&BIT1 > 0
}

func (item *Item) SetKeep(doKeep bool) {
	if item.Keep() != doKeep {
		item.info ^= BIT1
	}
}
