package ynotgo

import (
	"bufio"
	"bytes"
	"encoding/json"

	"github.com/thobstudio/ynotgo/lib0"
)

type UpdateEncoderV1 struct {
	buf     *bytes.Buffer
	writter *bufio.Writer
}

func newUpdateEncoderv1() *UpdateEncoderV1 {
	buf := &bytes.Buffer{}
	return &UpdateEncoderV1{
		buf:     buf,
		writter: bufio.NewWriter(buf),
	}
}

func (ue *UpdateEncoderV1) Bytes() ([]byte, error) {
	if err := ue.writter.Flush(); err != nil {
		return []byte{}, err
	}
	return ue.buf.Bytes(), nil
}

func (ue *UpdateEncoderV1) ResetDsCurVal() {}

func (ue *UpdateEncoderV1) WriteDsClock(clock uint32) error {
	return lib0.WriteVarUint(ue.writter, clock)
}

func (ue *UpdateEncoderV1) WriteDsLen(length uint32) error {
	return lib0.WriteVarUint(ue.writter, length)
}

func (ue UpdateEncoderV1) WriteLeftId(id *ID) error {
	if err := lib0.WriteVarUint(ue.writter, id.client); err != nil {
		return err
	}
	if err := lib0.WriteVarUint(ue.writter, id.clock); err != nil {
		return err
	}
	return nil
}

func (ue UpdateEncoderV1) WriteRightId(id *ID) error {
	if err := lib0.WriteVarUint(ue.writter, id.client); err != nil {
		return err
	}
	if err := lib0.WriteVarUint(ue.writter, id.clock); err != nil {
		return err
	}
	return nil
}

func (ue UpdateEncoderV1) WriteClient(client uint32) error {
	return lib0.WriteVarUint(ue.writter, client)
}

func (ue UpdateEncoderV1) WriteInfo(info uint8) error {
	return lib0.WriteUint8(ue.writter, info)
}

func (ue UpdateEncoderV1) WriteString(str string) error {
	return lib0.WriteVarString(ue.writter, str)
}

func (ue UpdateEncoderV1) WriteParentInfo(isKey bool) error {
	var b uint32 = 0
	if isKey {
		b = 1
	}
	return lib0.WriteVarUint(ue.writter, b)
}

func (ue UpdateEncoderV1) WriteTypeRef(ref uint32) error {
	return lib0.WriteVarUint(ue.writter, ref)
}

func (ue UpdateEncoderV1) WriteLength(length uint32) error {
	return lib0.WriteVarUint(ue.writter, length)
}

func (ue UpdateEncoderV1) WriteAny(object any) error {
	return lib0.WriteAny(ue.writter, object)
}

func (ue UpdateEncoderV1) WriteBuf(buf []byte) error {
	return lib0.WriteAny(ue.writter, buf)
}

func (ue UpdateEncoderV1) WriteJson(object any) error {
	str, err := json.Marshal(object)
	if err != nil {
		return err
	}
	if err := lib0.WriteVarString(ue.writter, string(str)); err != nil {
		return err
	}
	return nil
}

func (ue UpdateEncoderV1) WriteKey(key string) error {
	return lib0.WriteVarString(ue.writter, key)
}
