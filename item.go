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
	}
}

