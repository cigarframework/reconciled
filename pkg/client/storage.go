package client

import (
	"context"
	"time"

	"github.com/cenkalti/backoff"
	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
	"go.uber.org/zap"
)

type MemorizedClient struct {
	*Client
	backoff         *backoff.ExponentialBackOff
	internalStorage *storage.StateMap
	listOptions     *api.ListOptions
	cancelFunc      context.CancelFunc
	logger          *zap.Logger
	stopped         bool
}

func NewStorage(client *Client, listOptions *api.ListOptions, logger *zap.Logger) *MemorizedClient {
	return &MemorizedClient{
		Client:          client,
		backoff:         backoff.NewExponentialBackOff(),
		internalStorage: &storage.StateMap{},
		listOptions:     listOptions,
		logger:          logger,
	}
}

func (c *MemorizedClient) Start(ctx context.Context) error {

	return nil
}

func (c *MemorizedClient) watch() error {
	runCtx, cancelFunc := context.WithCancel(context.Background())
	c.cancelFunc = cancelFunc

	list, ch, err := c.List(runCtx, c.listOptions, &api.WatchOptions{BufferSize: 32})
	if err != nil {
		return err
	}

	for _, state := range list {
		c.internalStorage.Store(joinKey(storage.SparseKey(state)), state)
	}

	go func() {
		for n := range ch {
			if optional.UseBool(n.IsDelete) {
				c.internalStorage.Delete(joinKey(storage.SparseKey(n.State)))
				continue
			}
			if n.Error != nil {
				c.logger.Error("stream error")
				continue
			}
			c.internalStorage.Store(joinKey(storage.SparseKey(n.State)), n.State)
		}
		if !c.stopped {
			for {
				if err := c.watch(); err != nil {
					time.Sleep(c.backoff.NextBackOff())
					continue
				}
				return
			}

		}
	}()

	return nil
}

func (c *MemorizedClient) Stop(ctx context.Context) error {
	c.stopped = true
	c.cancelFunc()
	return c.Close()
}

func (c *MemorizedClient) Get(group, kind, name string) (storage.State, bool) {
	return c.internalStorage.Load(joinKey(group, kind, name))
}

func (c *MemorizedClient) Remove(ctx context.Context, group, kind, name string) error {
	c.internalStorage.Delete(joinKey(group, kind, name))
	return c.Remove(ctx, group, kind, name)
}

func joinKey(group, kind, name string) string {
	return group + "/" + kind + "/" + name
}
