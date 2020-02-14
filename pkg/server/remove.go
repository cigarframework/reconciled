package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/storage"
	etcdv3 "go.etcd.io/etcd/clientv3"
)

func (s *Server) Remove(ctx context.Context, group, kind, name string) error {
	record, err := s.runPlugins(ctx, &api.ReviewRequest{Action: api.RemoveAction, State: storage.NewMetaState(group, kind, name)})
	if err != nil {
		return err
	}
	group, kind, name = storage.SparseKey(record.State)

	if group == "" {
		return fmt.Errorf("%w: group is required", api.ErrBadData)
	}

	if kind == "" {
		if _, err := s.etcdClient.Delete(ctx, s.options.etcdPrefix+group, etcdv3.WithPrefix()); err != nil {
			return errors.New("unable to delete data from etcd")
		}
		s.storage.DeleteGroup(group)
		return nil
	}

	if name == "" {
		if _, err := s.etcdClient.Delete(ctx, s.options.etcdPrefix+joinKeys(group, kind), etcdv3.WithPrefix()); err != nil {
			return errors.New("unable to delete data from etcd")
		}
		s.storage.DeleteKind(group, kind)
		return nil
	}

	if _, err := s.etcdClient.Delete(ctx, s.options.etcdPrefix+joinKeys(group, kind, name), etcdv3.WithPrefix()); err != nil {
		return errors.New("unable to delete data from etcd")
	}
	s.storage.Delete(group, kind, name)
	return nil
}
