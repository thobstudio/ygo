package ynotgo

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/thobstudio/ynotgo/lib0"
)

type UpdateEncoderV1 struct {
	buf    *bytes.Buffer
	writer *bufio.Writer
}

func newUpdateEncoderV1() *UpdateEncoderV1 {
	buf := bytes.NewBuffer(nil)
	return &UpdateEncoderV1{
		buf:    buf,
		writer: bufio.NewWriter(buf),
	}
}

func (e UpdateEncoderV1) ResetDsCurVal() error {
	// This is noop in v1
	return nil
}

func (e UpdateEncoderV1) WriteDsClock(num uint32) error {
	return lib0.WriteVarUint(e.writer, num)
}

func (e UpdateEncoderV1) WriteDsLen(num uint32) error {
	return lib0.WriteVarUint(e.writer, num)
}

func (e UpdateEncoderV1) ToUint8Array() ([]byte, error) {
	if err := e.writer.Flush(); err != nil {
		return nil, err
	}
	return e.buf.Bytes(), nil
}

func (e UpdateEncoderV1) Writer() *bufio.Writer {
	return e.writer
}

func (e UpdateEncoderV1) WriteLeftId(id *ID) error {
	if err := lib0.WriteVarUint(e.writer, id.client); err != nil {
		return nil
	}
	if err := lib0.WriteVarUint(e.writer, id.clock); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteRightId(id *ID) error {
	if err := lib0.WriteVarUint(e.writer, id.client); err != nil {
		return nil
	}
	if err := lib0.WriteVarUint(e.writer, id.clock); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteClient(client uint32) error {
	if err := lib0.WriteVarUint(e.writer, client); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteInfo(info uint32) error {
	if err := lib0.WriteVarUint(e.writer, info); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteString(str string) error {
	if err := lib0.WriteVarString(e.writer, str); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteParentInfo(isYKey bool) error {
	var b uint32 = 0
	if isYKey {
		b = 1
	}
	if err := lib0.WriteVarUint(e.writer, b); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteTypeRef(ref uint32) error {
	if err := lib0.WriteVarUint(e.writer, ref); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteLen(length uint32) error {
	if err := lib0.WriteVarUint(e.writer, length); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteAny(object any) error {
	if err := lib0.WriteAny(e.writer, object); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteBuf(buf []byte) error {
	if err := lib0.WriteUint8Array(e.writer, buf); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteJSON(object any) error {
	buf, err := json.Marshal(object)
	if err != nil {
		return err
	}
	if err := lib0.WriteUint8Array(e.writer, buf); err != nil {
		return nil
	}
	return nil
}

func (e UpdateEncoderV1) WriteKey(key string) error {
	if err := lib0.WriteVarString(e.writer, key); err != nil {
		return nil
	}
	return nil
}
