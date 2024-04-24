package ynotgo

import "bufio"

type DsDecoder interface {
	ResetDsCurVal() error
	ReadDsClock() (uint32, error)
	ReadDsLen() (uint32, error)
}

type UpdateDecoder interface {
	DsDecoder
	Reader() *bufio.Reader
	ReadLeftId() (*ID, error)
	ReadRightId() (*ID, error)
	ReadClient() (uint32, error)
	ReadInfo() (uint32, error)
	ReadString() (string, error)
	ReadParentInfo() (uint32, error)
	ReadTypeRef() (uint32, error)
	ReadLen() (uint32, error)
	ReadAny() (any, error)
	ReadBuf() ([]byte, error)
	ReadJSON() (any, error)
	ReadKey() (string, error)
}
