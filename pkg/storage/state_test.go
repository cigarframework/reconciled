package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestState(t *testing.T) {
	group, kind, name := SparseKey(newTestState())
	assert.Equal(t, "group", group)
	assert.Equal(t, "kind", kind)
	assert.Equal(t, "name", name)
}
