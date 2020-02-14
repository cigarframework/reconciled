package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
)

func (s *Server) Update(ctx context.Context, state storage.State) (storage.State, error) {
	record, err := s.runPlugins(ctx, &api.ReviewRequest{Action: api.UpdateAction, State: state})
	if err != nil {
		return nil, err
	}
	state = record.State

	group, kind, name := storage.SparseKey(state)
	key := joinKeys(group, kind, name)
	if err := storage.Validate(state); err != nil {
		return nil, api.ErrBadData
	}

	oldState, ok := s.storage.Load(group, kind, name)
	if !ok {
		return nil, fmt.Errorf("%w: state %s does not exist", api.ErrNotExist, key)
	}
	state.GetMeta().Version = oldState.GetMeta().Version + 1
	state.GetMeta().UpdatedAt = optional.Now()
	state.GetMeta().CreatedAt = oldState.GetMeta().GetCreatedAt()

	body, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	if _, err := s.etcdClient.Put(ctx, s.options.etcdPrefix+key, string(body)); err != nil {
		return nil, errors.New(" unable to save data to etcd")
	}

	s.storage.Store(state)
	return state, nil
}
