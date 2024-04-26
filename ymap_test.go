package ynotgo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestYMapSetExistingValue(t *testing.T) {
	doc := newDoc()

	p := doc.GetMap("project")
	p.Set("boolean", true)
	b := p.Get("boolean")
	assert.True(t, b.(bool))
	p.Set("boolean", false)
	b = p.Get("boolean")
	assert.False(t, b.(bool))

	p.Set("int", 10000)
	assert.Equal(t, p.Get("int"), 10000)
	p.Set("int", -10000)
	assert.Equal(t, p.Get("int"), -10000)

	p.Set("bytearray", []byte("this is string"))
	assert.Equal(t, string(p.Get("bytearray").([]byte)), "this is string")

	p.Set("assets", newYMap())
	assets := p.Get("assets").(*YMap)
	assets.Set("url", "http://localhost:5001")
	assert.IsType(t, p.Get("assets"), &YMap{})
	assert.Equal(t, assets.Get("url"), "http://localhost:5001")
}
