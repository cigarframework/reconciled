package server

import (
	"context"
	"errors"
	"fmt"

	"github.com/antonmedv/expr"
	"github.com/antonmedv/expr/vm"
	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
)

func (s *Server) List(ctx context.Context, listOptions *api.ListOptions, watch *api.WatchOptions) ([]storage.State, <-chan *api.Notification, error) {
	record, err := s.runPlugins(
		ctx,
		&api.ReviewRequest{
			Action:     api.ListAction,
			Expression: listOptions.Expression,
			State:      storage.NewMetaState(listOptions.Group, listOptions.Kind, listOptions.Name),
		})
	if err != nil {
		return nil, nil, err
	}
	listOptions.Expression = record.Expression
	listOptions.Group = record.State.GetMeta().GetGroup()
	listOptions.Kind = record.State.GetMeta().GetKind()
	listOptions.Name = record.State.GetMeta().GetName()

	var matcher *vm.Program
	if listOptions.Expression != "" {
		var err error
		matcher, err = expr.Compile(listOptions.Expression)
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
		}
	}

	list := make([]storage.State, 0)
	var evaluateError error
	var evaluateResult bool
	s.storage.Range(func(_, _, _ string, state storage.State) bool {
		evaluateResult, evaluateError = ifStateMatch(matcher, listOptions, state)
		if evaluateError != nil {
			return false
		}
		if evaluateResult {
			list = append(list, state)
		}
		return true
	})

	if evaluateError != nil {
		return nil, nil, fmt.Errorf("%w: %s", api.ErrBadData, evaluateError.Error())
	}

	if watch == nil {
		if matcher != nil {
			matcher.Disassemble()
		}
		return list, nil, nil
	}

	ch := make(chan *api.Notification, watch.BufferSize)
	intermediateCh := make(chan *api.Notification, watch.BufferSize)
	s.subscription.Subscribe(intermediateCh)
	go func(ctx context.Context, options *api.ListOptions, matcher *vm.Program, out chan *api.Notification, in chan *api.Notification) {
		defer func() {
			s.subscription.Cancel(in)
			if matcher != nil {
				matcher.Disassemble()
			}
			close(in)
			close(out)
		}()
		for {
			select {
			case <-ctx.Done():
				if err := ctx.Err(); err != nil {
					out <- &api.Notification{Error: optional.Error(err)}
				}
				return
			case n := <-in:
				if n.State != nil {
					res, err := ifStateMatch(matcher, options, n.State)
					if err != nil {
						out <- &api.Notification{Error: optional.Error(err)}
						continue
					}
					if res {
						out <- n
					}
					continue
				}
				out <- n
			}
		}

	}(ctx, listOptions, matcher, ch, intermediateCh)
	return list, ch, nil
}

func ifStateMatch(matcher *vm.Program, options *api.ListOptions, state storage.State) (bool, error) {
	if options.Group != "" && state.GetMeta().GetGroup() != options.Group {
		return false, nil
	}

	if options.Kind != "" && state.GetMeta().GetKind() != options.Kind {
		return false, nil
	}

	if options.Name != "" && state.GetMeta().GetName() != options.Name {
		return false, nil
	}

	if matcher == nil {
		return true, nil
	}

	res, err := expr.Run(matcher, state)
	if err != nil {
		return false, err
	}

	b, ok := res.(bool)
	if !ok {
		return false, errors.New("expression should return a boolean result")
	}
	return b, nil
}
