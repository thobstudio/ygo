package ynotgo

import (
	"bufio"
	"errors"

	"github.com/thobstudio/ynotgo/lib0"
)

const (
	MessageSyncStep1 = 0
	MessageSyncStep2 = 1
	MessageUpdate    = 2
)

func WriteSyncStep1(rw *bufio.ReadWriter, doc *Doc) error {
	if err := lib0.WriteVarUint(rw.Writer, MessageSyncStep1); err != err {
		return err
	}
	stateVector, err := doc.EncodeStateVectorV1()
	if err != nil {
		return err
	}
	if err := lib0.WriteUint8Array(rw.Writer, stateVector); err != nil {
		return err
	}
	return nil
}

func WriteSyncStep2(rw *bufio.ReadWriter, doc *Doc, encodedStateVector []byte) error {
	if err := lib0.WriteVarUint(rw.Writer, MessageSyncStep2); err != nil {
		return err
	}
	buf, err := doc.EncodeStateAsUpdateV1(encodedStateVector)
	if err != nil {
		return err
	}
	if err := lib0.WriteUint8Array(rw.Writer, buf); err != nil {
		return err
	}

	return nil
}

func ReadSyncStep1(rw *bufio.ReadWriter, doc *Doc) error {
	buf, err := lib0.ReadUint8Array(rw.Reader)
	if err != nil {
		return err
	}
	return WriteSyncStep2(rw, doc, buf)
}

func ReadSyncStep2(rw *bufio.ReadWriter, doc *Doc, txOrigin any) error {
	update, err := lib0.ReadUint8Array(rw.Reader)
	if err != nil {
		return err
	}
	return doc.ApplyUpdateV1(update, txOrigin)
}

func ReadUpdate(rw *bufio.ReadWriter, doc *Doc, txOrigin any) error {
	return ReadSyncStep2(rw, doc, txOrigin)
}

func WriteUpdate(rw *bufio.ReadWriter, update []byte) error {
	if err := lib0.WriteVarUint(rw.Writer, MessageUpdate); err != nil {
		return err
	}
	if err := lib0.WriteUint8Array(rw.Writer, update); err != nil {
		return err
	}
	return nil
}

func ReadSyncMessage(rw *bufio.ReadWriter, doc *Doc, txOrigin any) (uint32, error) {
	messageType, err := lib0.ReadVarUint(rw.Reader)
	if err != nil {
		return 0, err
	}
	switch messageType {
	case MessageSyncStep1:
		return messageType, ReadSyncStep1(rw, doc)
	case MessageSyncStep2:
		return messageType, ReadSyncStep2(rw, doc, txOrigin)
	case MessageUpdate:
		return messageType, ReadUpdate(rw, doc, txOrigin)
	default:
		return messageType, errors.New("unknown message type")
	}
}
