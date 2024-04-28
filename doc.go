package ynotgo

import (
	"bufio"
	"bytes"
	"cmp"
	"errors"
	"fmt"
	"io"
	"slices"

	"github.com/olebedev/emitter"
	"github.com/thobstudio/ynotgo/lib0"
)

type Doc struct {
	emitter.Emitter
	opts                *DocOptions
	share               map[string]SharedType
	store               *StructStore
	item                *Item
	transaction         *Transaction
	transactionCleanups []*Transaction
}

func newDoc() *Doc {
	return newDocWithOptions()
}

func NewDoc() *Doc {
	return newDoc()
}

func newDocWithOptions(options ...Option) *Doc {
	opts := newDocOptions(options...)
	doc := &Doc{
		opts:                opts,
		share:               make(map[string]SharedType),
		store:               newStructStore(),
		transactionCleanups: make([]*Transaction, 0),
	}
	doc.Use("*", emitter.Void)

	return doc
}

func NewDocWithOptions(options ...Option) *Doc {
	return newDocWithOptions(options...)
}

func (doc *Doc) Get(name string) SharedType {
	if m, ok := doc.share[name].(SharedType); ok {
		return m
	}

	t := newAbstractType()
	t.Integrate(doc, nil)
	doc.share[name] = t

	return t
}

func (doc *Doc) GetMap(name string) (*YMap, error) {
	var t *YMap
	val, ok := doc.share[name]

	if ok {
		if m, ok := val.(*YMap); ok {
			t = m
		} else if at, ok := val.(*AbstractType); ok {
			t = newYMap()
			t.length = at.length
			t.itemMap = at.itemMap
			t.ForEachItem(func(item *Item, key string) {
				n := item
				for n != nil {
					n.parent = t
					n = n.left
				}
			})
			t.start = at.start
			n := t.start

			for n != nil {
				n.parent = t
				n = n.right
			}
		} else {
			return nil, errors.New(fmt.Sprintf("type with name %s has already been defined with a different struct", name))
		}
	} else {
		t = newYMap()
	}

	doc.share[name] = t
	t.Integrate(doc, nil)
	return t, nil
}

func (doc *Doc) Transact(handler TransactionHandler, origin any, local bool) (any, error) {
	transactionCleanups := doc.transactionCleanups
	initialCall := false
	if doc.transaction == nil {
		initialCall = true
		doc.transaction = newTrasaction(doc, origin, local)
		transactionCleanups = append(transactionCleanups, doc.transaction)

		if len(doc.transactionCleanups) == 1 {
			doc.Emit("beforeAllTransactions", doc)
		}

		doc.Emit("beforeTransaction", doc.transaction)
	}

	result, err := handler(doc.transaction)

	if initialCall {
		finishCleaup := transactionCleanups[0] == doc.transaction
		doc.transaction = nil
		if finishCleaup {
			cleanupTransactions(transactionCleanups, 0)
		}
	}

	return result, err
}

func (doc *Doc) readClientsStructRefs(decoder UpdateDecoder) (map[uint32]*ClientStructRef, error) {
	reader := decoder.Reader()

	clientRefs := make(map[uint32]*ClientStructRef)
	numOfStateUpdates, err := lib0.ReadVarUint(reader)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(numOfStateUpdates); i++ {
		numOfStructs, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, err
		}

		refs := make([]SharedStruct, numOfStructs)

		client, err := decoder.ReadClient()
		if err != nil {
			return nil, err
		}
		clock, err := lib0.ReadVarUint(reader)
		if err != nil {
			return nil, err
		}

		clientRefs[client] = &ClientStructRef{i: 0, refs: refs}

		for i := 0; i < int(numOfStructs); i++ {
			info, err := decoder.ReadInfo()
			if err != nil {
				return nil, err
			}
			switch lib0.Bits5 & info {
			case 0:
				length, err := decoder.ReadLen()
				if err != nil {
					return nil, err
				}
				refs[i] = newGc(newId(client, clock), length)
				clock += length
				break
			case 10:
				panic("not implemented")
			default:
				var (
					leftOrigin         *ID
					rightOrigin        *ID
					hasParentYKey      bool   = false
					cantCopyParentInfo bool   = (info & (lib0.Bit7 | lib0.Bit8)) == 0
					parentYKey         string = ""
					parent             any
					parentSub          string
				)

				if info&lib0.Bit8 == lib0.Bit8 {
					leftOrigin, err = decoder.ReadLeftId()
					if err != nil {
						return nil, err
					}
				}

				if info&lib0.Bit7 == lib0.Bit7 {
					rightOrigin, err = decoder.ReadRightId()
					if err != nil {
						return nil, err
					}
				}

				if cantCopyParentInfo {
					hasParentYKey, err = decoder.ReadParentInfo()
				}

				if cantCopyParentInfo && hasParentYKey {
					parentYKey, err = decoder.ReadString()
				}

				if cantCopyParentInfo && !hasParentYKey {
					parent, err = decoder.ReadLeftId()
					if err != nil {
						return nil, err
					}
				} else if parentYKey != "" {
					parent = doc.Get(parentYKey)
				}

				if cantCopyParentInfo && (info&lib0.Bit6) == lib0.Bit6 {
					parentSub, err = decoder.ReadString()
					if err != nil {
						return nil, err
					}
				}
				content, err := readItemContent(decoder, info)
				if err != nil {
					return nil, err
				}
				s := newItem(newId(client, clock), nil, nil, leftOrigin, rightOrigin, parent, parentSub, content)
				refs[i] = s
				clock += s.length
			}
		}
	}

	return clientRefs, nil
}

