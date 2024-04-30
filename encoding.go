package ynotgo

import (
	"bufio"
	"bytes"
	"io"

	"github.com/thobstudio/ynotgo/lib0"
)

func readStateVectorUpdate(update []byte) (map[uint32]uint32, error) {
	return readStateVector(bufio.NewReader(bytes.NewBuffer(update)))
}

func readStateVector(reader *bufio.Reader) (map[uint32]uint32, error) {
	sslength, err := lib0.ReadVarUint(reader)
	if err != nil {
		if err == io.EOF {
			return make(map[uint32]uint32, 0), nil
		}
		return nil, err
	}
	stateVector := make(map[uint32]uint32, sslength)
	for i := 0; i < int(sslength); i++ {
		client, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, err
		}
		clock, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, err
		}
		stateVector[client] = clock
	}

	return stateVector, nil
}
