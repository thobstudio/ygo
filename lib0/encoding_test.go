package lib0

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodingVarUint(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	WriteVarUint(writer, uint32(256))
	writer.Flush()

	num, _ := ReadVarUint(bufio.NewReader(buf))
	assert.Equal(t, uint32(256), num)
}

func TestEncodingVarint(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	WriteVarint(writer, int32(-45000))
	writer.Flush()

	num, _ := ReadVarInt(bufio.NewReader(buf))
	assert.Equal(t, int32(-45000), num)
}

func TestEncodingFloat32(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	WriteFloat32(writer, 32.000012)
	writer.Flush()

	num, _ := ReadFloat32(bufio.NewReader(buf))
	assert.Equal(t, float32(32.000012), num)
}

func TestEncodingFloat64(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	WriteFloat64(writer, 32.000012)
	writer.Flush()

	num, _ := ReadFloat64(bufio.NewReader(buf))
	assert.Equal(t, float64(32.000012), num)
}

func TestEncodingUint8Array(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	WriteUint8Array(writer, []byte{0x80, 0x40})
	writer.Flush()

	num, _ := ReadUint8Array(bufio.NewReader(buf))
	assert.Equal(t, []byte{0x80, 0x40}, num)
}

func TestEncodingString(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	WriteVarString(writer, "hello world 🍆")
	writer.Flush()

	num, _ := ReadVarString(bufio.NewReader(buf))
	assert.Equal(t, "hello world 🍆", num)
}

func TestEncodingStringArray(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	payload := []any{"one", "two", "three"}
	WriteAny(writer, payload)
	writer.Flush()

	num, _ := ReadAny(bufio.NewReader(buf))
	assert.Equal(t, payload, num)
}

func TestEncodingBoolMap(t *testing.T) {
	buf := bytes.NewBuffer(nil)
	writer := bufio.NewWriter(buf)
	payload := make(map[string]any, 2)
	payload["true"] = true
	payload["false"] = false
	WriteAny(writer, payload)
	writer.Flush()

	num, _ := ReadAny(bufio.NewReader(buf))
	assert.Equal(t, payload, num)
}
