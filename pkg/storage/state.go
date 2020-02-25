package storage

import "github.com/cigarframework/reconciled/pkg/proto"

//go:generate go-generate -name StateMap -generator sync/map map[string]State

type State interface {
	GetMeta() *proto.Meta
}

func SparseKey(state State) (group, kind, name string) {
	meta := state.GetMeta()
	return meta.GetGroup(), meta.GetKind(), meta.GetName()
}
