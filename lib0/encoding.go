package lib0

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"reflect"
)

func WriteVarUint(writer *bufio.Writer, x uint) error {
	for x >= 0x80 {
		if err := writer.WriteByte(byte(x) | 0x80); err != nil {
			return err
		}
		x >>= 7
	}
	if err := writer.WriteByte(byte(x)); err != nil {
		return nil
	}
	return nil
}

func WriteVarint(writer *bufio.Writer, x int32) error {
	isNegative := x < 0
	if isNegative {
		x = -x
	}

	cr := int32(0)
	if x > int32(Bits6) {
		cr = int32(Bit8)
	}
	neg := int32(0)
	if isNegative {
		neg = int32(Bit7)
	}
	if err := writer.WriteByte(byte(cr | neg | (int32(Bits6) & x))); err != nil {
		return err
	}
	x >>= 6
	for x > 0 {
		pre := int32(0)
		if x > int32(Bits7) {
			pre = int32(Bit8)
		}
		if err := writer.WriteByte(byte(pre | (int32(Bits7) & x))); err != nil {
			return err
		}
		x >>= 7
	}

	return nil
}

func WriteUint8Array(writer *bufio.Writer, buf []byte) error {
	length := len(buf)
	if err := WriteVarUint(writer, uint(length)); err != nil {
		return err
	}
	_, err := writer.Write(buf)
	return err
}

func WriteVarString(writer *bufio.Writer, s string) error {
	length := len(s)
	if err := WriteVarUint(writer, uint(length)); err != nil {
		return err
	}
	_, err := writer.Write([]byte(s))
	return err
}

func WriteFloat32(writer *bufio.Writer, f float32) error {
	return binary.Write(writer, binary.BigEndian, f)
}

func WriteFloat64(writer *bufio.Writer, f float64) error {
	return binary.Write(writer, binary.BigEndian, f)
}

/**
 * Encode data with efficient binary format.
 *
 * Differences to JSON:
 * • Transforms data to a binary format (not to a string)
 * • Encodes undefined, NaN, and ArrayBuffer (these can't be represented in JSON)
 * • Numbers are efficiently encoded either as a variable length integer, as a
 *   32 bit float, as a 64 bit float, or as a 64 bit bigint.
 *
 * Encoding table:
 *
 * | Data Type           | Prefix   | Encoding Method    | Comment |
 * | ------------------- | -------- | ------------------ | ------- |
 * | undefined           | 127      |                    | Functions, symbol, and everything that cannot be identified is encoded as undefined |
 * | null                | 126      |                    | |
 * | integer             | 125      | writeVarInt        | Only encodes 32 bit signed integers |
 * | float32             | 124      | writeFloat32       | |
 * | float64             | 123      | writeFloat64       | |
 * | bigint              | 122      | writeBigInt64      | |
 * | boolean (false)     | 121      |                    | True and false are different data types so we save the following byte |
 * | boolean (true)      | 120      |                    | - 0b01111000 so the last bit determines whether true or false |
 * | string              | 119      | writeVarString     | |
 * | object<string,any>  | 118      | custom             | Writes {length} then {length} key-value pairs |
 * | array<any>          | 117      | custom             | Writes {length} then {length} json values |
 * | Uint8Array          | 116      | writeVarUint8Array | We use Uint8Array for any kind of binary data |
 *
 * Reasons for the decreasing prefix:
 * We need the first bit for extendability (later we may want to encode the
 * prefix with writeVarUint). The remaining 7 bits are divided as follows:
 * [0-30]   the beginning of the data range is used for custom purposes
 *          (defined by the function that uses this library)
 * [31-127] the end of the data range is used for data encoding by
 *          lib0/encoding.js
 *
 * @param {Encoder} encoder
 * @param {undefined|null|number|bigint|boolean|string|Object<string,any>|Array<any>|Uint8Array} data
 */
func WriteAny(writer *bufio.Writer, o any) error {
	switch t := o.(type) {
	// 127
	// 126
	case nil:
		if err := writer.WriteByte(126); err != nil {
			return err
		}
	// 125
	case int:
		if err := writer.WriteByte(125); err != nil {
			return err
		}
		return WriteVarint(writer, int32(t))
	case int64:
		if err := writer.WriteByte(125); err != nil {
			return err
		}
		return WriteVarint(writer, int32(t))
		// 124
	case float32:
		if err := writer.WriteByte(124); err != nil {
			return err
		}
		return WriteFloat32(writer, t)
		// 123
	case float64:
		if err := writer.WriteByte(123); err != nil {
			return err
		}
		return WriteFloat64(writer, t)
		// 122 We dont handle this case
		// 121 false
		// 120 true
	case bool:
		var b byte = 121
		if t {
			b = 120
		}
		writer.WriteByte(b)
		// 119
	case string:
		if err := writer.WriteByte(119); err != nil {
			return err
		}
		return WriteVarString(writer, t)
		// 118
	case map[string]any:
		if err := writer.WriteByte(118); err != nil {
			return err
		}
		if err := WriteVarUint(writer, uint(len(t))); err != nil {
			return err
		}
		for key, value := range t {
			if err := WriteVarString(writer, key); err != nil {
				return err
			}
			if err := WriteAny(writer, value); err != nil {
				return err
			}
		}
		// 117
	case []any:
		if err := writer.WriteByte(117); err != nil {
			return err
		}
		if err := WriteVarUint(writer, uint(len(t))); err != nil {
			return err
		}
		for _, item := range t {
			if err := WriteAny(writer, item); err != nil {
				return err
			}
		}
		// 116
	case []byte:
		if err := writer.WriteByte(116); err != nil {
			return err
		}
		WriteUint8Array(writer, t)
	default:
		return errors.New(fmt.Sprintf("unsupported data type : %v \n", reflect.TypeOf(t)))
	}

	return nil
}
