package caseio

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderInt(t *testing.T) {
	r := Example("1 2 3 4 5")
	assert.Equal(t, 1, r.Int())
	assert.Equal(t, 2, r.Int())
	assert.Equal(t, 3, r.Int())
	assert.Equal(t, 4, r.Int())
	assert.Equal(t, 5, r.Int())

	r = Example("1\n2\n3 4\n5")
	assert.Equal(t, 1, r.Int())
	assert.Equal(t, 2, r.Int())
	assert.Equal(t, 3, r.Int())
	assert.Equal(t, 4, r.Int())
	assert.Equal(t, 5, r.Int())
}

func TestReaderString(t *testing.T) {
	r := Example("123 456\n789\nabc")
	assert.Equal(t, "123", r.String())
	assert.Equal(t, "456", r.String())
	assert.Equal(t, "789", r.String())
	assert.Equal(t, "abc", r.String())
}
