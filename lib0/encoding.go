package lib0

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"math"
	"reflect"
	"unsafe"
)

func WriteVarUint(writter *bufio.Writer, num uint32) error {
	for num > Bits7 {
		if err := writter.WriteByte((byte)(Bit8 | (Bits7 & num))); err != nil {
			return err
		}
		num >>= 7
	}
	return writter.WriteByte((byte)(Bits7 & num))
}

func WriteVarUint8Array(writter *bufio.Writer, buf []byte) error {
	if err := WriteVarUint(writter, uint32(len(buf))); err != nil {
		return err
	}
	if _, err := writter.Write(buf); err != nil {
		return err
	}
	return nil
}

func WriteVarString(writter *bufio.Writer, str string) error {
	buf := []byte(str)
	return WriteVarUint8Array(writter, buf)
}

func WriteVarInt(writer *bufio.Writer, num int, treatZeroAsNegative bool) {
	isNegative := treatZeroAsNegative
	if num < 0 {
		isNegative = true
	}
	if isNegative {
		num = -num
	}

	// |   whether to continue reading   |         is negative         | value.
	var tmpA uint32
	var tmpB uint32
	if uint32(num) > Bits6 {
		tmpA = Bit8
	}
	if isNegative {
		tmpB = Bit7
	}

	// |   whether to continue reading   |         is negative         | value.
	writer.WriteByte((byte)(tmpA | tmpB | Bit6&uint32(num)))
	num >>= 6

	// We don't need to consider the case of num == 0 so we can use a different pattern here than above.
	var tmpC uint32
	if uint32(num) > Bits7 {
		tmpC = Bit8
	}
	for num > 0 {
		writer.WriteByte((byte)(tmpC | (Bits7 & uint32(num))))
		num >>= 7
	}
}

func WriteAny(writter *bufio.Writer, o any) error {
	switch t := o.(type) {
	case string:
		if err := writter.WriteByte(119); err != nil {
			return err
		}
		if err := WriteVarString(writter, t); err != nil {
			return err
		}
	case bool:
		var b uint32 = 121
		if t {
			b = 120
		}
		if err := writter.WriteByte((byte)(b)); err != nil {
			return err
		}
	case float64: // TYPE 123: FLOAT64
		bits := math.Float64bits(t)
		bytes := make([]byte, 8)
		binary.LittleEndian.PutUint64(bytes, bits)
		if !IsLittleEndian() {
			bytes = Reverse(bytes)
		}
		writter.WriteByte(123)
		writter.Write(bytes)
	case float32: // TYPE 124: FLOAT32
		bits := math.Float32bits(t)
		bytes := make([]byte, 4)
		binary.LittleEndian.PutUint32(bytes, bits)
		if !IsLittleEndian() {
			bytes = Reverse(bytes)
		}
		writter.WriteByte(124)
		writter.Write(bytes)
	case int: // TYPE 125: INTEGER
		writter.WriteByte(125)
		WriteVarInt(writter, t, false)
	case int64: // Special case: treat LONG as INTEGER.
		writter.WriteByte(125)
		WriteVarInt(writter, int(t), false)
	case nil: // TYPE 126: null
		// TYPE 127: undefined
		writter.WriteByte(126)
	case []byte: // TYPE 116: ArrayBuffer
		writter.WriteByte(116)
		WriteVarUint8Array(writter, t)
	case map[string]any: // TYPE 118: object (Dictionary<string, object>)
		writter.WriteByte(118)
		WriteVarUint(writter, uint32(len(t)))
		for key, value := range t {
			WriteVarString(writter, key)
			WriteAny(writter, value)
		}
	case []any: // TYPE 117: Array
		writter.WriteByte(117)
		WriteVarUint(writter, uint32(len(t)))
		for _, item := range t {
			WriteAny(writter, item)
		}
	default:
		return errors.New(fmt.Sprintf("unsupported type : %v \n", reflect.TypeOf(o)))
	}
	return errors.New(fmt.Sprintf("unsupported type : %v \n", reflect.TypeOf(o)))
}

func WriteUint8(writter *bufio.Writer, b byte) error {
	return writter.WriteByte(b)
}

func IsLittleEndian() bool {
	var value int32 = 1 // 占4byte 转换成16进制 0x00 00 00 01
	// 大端(16进制)：00 00 00 01
	// 小端(16进制)：01 00 00 00
	pointer := unsafe.Pointer(&value)
	pb := (*byte)(pointer)
	if *pb != 1 {
		return false
	}
	return true
}

func Reverse(bytes []byte) []byte {
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return bytes
}
