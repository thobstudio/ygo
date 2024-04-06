package ynotgo

import (
	"math/rand"

	"github.com/google/uuid"
)

type Doc struct {
	gc                  bool
	guid                string
	clientId            uint32
	collectionId        string
	autoLoad            bool
	shouldLoad          bool
	share               map[string]DocIntegrator
	store               *StructStore
}

type Option func(*Doc)

func newDoc(options ...Option) *Doc {
	doc := &Doc{
		autoLoad:            false,
		shouldLoad:          true,
		gc:                  true,
		guid:                uuid.NewString(),
		clientId:            rand.Uint32(),
		share:               make(map[string]DocIntegrator, 0),
		store:               newStructStore(),
	}

	for _, o := range options {
		o(doc)
	}

	return doc
}

func NewDoc(options ...Option) *Doc {
	return newDoc(options...)
}

func WithGc(gc bool) Option {
	return func(doc *Doc) {
		doc.gc = gc
	}
}

func WithGuid(guid string) Option {
	return func(doc *Doc) {
		doc.guid = guid
	}
}

func WithClientId(clientId uint32) Option {
	return func(doc *Doc) {
		doc.clientId = clientId
	}
}

func WithCollectionId(collectionId string) Option {
	return func(doc *Doc) {
		doc.collectionId = collectionId
	}
}

func WithAutoLoad(autoLoad bool) Option {
	return func(doc *Doc) {
		doc.autoLoad = autoLoad
	}
}

func WithShouldLoad(shouldLoad bool) Option {
	return func(doc *Doc) {
		doc.shouldLoad = shouldLoad
	}
}

func (doc *Doc) GetMap(name string) *YMap {
	if m, ok := doc.share[name].(*YMap); ok {
		return m
	}

	m := newYMap()
	m.integrate(doc)
	doc.share[name] = m

	return m
}
