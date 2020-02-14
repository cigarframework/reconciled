package storage

import (
	"encoding/json"

	"github.com/cigarframework/reconciled/pkg/proto"
)

type JSON struct {
	Meta *proto.Meta     `json:"Meta,omitempty"`
	Spec json.RawMessage `json:"Spec,omitempty"`
}

func (s *JSON) GetMeta() *proto.Meta {
	if s == nil {
		return nil
	}
	return s.Meta
}
