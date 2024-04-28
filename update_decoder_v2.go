package ynotgo

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/thobstudio/ynotgo/lib0"
)

type DsDecoderV2 struct {
	buf    *bytes.Buffer
	reader *bufio.Reader
}

func newDsDecoderV2(update []byte) *DsDecoderV2 {
	buf := bytes.NewBuffer(update)
	return &DsDecoderV2{
		buf:    buf,
		reader: bufio.NewReader(buf),
	}
}

func (d *DsDecoderV2) Reader() *bufio.Reader {
	return d.reader
}

func (d *DsDecoderV2) ResetDsCurVal() error {
	return nil // This is a noop in v1
}

func (d *DsDecoderV2) ReadDsClock() (uint32, error) {
	return lib0.ReadVarUint(d.reader)
}

func (d *DsDecoderV2) ReadDsLen() (uint32, error) {
	return lib0.ReadVarUint(d.reader)
}

/*--------------------------------------------------------------*/

type UpdateDecoderV2 struct {
	*DsDecoderV2
}

func newUpdateDecoderV2(update []byte) *UpdateDecoderV2 {
	return &UpdateDecoderV2{
		DsDecoderV2: newDsDecoderV2(update),
	}
}

func NewUpdateDecoderV2(update []byte) *UpdateDecoderV2 {
	return newUpdateDecoderV2(update)
}

func (d *UpdateDecoderV2) ReadLeftId() (*ID, error) {
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

func (d *UpdateDecoderV2) ReadRightId() (*ID, error) {
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

func (d *UpdateDecoderV2) ReadClient() (uint32, error) {
	client, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return client, nil
}

func (d *UpdateDecoderV2) ReadInfo() (uint32, error) {
	info, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return info, nil
}

func (d *UpdateDecoderV2) ReadString() (string, error) {
	s, err := lib0.ReadVarString(d.reader)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (d *UpdateDecoderV2) ReadParentInfo() (bool, error) {
	info, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return false, err
	}
	return info == 1, nil
}

func (d *UpdateDecoderV2) ReadTypeRef() (uint32, error) {
	ref, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return ref, nil
}

func (d *UpdateDecoderV2) ReadLen() (uint32, error) {
	length, err := lib0.ReadVarUint(d.reader)
	if err != nil {
		return 0, err
	}
	return length, nil
}

func (d *UpdateDecoderV2) ReadAny() (any, error) {
	a, err := lib0.ReadAny(d.reader)
	if err != nil {
		return 0, err
	}
	return a, nil
}

func (d *UpdateDecoderV2) ReadBuf() ([]byte, error) {
	buf, err := lib0.ReadUint8Array(d.reader)
	if err != nil {
		return []byte{}, err
	}
	return buf, nil
}

func (d *UpdateDecoderV2) ReadJSON() (any, error) {
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

func (d *UpdateDecoderV2) ReadKey() (string, error) {
	s, err := lib0.ReadVarString(d.reader)
	if err != nil {
		return "", err
	}
	return s, nil
}
