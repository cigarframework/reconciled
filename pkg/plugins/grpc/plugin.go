package grpc

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/common"
	"github.com/cigarframework/reconciled/pkg/grpclb"
	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/cigarframework/reconciled/pkg/storage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/resolver"
)

type plugin struct {
	resolver *grpclb.ResolverBuilder
	config   *Config
	client   proto.PluginClient
	conn     *grpc.ClientConn
}

func New(config *Config) (api.Plugin, error) {
	tls, err := common.LoadTLSConfig(config.TLS)
	if err != nil {
		return nil, err
	}

	options := make([]grpc.DialOption, 0, len(config.DialOptions))
	copy(options, config.DialOptions)
	p := &plugin{
		config: &Config{
			Audit:       config.Audit,
			Admission:   config.Admission,
			TLS:         config.TLS,
			Addr:        config.Addr,
			DialOptions: options,
		},
	}
	p.resolver = grpclb.NewResolverBuilder(config.Addr)
	if tls != nil {
		config.DialOptions = append(config.DialOptions, grpc.WithTransportCredentials(credentials.NewTLS(tls)))
	} else {
		config.DialOptions = append(config.DialOptions, grpc.WithInsecure())
	}

	config.DialOptions = append(config.DialOptions, grpc.WithBlock())
	resolver.Register(p.resolver)

	return p, nil
}

func (p *plugin) IsPlugin() {}

func (p *plugin) Audit(ctx context.Context, review *api.ReviewRequest) error {
	if !p.config.Audit {
		return nil
	}

	req, opts, err := prepareRequest(ctx, review)
	if err != nil {
		return err
	}

	_, err = p.client.Audit(ctx, req, opts...)
	return err
}

func (p *plugin) Admission(ctx context.Context, review *api.ReviewRequest) (*api.ReviewRequest, error) {
	if !p.config.Admission {
		return review, nil
	}

	req, opts, err := prepareRequest(ctx, review)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Admission(ctx, req, opts...)
	if err != nil {
		return nil, err
	}

	action := api.Action(res.Action)
	if !action.Valid() {
		return nil, fmt.Errorf("%w: invalid action", api.ErrBadData)
	}

	state, err := storage.FromProto(res.State)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	var patch []*api.Patch
	if v := res.GetPatch(); v != "" {
		if err := json.Unmarshal([]byte(v), &patch); err != nil {
			return nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
		}
	}

	var user *api.User
	if u := res.GetUser(); u != nil {
		user = &api.User{
			Name:  u.GetName(),
			Group: u.GetGroup(),
			Token: u.GetToken(),
		}
	}

	return &api.ReviewRequest{
		Action:     action,
		State:      state,
		Patch:      patch,
		User:       user,
		Expression: res.GetExpression(),
	}, nil
}

func (p *plugin) Start(ctx context.Context) error {
	options := []grpc.DialOption{
		grpc.WithDefaultServiceConfig(fmt.Sprintf(`{"LoadBalancingPolicy": "%s"}`, roundrobin.Name)),
	}

	conn, err := grpc.DialContext(ctx, p.resolver.Addr(), append(options, p.config.DialOptions...)...)
	if err != nil {
		return err
	}
	client := proto.NewPluginClient(conn)
	p.client = client
	p.conn = conn

	return nil
}

func (p *plugin) Stop(ctx context.Context) error {
	return p.conn.Close()
}

func prepareRequest(ctx context.Context, review *api.ReviewRequest) (*proto.ReviewRequest, []grpc.CallOption, error) {
	var options []grpc.CallOption

	state, err := storage.ToProto(review.State)
	if err != nil {
		return nil, nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	var patch []byte
	if review.Patch != nil {
		var err error
		patch, err = json.Marshal(&review.Patch)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
		}
	}

	var user *proto.User
	if u := review.User; u != nil {
		user = &proto.User{
			Name:  u.Name,
			Group: u.Group,
			Token: u.Token,
		}
	}

	return &proto.ReviewRequest{
		Action:     review.Action.String(),
		State:      state,
		Patch:      string(patch),
		User:       user,
		Expression: review.Expression,
	}, options, nil
}
