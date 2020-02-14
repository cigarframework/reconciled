package grpcserver

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/cigarframework/reconciled/pkg/storage"
	"github.com/gogo/protobuf/types"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type Server struct {
	stateServer api.Server
	options     *options
	logger      *zap.Logger
}

func New(stateServer api.Server, logger *zap.Logger, opts ...optionFunc) *Server {
	opt := &options{}
	for _, o := range opts {
		opt = o(opt)
	}

	return &Server{
		stateServer: stateServer,
		options:     opt,
		logger:      logger,
	}
}

func (s *Server) Get(ctx context.Context, req *proto.Meta) (res *proto.State, err error) {
	ctx = patchContextUser(ctx)

	state, err := s.stateServer.Get(ctx, req.GetGroup(), req.GetKind(), req.GetName())
	if err != nil {
		return nil, wrapError(err)
	}

	res, err = storage.ToProto(state)
	err = wrapError(err)
	return
}

func (s *Server) Remove(ctx context.Context, req *proto.Meta) (res *types.Empty, err error) {
	ctx = patchContextUser(ctx)

	if err := s.stateServer.Remove(ctx, req.GetGroup(), req.GetKind(), req.GetName()); err != nil {
		return nil, wrapError(err)
	}

	return &types.Empty{}, nil
}

func (s *Server) List(req *proto.ListRequest, stream proto.StateService_ListServer) error {
	ctx := patchContextUser(stream.Context())

	var watch *api.WatchOptions
	if req.Watch {
		watch = &api.WatchOptions{BufferSize: s.options.streamBufferSize}
	}

	list, ch, err := s.stateServer.List(
		ctx,
		&api.ListOptions{
			Expression: req.GetExpression(),
			Group:      req.GetGroup(),
			Kind:       req.GetKind(),
			Name:       req.GetName(),
		},
		watch)
	if err != nil {
		return wrapError(err)
	}

	for _, state := range list {
		res, err := storage.ToProto(state)
		if err != nil {
			return wrapError(err)
		}
		if err := stream.Send(&proto.ListResponse{State: res}); err != nil {
			return wrapError(err)
		}
	}

	if !req.Watch {
		return nil
	}

	for n := range ch {
		if n.State != nil {
			res, err := storage.ToProto(n.State)
			if err != nil {
				return wrapError(err)
			}
			if err := stream.Send(&proto.ListResponse{
				State:    res,
				IsUpdate: optional.UseBool(n.IsUpdate),
				IsDelete: optional.UseBool(n.IsDelete),
			}); err != nil {
				return wrapError(err)
			}
		}
		if n.Error != nil {
			if err := stream.Send(&proto.ListResponse{Error: optional.UseString(n.Error)}); err != nil {
				return wrapError(err)
			}
		}
	}

	return nil
}

func (s *Server) Create(ctx context.Context, req *proto.State) (res *proto.State, err error) {
	ctx = patchContextUser(ctx)

	in, err := storage.FromProto(req)
	if err != nil {
		return nil, wrapError(err)
	}

	state, err := s.stateServer.Create(ctx, in)
	if err != nil {
		return nil, wrapError(err)
	}

	res, err = storage.ToProto(state)
	err = wrapError(err)
	return
}

func (s *Server) Update(ctx context.Context, req *proto.State) (res *proto.State, err error) {
	ctx = patchContextUser(ctx)

	in, err := storage.FromProto(req)
	if err != nil {
		return nil, wrapError(err)
	}

	state, err := s.stateServer.Update(ctx, in)
	if err != nil {
		return nil, wrapError(err)
	}

	res, err = storage.ToProto(state)
	err = wrapError(err)
	return
}

func (s *Server) Patch(ctx context.Context, req *proto.PatchRequest) (res *proto.State, err error) {
	ctx = patchContextUser(ctx)

	in := storage.NewMetaStateFromMeta(req.GetMeta())
	var patch []*api.Patch
	if err := json.Unmarshal([]byte(req.GetJSONPatch().Value), patch); err != nil {
		return nil, wrapError(fmt.Errorf("%w: %s", api.ErrBadData, err.Error()))
	}

	state, err := s.stateServer.Patch(ctx, in, patch)
	if err != nil {
		return nil, wrapError(err)
	}

	res, err = storage.ToProto(state)
	err = wrapError(err)
	return
}

func patchContextUser(ctx context.Context) context.Context {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ctx
	}

	tokens := md.Get(api.AuthHeader)
	if len(tokens) > 0 {
		return api.WithToken(ctx, tokens[0])
	}

	return ctx
}
