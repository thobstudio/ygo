package ygo

import "bufio"

type DsEncoder interface {
	Writer() *bufio.Writer
	ToUint8Array() ([]byte, error)
	ResetDsCurVal() error
	WriteDsClock(num uint) error
	WriteDsLen(num uint) error
}

type UpdateEncoder interface {
	DsEncoder
	WriteLeftId(id *ID) error
	WriteRightId(id *ID) error
	WriteClient(client uint) error
	WriteInfo(info uint) error
	WriteString(str string) error
	WriteParentInfo(isYKey bool) error
	WriteTypeRef(info uint) error
	WriteLen(length uint) error
	WriteAny(object any) error
	WriteBuf(buf []byte) error
	WriteJSON(object any) error
	WriteKey(key string) error
}
