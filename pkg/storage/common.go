package storage

import "github.com/cigarframework/reconciled/pkg/proto"

type metaState struct {
	meta *proto.Meta
}

func (m *metaState) GetMeta() *proto.Meta {
	return m.meta
}

func NewMetaState(group, kind, name string) State {
	return &metaState{
		meta: NewMeta(group, kind, name),
	}
}

func NewMetaStateFromMeta(meta *proto.Meta) State {
	return &metaState{
		meta: meta,
	}
}

func NewMeta(group, kind, name string) *proto.Meta {
	return &proto.Meta{
		Group: group,
		Kind:  kind,
		Name:  name,
	}
}
