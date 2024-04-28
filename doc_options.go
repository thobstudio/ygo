package ynotgo

import (
	"math/rand"

	"github.com/google/uuid"
)

func generateNewClientId() uint32 {
	return rand.Uint32()
}

type DocOptions struct {
	gc           bool
	guid         string
	clientId     uint32
	collectionId string
	autoLoad     bool
	shouldLoad   bool
	meta         any
}

type Option func(*DocOptions)

func newDocOptions(options ...Option) *DocOptions {
	opts := &DocOptions{
		autoLoad:   false,
		shouldLoad: true,
		gc:         true,
		guid:       uuid.NewString(),
		clientId:   generateNewClientId(),
	}
	for _, o := range options {
		o(opts)
	}

	return opts
}

func WithGc(gc bool) Option {
	return func(doc *DocOptions) {
		doc.gc = gc
	}
}

func WithGuid(guid string) Option {
	return func(doc *DocOptions) {
		doc.guid = guid
	}
}

func WithClientId(clientId uint32) Option {
	return func(doc *DocOptions) {
		doc.clientId = clientId
	}
}

func WithCollectionId(collectionId string) Option {
	return func(doc *DocOptions) {
		doc.collectionId = collectionId
	}
}

func WithAutoLoad(autoLoad bool) Option {
	return func(doc *DocOptions) {
		doc.autoLoad = autoLoad
	}
}

func WithShouldLoad(shouldLoad bool) Option {
	return func(doc *DocOptions) {
		doc.shouldLoad = shouldLoad
	}
}
