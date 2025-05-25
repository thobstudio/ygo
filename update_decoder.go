package ygo

import "bufio"

type DsDecoder interface {
	Reader() *bufio.Reader
	ResetDsCurVal() error
	ReadDsClock() (uint, error)
	ReadDsLen() (uint, error)
}

type UpdateDecoder interface {
	DsDecoder
	ReadLeftId() (*ID, error)
	ReadRightId() (*ID, error)
	ReadClient() (uint, error)
	ReadInfo() (uint, error)
	ReadString() (string, error)
	ReadParentInfo() (bool, error)
	ReadTypeRef() (uint, error)
	ReadLen() (uint, error)
	ReadAny() (any, error)
	ReadBuf() ([]byte, error)
	ReadJSON() (any, error)
	ReadKey() (string, error)
}
