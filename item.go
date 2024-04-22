package ynotgo

import (
	"errors"

	"github.com/thobstudio/ynotgo/lib0"
)

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
	info        uint32
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
	var info uint32 = 0
	if content.Countable() {
		info = lib0.Bit2
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
	if offset > 0 {
		item.id.clock += offset
		left, err := tx.doc.store.GetItemCleanEnd(tx, newId(item.id.client, item.id.clock-1))
		if err != nil {
			return err
		}
		item.left = left.(*Item)
		item.leftOrigin = item.left.LastId()
		content, err := item.content.Splice(offset)
		if err != nil {
			return err
		}
		item.content = content
		item.length -= offset
	}

	if item.parent != nil {
		parent, _ := item.parent.(SharedType)
		if (item.left == nil && (item.right == nil || item.right.left != nil)) || (item.left != nil && item.left.right != item.right) {
			left := item.left

			var o *Item

			if left != nil {
				o = left.right
			} else if item.parentSub != "" {
				o = item.parent.(SharedType).GetItem(item.parentSub)
				for o != nil && o.left != nil {
					o = o.left
				}
			} else {
				o = item.parent.(SharedType).Start()
			}

			conflictingItems := make(map[*Item]bool)
			itemsBeforeOrigin := make(map[*Item]bool)

			for o != nil && o != item.right {
				itemsBeforeOrigin[o] = true
				conflictingItems[o] = true

				if compareIds(item.leftOrigin, o.leftOrigin) {
					if o.id.client < item.id.client {
						left = o
						conflictingItems = make(map[*Item]bool)
					} else if compareIds(item.rightOrigin, o.rightOrigin) {
						break
					}
				} else if o.leftOrigin != nil {
					leftOriginItem, err := tx.doc.store.GetItem(o.leftOrigin)
					if err != nil {
						return err
					}
					if hasItemBeforeOrigin, ok := itemsBeforeOrigin[leftOriginItem.(*Item)]; ok && hasItemBeforeOrigin {
						if hasItemConflicting, ok := conflictingItems[leftOriginItem.(*Item)]; !ok && !hasItemConflicting {
							left = o
							conflictingItems = make(map[*Item]bool)
						}
					}
				} else {
					break
				}
				o = o.right
			}

			item.left = left
		}

		if item.left != nil {
			right := item.left.right
			item.right = right
			item.left.right = item
		} else {
			var r *Item
			if item.parentSub != "" {
				r = parent.GetItem(item.parentSub)
				for r != nil && r.left != nil {
					r = r.left
				}
			} else {
				r = parent.Start()
				parent.SetStart(item)
			}

			item.right = r
		}

		if item.right != nil {
			item.right.left = item
		} else if item.parentSub != "" {
			parent.SetItem(item.parentSub, item)
			if item.left != nil {
				item.left.Delete(tx)
			}
		}

		if item.parentSub == "" && item.Countable() && !item.Deleted() {
			parent.SetLength(parent.Length() + item.length)
		}

		tx.doc.store.AddStructItem(item)
		item.content.Integrate(tx, item)

		tx.AddChangedType(item.parent.(SharedType), item.parentSub)

		if parent.Item() != nil && parent.Item().Deleted() || (item.parentSub != "" && item.right != nil) {
			item.Delete(tx)
		}
	} else {
		// If parent is not defined. We need Integrate GC structs instead
		newGc(item.id, item.length).Integrate(tx, offset)
	}
	return nil
}

func (item *Item) Gc(store *StructStore, parentGcd bool) error {
	if !item.Deleted() {
		return errors.New("Gc item not deleted yet")
	}

	item.content.Gc(store)
	if parentGcd {
		if err := store.ReplaceStruct(item, newGc(item.id, item.length)); err != nil {
			return err
		}
	} else {
		item.content = newContentDeleted(item.length)
	}
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

func (item *Item) SplitItem(tx *Transaction, offset uint32) *Item {
	client := item.id.client
	clock := item.id.clock
	content, _ := item.content.Splice(offset)
	rightItem := newItem(newId(client, clock+offset), item, item.right, newId(client, clock+offset-1), item.rightOrigin, item.parent, item.parentSub, content)

	if item.Deleted() {
		rightItem.markDeleted()
	}

	if item.Keep() {
		rightItem.SetKeep(true)
	}

	item.right = rightItem
	if rightItem.right != nil {
		rightItem.right.left = rightItem
	}

	tx.mergeStructs = append(tx.mergeStructs, rightItem)

	if rightItem.parentSub != "" && rightItem.right == nil {
		rightItem.parent.(SharedType).SetItem(rightItem.parentSub, rightItem)
	}

	item.length = offset

	return rightItem
}

func (item *Item) LastId() *ID {
	if item.length == 1 {
		return item.id
	}

	return newId(item.id.client, item.id.clock+item.length-1)
}

func (item *Item) Deleted() bool {
	// BIT3 is bitmask for deleted
	return item.info&lib0.Bit3 > 0
}

func (item *Item) SetDeleted(doDelete bool) {
	if item.Deleted() != doDelete {
		item.info ^= lib0.Bit3
	}
}

func (item *Item) markDeleted() {
	item.info |= lib0.Bit3
}

func (item *Item) Countable() bool {
	return item.info&lib0.Bit2 > 0
}

func (item *Item) Keep() bool {
	return item.info&lib0.Bit1 > 0
}

func (item *Item) SetKeep(doKeep bool) {
	if item.Keep() != doKeep {
		item.info ^= lib0.Bit1
	}
}

func (item *Item) Write(encoder UpdateEncoder, offset uint32) error {
	return nil
}
