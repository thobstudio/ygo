package ynotgo

import (
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDocDefaultOptions(t *testing.T) {
	doc := NewDoc()

	assert.True(t, doc.opts.gc)
	assert.True(t, doc.opts.shouldLoad)
	assert.False(t, doc.opts.autoLoad)
	assert.NotEmpty(t, doc.opts.guid)
	assert.GreaterOrEqual(t, doc.opts.clientId, uint32(0))
	assert.Empty(t, doc.opts.collectionId)
}

func TestDocOptions(t *testing.T) {
	guid := uuid.NewString()
	clientId := rand.Uint32()
	doc := NewDocWithOptions(
		WithGuid(guid),
		WithClientId(clientId),
		WithAutoLoad(true),
		WithShouldLoad(false),
		WithCollectionId("dummy"),
		WithGc(false),
	)

	assert.False(t, doc.opts.gc)
	assert.False(t, doc.opts.shouldLoad)
	assert.True(t, doc.opts.autoLoad)
	assert.Equal(t, doc.opts.guid, guid)
	assert.Equal(t, doc.opts.clientId, clientId)
	assert.Equal(t, doc.opts.collectionId, "dummy")
}

func TestDocGetMap(t *testing.T) {
	doc := newDoc()

	p := doc.GetMap("project")
	assert.Equal(t, doc.share["project"], p)
	assert.Equal(t, p.doc, doc)
}
