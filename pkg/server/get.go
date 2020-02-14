package server

import (
	"context"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/storage"
)

func (s *Server) Get(ctx context.Context, group, kind, name string) (storage.State, error) {
	record, err := s.runPlugins(ctx, &api.ReviewRequest{Action: api.GetAction, State: storage.NewMetaState(group, kind, name)})
	if err != nil {
		return nil, err
	}
	group, kind, name = storage.SparseKey(record.State)

	if err := storage.Validate(storage.NewMetaState(group, kind, name)); err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	state, ok := s.storage.Load(group, kind, name)
	if !ok {
		return nil, fmt.Errorf("%w: state %s does not exist", api.ErrNotExist, joinKeys(group, kind, name))
	}
	return state, nil
}
