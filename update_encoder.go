package ynotgo

import "bufio"

type DsEncoder interface {
	ToUint8Array() ([]byte, error)
	ResetDsCurVal() error
	WriteDsClock(num uint32) error
	WriteDsLen(num uint32) error
}

type UpdateEncoder interface {
	DsEncoder
	Writer() *bufio.Writer
	WriteLeftId(id *ID) error
	WriteRightId(id *ID) error
	WriteClient(client uint32) error
	WriteInfo(info uint32) error
	WriteString(str string) error
	WriteParentInfo(isYKey bool) error
	WriteTypeRef(info uint32) error
	WriteLen(length uint32) error
	WriteAny(object any) error
	WriteBuf(buf []byte) error
	WriteJSON(object any) error
	WriteKey(key string) error
}
