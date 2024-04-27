package ynotgo

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/thobstudio/ynotgo/lib0"
)

type DsDecoderV1 struct {
	buf    *bytes.Buffer
	reader *bufio.Reader
}

func newDsDecoderV1(update []byte) *DsDecoderV1 {
	buf := bytes.NewBuffer(update)
	return &DsDecoderV1{
		buf:    buf,
		reader: bufio.NewReader(buf),
	}
}

func NewDsDecoderV1(update []byte) *DsDecoderV1 {
	return newDsDecoderV1(update)
}

func (d *DsDecoderV1) Reader() *bufio.Reader {
	return d.reader
}

func (d *DsDecoderV1) ResetDsCurVal() error {
	return nil // This is a noop in v1
}

func (d *DsDecoderV1) ReadDsClock() (uint32, error) {
	return lib0.ReadVarUint(d.reader)
}

func (d *DsDecoderV1) ReadDsLen() (uint32, error) {
	return lib0.ReadVarUint(d.reader)
}

/*---------------------------------------------------------------------*/

type UpdateDecoderV1 struct {
	DsDecoderV1
}

func newUpdateDecoderV1(update []byte) *UpdateDecoderV1 {
	return &UpdateDecoderV1{
		DsDecoderV1: *newDsDecoderV1(update),
	}
}

func NewUpdateDecoderV1(update []byte) *UpdateDecoderV1 {
	return newUpdateDecoderV1(update)
}

func (d *UpdateDecoderV1) ReadLeftId() (*ID, error) {
	client, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return nil, err
	}
	clock, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return nil, err
	}
	return NewID(client, clock), nil
}

func (d *UpdateDecoderV1) ReadRightId() (*ID, error) {
	client, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return nil, err
	}
	clock, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return nil, err
	}
	return NewID(client, clock), nil
}

func (d *UpdateDecoderV1) ReadClient() (uint32, error) {
	client, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return client, nil
}

func (d *UpdateDecoderV1) ReadInfo() (uint32, error) {
	info, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return info, nil
}

func (d *UpdateDecoderV1) ReadString() (string, error) {
	s, err := lib0.ReadVarString(d.reader)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (d *UpdateDecoderV1) ReadParentInfo() (bool, error) {
	info, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return false, err
	}
	return info == 1, nil
}

func (d *UpdateDecoderV1) ReadTypeRef() (uint32, error) {
	ref, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return ref, nil
}

func (d *UpdateDecoderV1) ReadLen() (uint32, error) {
	length, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return length, nil
}

func (d *UpdateDecoderV1) ReadAny() (any, error) {
	a, err := lib0.ReadAny(d.reader)
	if err != nil {
		return 0, err
	}
	return a, nil
}

func (d *UpdateDecoderV1) ReadBuf() ([]byte, error) {
	buf, err := lib0.ReadUint8Array(d.reader)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

func (d *UpdateDecoderV1) ReadJSON() (any, error) {
	buf, err := lib0.ReadUint8Array(d.reader)
	if err != nil {
		return []byte{}, err
	}
	var j any
	err = json.Unmarshal(buf, &j)
	if err != nil {
		return nil, err
	}
	return j, nil
}

func (d *UpdateDecoderV1) ReadKey() (string, error) {
	s, err := lib0.ReadVarString(d.reader)
	if err != nil {
		return "", err
	}
	return s, nil
}
