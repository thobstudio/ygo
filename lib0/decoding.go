package lib0

import (
	"bufio"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

func ReadVarUint(reader *bufio.Reader) (uint64, error) {
	return binary.ReadUvarint(reader)
}

func ReadVarInt(reader *bufio.Reader) (int32, error) {
	r, err := reader.ReadByte()
	if err != nil {
		return 0, err
	}
	num := uint32(r) & Bits6
	mult := uint32(64)
	var sign int32 = 1
	if r&byte(Bit7) > 0 {
		sign = -1
	}
	if r&byte(Bit8) == 0 {
		return sign * int32(num), nil
	}

	for {
		r, err := reader.ReadByte()
		if err == io.EOF {
			return sign * int32(num), nil
		}
		if err != nil {
			return 0, nil
		}

		num = num + uint32(r&byte(Bits7))*mult
		mult *= 128

		if r < byte(Bit8) {
			return sign * int32(num), nil
		}
	}
}

func ReadUint8Array(reader *bufio.Reader) ([]byte, error) {
	length, err := ReadVarUint(reader)
	if err != nil {
		return []byte{}, err
	}
	buf := make([]byte, length)
	_, err = reader.Read(buf)
	if err != nil {
		return []byte{}, nil
	}
	return buf, nil
}

func ReadVarString(reader *bufio.Reader) (string, error) {
	length, err := ReadVarUint(reader)
	if err != nil {
		return "", err
	}
	buf := make([]byte, length)
	_, err = reader.Read(buf)
	if err != nil {
		return "", nil
	}
	return string(buf), nil
}

func ReadFloat32(reader *bufio.Reader) (float32, error) {
	var f float32
	if err := binary.Read(reader, binary.BigEndian, &f); err != nil {
		return 0, err
	}
	return f, nil
}

func ReadFloat64(reader *bufio.Reader) (float64, error) {
	var f float64
	if err := binary.Read(reader, binary.BigEndian, &f); err != nil {
		return 0, err
	}
	return f, nil
}

func ReadAny(reader *bufio.Reader) (any, error) {
	b, err := reader.ReadByte()
	if err != nil {
		return nil, err
	}
	switch b {
	case 127, 126:
		return nil, nil
	case 125:
		return ReadVarInt(reader)
	case 124:
		return ReadFloat32(reader)
	case 123:
		return ReadFloat64(reader)
	case 121:
		return false, nil
	case 120:
		return true, nil
	case 119:
		return ReadVarString(reader)
	case 118: // map[string]any
		length, err := ReadVarUint(reader)
		if err != nil {
			return nil, err
		}
		m := make(map[string]any, length)
		for i := 0; i < int(length); i++ {
			s, err := ReadVarString(reader)
			if err != nil {
				return nil, err
			}
			v, err := ReadAny(reader)
			if err != nil {
				return nil, err
			}
			m[s] = v
		}
		return m, nil
	case 117:
		length, err := ReadVarUint(reader)
		if err != nil {
			return nil, err
		}
		a := make([]any, length)
		for i := 0; i < int(length); i++ {
			v, err := ReadAny(reader)
			if err != nil {
				return nil, err
			}
			a[i] = v
		}
		return a, nil
	case 116:
		return ReadUint8Array(reader)
	default:
		return nil, errors.New(fmt.Sprintf("unsupported data type : %v ", b))
	}
}
