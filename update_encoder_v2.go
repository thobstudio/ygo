package ygo

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/thobstudio/ygo/lib0"
)

type DsEncoderV2 struct {
	buf    *bytes.Buffer
	writer *bufio.Writer
}

func newDsEncoderV2() *DsEncoderV2 {
	buf := bytes.NewBuffer(nil)
	return &DsEncoderV2{
		buf:    buf,
		writer: bufio.NewWriter(buf),
	}
}

func NewDsEncoderV2() *DsEncoderV2 {
	return newDsEncoderV2()
}

func (e DsEncoderV2) ResetDsCurVal() error {
	// This is noop in v1
	return nil
}

func (e DsEncoderV2) WriteDsClock(num uint) error {
	return lib0.WriteVarUint(e.writer, num)
}

func (e DsEncoderV2) WriteDsLen(num uint) error {
	return lib0.WriteVarUint(e.writer, num)
}

func (e DsEncoderV2) ToUint8Array() ([]byte, error) {
	if err := e.writer.Flush(); err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

func (e DsEncoderV2) Writer() *bufio.Writer {
	return e.writer
}

/*--------------------------------------------------------------------------*/

type UpdateEncoderV2 struct {
	*DsEncoderV2
}

func newUpdateEncoderV2() *UpdateEncoderV2 {
	return &UpdateEncoderV2{
		DsEncoderV2: newDsEncoderV2(),
	}
}

func NewUpdateEncoderV2() *UpdateEncoderV2 {
	return newUpdateEncoderV2()
}

func (e UpdateEncoderV2) ToUint8Array() ([]byte, error) {
	if err := e.writer.Flush(); err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

func (e UpdateEncoderV2) WriteLeftId(id *ID) error {
	if err := lib0.WriteVarUint(e.writer, id.client); err != nil {
		return nil
	}
	if err := lib0.WriteVarUint(e.writer, id.clock); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteRightId(id *ID) error {
	if err := lib0.WriteVarUint(e.writer, id.client); err != nil {
		return nil
	}
	if err := lib0.WriteVarUint(e.writer, id.clock); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteClient(client uint) error {
	if err := lib0.WriteVarUint(e.writer, client); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteInfo(info uint) error {
	if err := lib0.WriteVarUint(e.writer, info); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteString(str string) error {
	if err := lib0.WriteVarString(e.writer, str); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteParentInfo(isYKey bool) error {
	var b uint = 0
	if isYKey {
		b = 1
	}
	if err := lib0.WriteVarUint(e.writer, b); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteTypeRef(ref uint) error {
	if err := lib0.WriteVarUint(e.writer, ref); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteLen(length uint) error {
	if err := lib0.WriteVarUint(e.writer, length); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteAny(object any) error {
	if err := lib0.WriteAny(e.writer, object); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteBuf(buf []byte) error {
	if err := lib0.WriteUint8Array(e.writer, buf); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteJSON(object any) error {
	buf, err := json.Marshal(object)
	if err != nil {
		return err
	}
	if err := lib0.WriteUint8Array(e.writer, buf); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV2) WriteKey(key string) error {
	if err := lib0.WriteVarString(e.writer, key); err != nil {
		return nil
	}
	return nil
}
