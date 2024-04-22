package ynotgo

import "bufio"

type DsDecoder interface {
	ResetDsClock() error
	ReadDsClock() error
	ReadDsLen() error
}

type UpdateDecoder interface {
	DsDecoder
	Reader() *bufio.Reader
	ReadLeftId() (*ID, error)
	ReadRightId() (*ID, error)
	ReadClient() (uint32, error)
	ReadInfo() (uint8, error)
	ReadString() (string, error)
	ReadParentInfo() (uint8, error)
	ReadTypeRef() (uint32, error)
	ReadLen() (uint32, error)
	ReadAny() (any, error)
	ReadBuf() ([]byte, error)
	ReadJSON() (any, error)
	ReadKey() (string, error)
}
