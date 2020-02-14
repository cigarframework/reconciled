package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
	jsonpatch "github.com/evanphx/json-patch"
)

func (s *Server) Patch(ctx context.Context, query storage.State, ops []*api.Patch) (storage.State, error) {
	record, err := s.runPlugins(ctx, &api.ReviewRequest{Action: api.PatchAction, State: query, Patch: ops})
	if err != nil {
		return nil, err
	}
	query = record.State

	group, kind, name := storage.SparseKey(query)
	key := joinKeys(group, kind, name)
	if err := storage.Validate(query); err != nil {
		return nil, api.ErrBadData
	}

	state, ok := s.storage.Load(group, kind, name)
	if !ok {
		return nil, fmt.Errorf("%w: state %s does not exist", api.ErrNotExist, key)
	}
	state.GetMeta().Version += 1
	state.GetMeta().UpdatedAt = optional.Now()

	body, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	jsonState := &storage.JSON{}
	if err := json.Unmarshal(body, jsonState); err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	specBytes, err := jsonState.Spec.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	patchBytes, err := json.Marshal(ops)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	patch, err := jsonpatch.DecodePatch(patchBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	specBytes, err = patch.Apply(specBytes)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	if err := jsonState.Spec.UnmarshalJSON(specBytes); err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	body, err = json.Marshal(jsonState)
	if err != nil {
		return nil, err
	}

	if _, err := s.etcdClient.Put(ctx, s.options.etcdPrefix+key, string(body)); err != nil {
		return nil, errors.New("unable to save data to etcd")
	}

	s.storage.Store(jsonState)
	return jsonState, nil
}
