package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
	uuid "github.com/satori/go.uuid"
)

func (s *Server) Create(ctx context.Context, state storage.State) (storage.State, error) {
	record, err := s.runPlugins(ctx, &api.ReviewRequest{Action: api.CreateAction, State: state})
	if err != nil {
		return nil, err
	}
	state = record.State

	group, kind, name := storage.SparseKey(state)
	key := joinKeys(group, kind, name)
	if name == "" {
		name = uuid.NewV4().String()
	}

	if err := storage.Validate(state); err != nil {
		return nil, api.ErrBadData
	}

	_, ok := s.storage.Load(group, kind, name)
	if ok {
		return nil, fmt.Errorf("%w: state %s exist", api.ErrExist, key)
	}

	state.GetMeta().Name = name
	state.GetMeta().CreatedAt = optional.Now()
	state.GetMeta().Version = 0
	state.GetMeta().UpdatedAt = optional.Now()

	body, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	if _, err := s.etcdClient.Put(ctx, s.options.etcdPrefix+key, string(body)); err != nil {
		return nil, errors.New("unable to save data to etcd")
	}

	s.storage.Store(state)
	return state, nil
}
