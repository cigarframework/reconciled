package storage

import (
	"encoding/json"

	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/gogo/protobuf/types"
)

func ToProto(state State) (*proto.State, error) {
	if state == nil {
		return nil, nil
	}
	g := &Protobuf{}
	b, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(b, g); err != nil {
		return nil, err
	}
	return g.State, nil
}

func FromProto(state *proto.State) (State, error) {
	if state == nil {
		return nil, nil
	}

	g := &Protobuf{state}
	b, err := json.Marshal(g)
	if err != nil {
		return nil, err
	}

	j := &JSON{}
	if err := json.Unmarshal(b, j); err != nil {
		return nil, err
	}
	return j, nil
}

type Protobuf struct {
	*proto.State
}

type protobufState struct {
	Meta *proto.Meta
	Spec map[string]interface{}
}

func (s *Protobuf) UnmarshalJSON(b []byte) error {
	g := &protobufState{}

	if err := json.Unmarshal(b, g); err != nil {
		return err
	}

	spec := &types.Struct{
		Fields: map[string]*types.Value{},
	}

	for k, v := range g.Spec {
		spec.Fields[k] = unmarshalStructValue(v)
	}

	st := &proto.State{
		Meta: g.Meta,
		Spec: spec,
	}
	s.State = st
	return nil
}

func unmarshalStructValue(in interface{}) *types.Value {
	if in == nil {
		return &types.Value{Kind: &types.Value_NullValue{NullValue: types.NullValue_NULL_VALUE}}
	}
	switch in.(type) {
	case string:
		return &types.Value{Kind: &types.Value_StringValue{StringValue: in.(string)}}
	case float64:
		return &types.Value{Kind: &types.Value_NumberValue{NumberValue: in.(float64)}}
	case bool:
		return &types.Value{Kind: &types.Value_BoolValue{BoolValue: in.(bool)}}
	case []interface{}:
		{
			jsonList := in.([]interface{})
			list := make([]*types.Value, len(jsonList))
			for i, val := range jsonList {
				list[i] = unmarshalStructValue(val)
			}
			return &types.Value{Kind: &types.Value_ListValue{ListValue: &types.ListValue{Values: list}}}
		}
	case map[string]interface{}:
		{
			jsonMap := in.(map[string]interface{})
			fields := make(map[string]*types.Value, len(jsonMap))
			for k, v := range jsonMap {
				fields[k] = unmarshalStructValue(v)
			}
			return &types.Value{Kind: &types.Value_StructValue{StructValue: &types.Struct{Fields: fields}}}
		}
	}
	return nil
}

func (s *Protobuf) MarshalJSON() ([]byte, error) {
	g := &protobufState{
		Meta: s.Meta,
		Spec: map[string]interface{}{},
	}
	for name, value := range s.GetSpec().GetFields() {
		g.Spec[name] = marshalStructValue(value)
	}
	return json.Marshal(g)
}

func marshalStructValue(value *types.Value) interface{} {
	switch value.Kind.(type) {
	case *types.Value_NumberValue:
		return value.GetNumberValue()
	case *types.Value_StringValue:
		return value.GetStringValue()
	case *types.Value_BoolValue:
		return value.GetBoolValue()
	case *types.Value_ListValue:
		{
			list := value.GetListValue().GetValues()
			jsonList := make([]interface{}, len(list))
			for i, val := range list {
				jsonList[i] = marshalStructValue(val)
			}
			return jsonList
		}
	case *types.Value_StructValue:
		{
			fields := value.GetStructValue().GetFields()
			jsonMap := make(map[string]interface{}, len(fields))
			for name, val := range fields {
				jsonMap[name] = marshalStructValue(val)
			}
			return jsonMap
		}
	default:
		return nil
	}
}