/*
* Encode State As Update
 */

func (doc *Doc) encodeStateAsUpdate(encoder UpdateEncoder, encodedTargetStateVector []byte) error {
	targetStateVector, err := readStateVector(bufio.NewReader(bytes.NewBuffer(encodedTargetStateVector)))
	if err != nil {
		return err
	}

	if err := doc.store.WriteClientStructs(encoder, targetStateVector); err != nil {
		return err
	}

	if err := doc.store.WriteDeleteSet(encoder); err != nil {
		return err
	}

	return nil
}

func (doc *Doc) EncodeStateAsUpdateV1(encodedTargetStateVector []byte) ([]byte, error) {
	encoder := NewUpdateEncoderV1()
	if err := doc.encodeStateAsUpdate(encoder, encodedTargetStateVector); err != nil {
		return []byte{}, err
	}
	return encoder.ToUint8Array()
}

/*
* Encoder State Vector
 */

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

func (doc *Doc) writeStateVector(encoder DsEncoder, sv StateVector) error {
	if err := lib0.WriteVarUint(encoder.Writer(), uint32(len(sv))); err != nil {
		return err
	}

	clients := make([]uint32, len(sv))
	i := 0
	for client := range sv {
		clients[i] = client
		i++
	}
	slices.SortFunc(clients, func(a, b uint32) int {
		return cmp.Compare(b, a)
	})

	for _, client := range clients {
		if err := lib0.WriteVarUint(encoder.Writer(), client); err != nil {
			return err
		}
		if err := lib0.WriteVarUint(encoder.Writer(), sv[client]); err != nil {
			return err
		}
	}

	return nil
}

func (doc *Doc) encodeStateVector(dsencoder DsEncoder) ([]byte, error) {
	sv := doc.store.StateVector()

	if err := doc.writeStateVector(dsencoder, sv); err != nil {
		return nil, err
	}

	return dsencoder.ToUint8Array()
}

func (doc *Doc) EncodeStateVectorV1() ([]byte, error) {
	dsencoder := newDsEncoderV1()
	return doc.encodeStateVector(dsencoder)
}

/*
* Apply Update
 */

func (doc *Doc) applyUpdate(decoder UpdateDecoder, txOrigin any) error {
	_, err := doc.Transact(func(tx *Transaction) (any, error) {
		retry := false
		ss, err := doc.readClientsStructRefs(decoder)
		if err != nil {
			return nil, err
		}
		restStructs, err := doc.store.integrateStructs(ss, tx)
		if err != nil {
			return nil, err
		}
		pending := doc.store.pendingStructs
		if pending != nil {
			for client, clock := range pending.missing {
				if clock < doc.store.State(client) {
					retry = true
					break
				}
			}

			if restStructs != nil {
				for client, clock := range restStructs.missing {
					if mclock, ok := pending.missing[client]; ok && mclock > clock {
						pending.missing[client] = clock
					}
				}
				update, err := mergeUpdatesV2([][]byte{pending.update, restStructs.update})
				if err != nil {
					return nil, err
				}
				pending.update = update
			}
		} else {
			doc.store.pendingStructs = restStructs
		}

		dsRest, err := doc.store.readAndApplyDeleteSet(decoder, tx)
		if err != nil {
			return nil, err
		}
		if doc.store.pendingDs != nil {
			pendingDsDecoder := newUpdateDecoderV2(doc.store.pendingDs)
			lib0.ReadVarUint(pendingDsDecoder.Reader())
			dsRest2, err := doc.store.readAndApplyDeleteSet(pendingDsDecoder, tx)
			if err != nil {
				return nil, err
			}
			if dsRest != nil && dsRest2 != nil {
				mergedUpdates, err := mergeUpdatesV2([][]byte{dsRest, dsRest2})
				if err != nil {
					return nil, err
				}
				doc.store.pendingDs = mergedUpdates
			} else {
				if dsRest != nil {
					doc.store.pendingDs = dsRest
				} else {
					doc.store.pendingDs = dsRest2
				}
			}
		} else {
			doc.store.pendingDs = dsRest
		}

		if retry {
			update := doc.store.pendingStructs.update
			doc.store.pendingStructs = nil
			if err := doc.ApplyUpdateV2(update, txOrigin); err != nil {
				return nil, err
			}
		}
		return nil, nil
	}, txOrigin, false)
	return err
}

func (doc *Doc) ApplyUpdateV1(update []byte, txOrigin any) error {
	decoder := newUpdateDecoderV1(update)
	return doc.applyUpdate(decoder, txOrigin)
}

func (doc *Doc) ApplyUpdateV2(update []byte, txOrigin any) error {
	panic("not implemented")
}
