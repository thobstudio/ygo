package ynotgo

import (
	"errors"
)

type ContentAny struct {
	arr []any
}

func newContentAny(arr []any) *ContentAny {
	return &ContentAny{
		arr: arr,
	}
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

func (content *ContentAny) Write(encoder UpdateEncoder, offset uint32) error {
	length := len(content.arr)
	if err := encoder.WriteLen(uint32(length)); err != nil {
		return err
	}

	for i := int(offset); i < length; i++ {
		c := content.arr[i]
		encoder.WriteAny(c)
	}

	return nil
}

func (content *ContentAny) Ref() uint8 {
	return 8
}

func readContentAny(decoder UpdateDecoder) (*ContentAny, error) {
	length, err := decoder.ReadLen()
	if err != nil {
		return nil, err
	}
	cs := make([]any, length)
	for i := 0; i < int(length); i++ {
		a, err := decoder.ReadAny()
		if err != nil {
			return nil, err
		}
		cs = append(cs, a)
	}

	return newContentAny(cs), nil
}
