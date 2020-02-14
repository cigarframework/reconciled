package server

import (
	"context"
	"encoding/json"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/common"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
	"github.com/coreos/etcd/mvcc/mvccpb"
	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

type Server struct {
	storage      *storage.Storage
	subscription SubscriptionManager
	options      *options
	etcdClient   *etcdv3.Client
	logger       *zap.Logger
	plugins      []api.Plugin
}

func New(opts ...optionFunc) *Server {
	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		panic(err)
	}

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	opt := &options{
		etcdClient: etcdClient,
		logger:     logger,
	}
	for _, o := range opts {
		opt = o(opt)
	}

	return &Server{
		storage:      storage.New(),
		subscription: newSubscriptionManager(1024),
		options:      opt,
		etcdClient:   opt.etcdClient,
		logger:       opt.logger,
	}
}

func (s *Server) Start(plugins []api.Plugin, ctx context.Context) error {
	s.plugins = plugins
	list, err := s.etcdClient.Get(ctx, s.options.etcdPrefix, etcdv3.WithPrefix())
	if err != nil {
		return err
	}

	for _, kv := range list.Kvs {
		if state := s.unmarshalKV(kv); state != nil {
			s.storage.Store(state)
		}
	}

	for _, p := range s.plugins {
		if err := p.Start(ctx); err != nil {
			return err
		}
	}

	go s.watch()
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	var errors []error
	for _, p := range s.plugins {
		errors = append(errors, p.Stop(ctx))
	}

	return common.ReduceErrors(
		s.logger.Sync(),
		s.etcdClient.Close(),
		common.ReduceErrors(errors...),
	)
}

func (s *Server) watch() {
	watcher := s.etcdClient.Watch(context.Background(), s.options.etcdPrefix, etcdv3.WithPrefix())
	for res := range watcher {
		if err := res.Err(); err != nil {
			s.logger.Error("unable to watch", zap.Error(err))
			continue
		}

		for _, ev := range res.Events {
			key := string(ev.Kv.Key)[len(s.options.etcdPrefix):]
			group, kind, name := parseKey(key)
			switch ev.Type {
			case etcdv3.EventTypePut:
				if state := s.unmarshalKV(ev.Kv); state != nil {
					if ev.IsModify() {
						oldState, ok := s.storage.Load(group, kind, name)

						if ok && oldState.GetMeta().GetVersion() > state.GetMeta().GetVersion() {
							continue
						}

						s.subscription.Publish(&api.Notification{
							State:    state,
							IsUpdate: optional.True(),
						})
					} else {
						s.subscription.Publish(&api.Notification{
							State: state,
						})
					}
					s.storage.Store(state)
				}
			case etcdv3.EventTypeDelete:
				{
					s.subscription.Publish(&api.Notification{
						State:    storage.NewMetaState(group, kind, name),
						IsDelete: optional.True(),
					})
					s.storage.Delete(group, kind, name)
				}
			}
		}
	}
}

func (s *Server) unmarshalKV(kv *mvccpb.KeyValue) storage.State {
	state := &storage.JSON{}
	if err := json.Unmarshal(kv.Value, state); err != nil {
		s.logger.Error("unable to unmarshal etcd data")
		return nil
	}

	if err := storage.Validate(state); err != nil {
		s.logger.Error("data " + string(kv.Key) + ":" + err.Error())
		return nil
	}

	return state
}

func (s *Server) runPlugins(ctx context.Context, review *api.ReviewRequest) (*api.ReviewRequest, error) {
	user := api.ContextUser(ctx)
	review.User = user
	var err error
	action := review.Action
	for _, p := range s.plugins {
		if err = p.Audit(ctx, review); err != nil {
			return nil, err
		}

		if review, err = p.Admission(ctx, review); err != nil {
			return nil, err
		}
		review.Action = action
	}

	return review, nil
}
