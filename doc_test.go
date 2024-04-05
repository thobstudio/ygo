package ynotgo

import (
	"math/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestDocDefaultOptions(t *testing.T) {
	doc := NewDoc()

	assert.True(t, doc.gc)
	assert.True(t, doc.shouldLoad)
	assert.False(t, doc.autoLoad)
	assert.NotEmpty(t, doc.guid)
	assert.GreaterOrEqual(t, doc.clientId, uint32(0))
	assert.Empty(t, doc.collectionId)
}

func TestDocOptions(t *testing.T) {
	guid := uuid.NewString()
	clientId := rand.Uint32()
	doc := NewDoc(
		WithGuid(guid),
		WithClientId(clientId),
		WithAutoLoad(true),
		WithShouldLoad(false),
		WithCollectionId("dummy"),
		WithGc(false),
	)

	assert.False(t, doc.gc)
	assert.False(t, doc.shouldLoad)
	assert.True(t, doc.autoLoad)
	assert.Equal(t, doc.guid, guid)
	assert.Equal(t, doc.clientId, clientId)
	assert.Equal(t, doc.collectionId, "dummy")
}

func TestDocGetMap(t *testing.T) {
	doc := newDoc()

	p := doc.GetMap("project")
	assert.Equal(t, doc.share["project"], p)
	assert.Equal(t, p.doc, doc)
}
