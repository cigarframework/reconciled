package storage

import "github.com/cigarframework/reconciled/pkg/proto"

//go:generate syncmap -name stateMap -pkg storage -o statemap_gen.go map[string]State

type State interface {
	GetMeta() *proto.Meta
}

func SparseKey(state State) (group, kind, name string) {
	meta := state.GetMeta()
	return meta.GetGroup(), meta.GetKind(), meta.GetName()
}
