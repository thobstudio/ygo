package ynotgo

import "errors"

type ContentAny struct {
	arr []any
}

func (content *ContentAny) Length() int {
	return len(content.arr)
}

func (content *ContentAny) Content() []any {
	return content.arr
}

func (content *ContentAny) Countable() bool {
	return true
}

func (content *ContentAny) Copy() Content {
	return &ContentAny{
		arr: content.arr,
	}
}

func (content *ContentAny) Splice(offset uint32) (Content, error) {
	right := &ContentAny{
		arr: content.arr[offset:],
	}
	content.arr = content.arr[:offset]
	return right, nil
}

func (content *ContentAny) MergeWith(right Content) (bool, error) {
	if c, ok := right.(*ContentAny); ok {
		content.arr = append(content.arr, c.arr...)
		return true, nil
	}
	return false, errors.New("unexpected type")
}

func (content *ContentAny) Integrate(tx *Transaction, item *Item) {}

func (content *ContentAny) Delete(tx *Transaction) {}

func (content *ContentAny) Gc(store *StructStore) {}

// TODO: Implement once the encoder has been taken care of
func (content *ContentAny) Write() {}

func (content *ContentAny) Ref() uint8 {
	return 8
}
