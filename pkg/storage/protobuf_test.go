package storage

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

func TestGRPCToProto(t *testing.T) {
	state, err := ToProto(newTestState())
	assert.Nil(t, err)
	assert.Equal(t, "group", state.GetMeta().GetGroup())
	assert.Equal(t, "kind", state.GetMeta().GetKind())
	assert.Equal(t, "name", state.GetMeta().GetName())

	assert.Equal(t, "string", state.GetSpec().GetFields()["string"].GetStringValue())
	assert.Equal(t, float64(123), state.GetSpec().GetFields()["number"].GetNumberValue())
	assert.Equal(t, true, state.GetSpec().GetFields()["bool"].GetBoolValue())

	list := state.GetSpec().GetFields()["array"].GetListValue().Values
	assert.Equal(t, "a", list[0].GetStringValue())
	assert.Equal(t, "b", list[1].GetStringValue())
	assert.Equal(t, "c", list[2].GetStringValue())

	assert.Equal(t, true, state.GetSpec().GetFields()["map"].GetStructValue().GetFields()["child"].GetBoolValue())
}

func TestGRPCFromProto(t *testing.T) {
	jsonState, err := ToProto(newTestState())
	assert.Nil(t, err)

	state, err := FromProto(jsonState)
	assert.Nil(t, err)

	body, err := json.Marshal(state)
	assert.Nil(t, err)

	assert.Equal(t, "group", state.GetMeta().GetGroup())
	assert.Equal(t, "kind", state.GetMeta().GetKind())
	assert.Equal(t, "name", state.GetMeta().GetName())

	assert.Equal(t, "string", gjson.Get(string(body), "Spec.string").String())
	assert.Equal(t, float64(123), gjson.Get(string(body), "Spec.number").Float())
	assert.Equal(t, true, gjson.Get(string(body), "Spec.bool").Bool())

	assert.Equal(t, "a", gjson.Get(string(body), "Spec.array.0").String())
	assert.Equal(t, "b", gjson.Get(string(body), "Spec.array.1").String())
	assert.Equal(t, "c", gjson.Get(string(body), "Spec.array.2").String())

	assert.Equal(t, true, gjson.Get(string(body), "Spec.map.child").Bool())
}
