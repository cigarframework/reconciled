package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJSON(t *testing.T) {
	data := `{"Meta":{"Group":"group","Kind":"kind","Name":"name"},"Spec":{"message":"hello world"}}`

	state := &JSON{}
	err := json.Unmarshal([]byte(data), state)
	assert.Nil(t, err)

	assert.Equal(t, "group", state.GetMeta().GetGroup())
	assert.Equal(t, "kind", state.GetMeta().GetKind())
	assert.Equal(t, "name", state.GetMeta().GetName())

	spec, err := state.Spec.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, `{"message":"hello world"}`, string(spec))
}
