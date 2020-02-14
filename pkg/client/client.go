package client

import (
	"context"
	"encoding/json"
	"io"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/grpclb"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/cigarframework/reconciled/pkg/storage"
	"github.com/gogo/protobuf/types"
	"google.golang.org/grpc"
)

type Client struct {
	option *option
	client proto.StateServiceClient
	conn   *grpc.ClientConn
}

func New(addrs []string, options ...optionFunc) (*Client, error) {
	opt := &option{}
	for _, o := range options {
		opt = o(opt)
	}
	if opt.streamBufferSize <= 0 {
		opt.streamBufferSize = 1
	}

	resolver := grpclb.NewResolverBuilder(addrs)
	opt.dialOptions = append(opt.dialOptions, grpc.WithDefaultCallOptions())
	conn, err := grpc.Dial(resolver.Addr(), opt.dialOptions...)
	if err != nil {
		return nil, err
	}
	client := proto.NewStateServiceClient(conn)
	return &Client{
		client: client,
		conn:   conn,
		option: opt,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) Get(ctx context.Context, group, kind, name string) (storage.State, error) {
	ctx = c.patchContext(ctx)
	res, err := c.client.Get(ctx, &proto.Meta{
		Group: group,
		Kind:  kind,
		Name:  name,
	})
	if err != nil {
		return nil, err
	}
	return storage.FromProto(res)
}

func (c *Client) Create(ctx context.Context, state storage.State) (storage.State, error) {
	ctx = c.patchContext(ctx)
	req, err := storage.ToProto(state)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Create(ctx, req)
	if err != nil {
		return nil, err
	}
	return storage.FromProto(res)
}

func (c *Client) List(ctx context.Context, listOptions *api.ListOptions, watch *api.WatchOptions) ([]storage.State, <-chan *api.Notification, error) {
	ctx = c.patchContext(ctx)
	stream, err := c.client.List(
		ctx,
		&proto.ListRequest{
			Expression: listOptions.Expression,
			Watch:      watch != nil,
			Group:      listOptions.Group,
			Kind:       listOptions.Kind,
			Name:       listOptions.Name,
		})
	if err != nil {
		return nil, nil, err
	}

	if watch == nil {
		var list []storage.State
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				return nil, nil, err
			}
			notification, err := parseListResponse(res)
			if err != nil {
				return nil, nil, err
			}
			if !optional.UseBool(notification.IsDelete) && notification.State != nil {
				list = append(list, notification.State)
			}
		}
		return list, nil, nil
	}

	ch := make(chan *api.Notification, watch.BufferSize)
	go func() {
		defer func() {
			close(ch)
		}()
		for {
			res, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					return
				}
				ch <- api.NewErrorNotification(err.Error())
				return
			}
			notification, err := parseListResponse(res)
			if err != nil {
				ch <- api.NewErrorNotification(err.Error())
				return
			}
			ch <- notification
		}
	}()
	return nil, ch, nil
}

func (c *Client) Remove(ctx context.Context, group, kind, name string) error {
	ctx = c.patchContext(ctx)
	_, err := c.client.Remove(ctx, &proto.Meta{
		Group: group,
		Kind:  kind,
		Name:  name,
	})
	return err
}

func (c *Client) Patch(ctx context.Context, query storage.State, ops []*api.Patch) (storage.State, error) {
	ctx = c.patchContext(ctx)

	patch, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	res, err := c.client.Patch(ctx, &proto.PatchRequest{
		Meta:      query.GetMeta(),
		JSONPatch: &types.StringValue{Value: string(patch)},
	})
	if err != nil {
		return nil, err
	}

	return storage.FromProto(res)
}

func (c *Client) Update(ctx context.Context, state storage.State) (storage.State, error) {
	ctx = c.patchContext(ctx)
	req, err := storage.ToProto(state)
	if err != nil {
		return nil, err
	}
	res, err := c.client.Update(ctx, req)
	if err != nil {
		return nil, err
	}
	return storage.FromProto(res)
}
