package storage

import "github.com/cigarframework/reconciled/pkg/proto"

type testState struct {
	Meta *proto.Meta
	Spec map[string]interface{}
}

func (s *testState) GetMeta() *proto.Meta {
	return s.Meta
}

func newTestState() State {
	return &testState{
		Meta: &proto.Meta{
			Group: "group",
			Kind:  "kind",
			Name:  "name",
		},
		Spec: map[string]interface{}{
			"string": "string",
			"number": 123,
			"bool":   true,
			"array":  []string{"a", "b", "c"},
			"map": map[string]interface{}{
				"child": true,
			},
		},
	}
}
