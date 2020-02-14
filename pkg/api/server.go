package api

import (
	"context"

	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
)

type WatchOptions struct {
	BufferSize int
}

type ListOptions struct {
	Expression string
	Group      string
	Kind       string
	Name       string
}

type Notification struct {
	State    storage.State `json:"State,omitempty"`
	IsDelete *bool         `json:"IsDelete,omitempty"`
	IsUpdate *bool         `json:"IsUpdate,omitempty"`
	Error    *string       `json:"Error,omitempty"`
}

func NewErrorNotification(err string) *Notification {
	return &Notification{
		Error: optional.String(err),
	}
}

type Server interface {
	Get(ctx context.Context, group, kind, name string) (storage.State, error)
	Create(ctx context.Context, state storage.State) (storage.State, error)
	List(ctx context.Context, listOptions *ListOptions, watch *WatchOptions) ([]storage.State, <-chan *Notification, error)
	Remove(ctx context.Context, group, kind, name string) error
	Patch(ctx context.Context, query storage.State, ops []*Patch) (storage.State, error)
	Update(ctx context.Context, state storage.State) (storage.State, error)
}
