package storage

import (
	"testing"

	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/stretchr/testify/assert"
)

func TestCommonNewMetaState(t *testing.T) {
	state := NewMetaState("group", "kind", "name")
	assert.Equal(t, "group", state.GetMeta().GetGroup())
	assert.Equal(t, "kind", state.GetMeta().GetKind())
	assert.Equal(t, "name", state.GetMeta().GetName())
}

func TestCommonNewMetaStateFromMeta(t *testing.T) {
	state := NewMetaStateFromMeta(&proto.Meta{
		Group: "group",
		Kind:  "kind",
		Name:  "name",
	})
	assert.Equal(t, "group", state.GetMeta().GetGroup())
	assert.Equal(t, "kind", state.GetMeta().GetKind())
	assert.Equal(t, "name", state.GetMeta().GetName())
}

func TestCommonNewMeta(t *testing.T) {
	meta := NewMeta("group", "kind", "name")
	assert.Equal(t, "group", meta.GetGroup())
	assert.Equal(t, "kind", meta.GetKind())
	assert.Equal(t, "name", meta.GetName())
}
