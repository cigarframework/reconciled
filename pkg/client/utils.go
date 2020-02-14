package client

import (
	"context"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/cigarframework/reconciled/pkg/storage"
	"google.golang.org/grpc/metadata"
)

func (c *Client) patchContext(ctx context.Context) context.Context {
	if c.option.getToken == nil {
		return ctx
	}
	return metadata.AppendToOutgoingContext(ctx, api.AuthHeader, c.option.getToken())
}

func parseListResponse(res *proto.ListResponse) (*api.Notification, error) {
	n := &api.Notification{}

	if res.GetIsDelete() {
		n.IsDelete = optional.True()
	}

	if res.GetIsUpdate() {
		n.IsUpdate = optional.True()
	}

	if res.GetError() != "" {
		n.Error = optional.String(res.GetError())
	}

	if res.GetState() != nil {
		state, err := storage.FromProto(res.GetState())
		if err != nil {
			return nil, err
		}
		n.State = state
	}

	return n, nil
}
