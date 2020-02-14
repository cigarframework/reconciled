package authz

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/storage"
	"github.com/ory/ladon"
	manager "github.com/ory/ladon/manager/memory"
	"go.uber.org/zap"
)

const Kind = "Policy"

func init() {
	storage.RegisterSystemKind(Kind)
}

type plugin struct {
	server             api.Server
	logger             *zap.Logger
	warden             *ladon.Ladon
	cancelSubscription context.CancelFunc
}

func New(server api.Server, logger *zap.Logger, config *Config) (api.Plugin, error) {
	p := &plugin{
		server: server,
		logger: logger,
	}
	return p, nil
}

func (p *plugin) IsPlugin() {}

func (p *plugin) Audit(context.Context, *api.ReviewRequest) error {
	return nil
}

func (p *plugin) Admission(ctx context.Context, review *api.ReviewRequest) (*api.ReviewRequest, error) {
	if review.User == nil || review.User.Group == "" || review.User.Name == "" {
		return nil, api.ErrUnAuthenticated
	}

	req := &ladon.Request{
		Action:  review.Action.String(),
		Subject: review.User.Group + "/" + review.User.Name,
		Context: map[string]interface{}{},
	}

	req.Resource = fmt.Sprintf("%s/%s/%s", review.State.GetMeta().GetGroup(), review.State.GetMeta().GetKind(), review.State.GetMeta().GetName())
	req.Context["Meta"] = map[string]interface{}{
		"Group": review.State.GetMeta().GetGroup(),
		"Kind":  review.State.GetMeta().GetKind(),
		"Name":  review.State.GetMeta().GetName(),
	}

	switch review.Action {
	case api.ListAction:
		req.Context["Expression"] = review.Expression
	case api.CreateAction:
		fallthrough
	case api.UpdateAction:
		{
			state := review.State.(*storage.JSON)
			spec := map[string]interface{}{}
			specBytes, err := state.Spec.MarshalJSON()
			if err != nil {
				return nil, err
			}
			if err := json.Unmarshal(specBytes, spec); err != nil {
				return nil, err
			}
			req.Context["Spec"] = spec
		}
	case api.PatchAction:
		{
			var patch []map[string]interface{}
			bytes, _ := json.Marshal(review.Patch)
			if err := json.Unmarshal(bytes, &patch); err != nil {
				return nil, err
			}
			req.Context["Patch"] = patch
		}
	default:
	}

	if err := p.warden.IsAllowed(req); err != nil {
		return nil, err
	}

	return review, nil
}

func (p *plugin) Start(ctx context.Context) error {
	c, cancel := context.WithCancel(context.Background())
	warden := &ladon.Ladon{
		Manager: manager.NewMemoryManager(),
	}

	list, ch, err := p.server.List(c, &api.ListOptions{Group: storage.SystemGroup, Kind: Kind}, &api.WatchOptions{BufferSize: 64})
	if err != nil {
		return err
	}
	p.warden = warden
	p.cancelSubscription = cancel

	for _, item := range list {
		if err := p.putPolicy(item); err != nil {
			return err
		}
	}

	go p.watchPolicy(ch)
	return nil
}

func (p *plugin) Stop(ctx context.Context) error {
	p.cancelSubscription()
	return nil
}

func (p *plugin) putPolicy(state storage.State) error {
	j := state.(*storage.JSON)
	b, _ := j.Spec.MarshalJSON()

	policy := &ladon.DefaultPolicy{}
	if err := json.Unmarshal(b, policy); err != nil {
		return err
	}

	policy.ID = state.GetMeta().GetName()
	if err := p.warden.Manager.Update(policy); err != nil {
		return err
	}
	return nil
}

func (p *plugin) watchPolicy(ch <-chan *api.Notification) {
	for n := range ch {
		if n.Error != nil {
			p.logger.Error(optional.UseString(n.Error))
		} else if n.State != nil {
			if err := p.putPolicy(n.State); err != nil {
				p.logger.Error("unable to put policy", zap.Error(err))
			}
		} else if optional.UseBool(n.IsDelete) {
			if err := p.warden.Manager.Delete(n.State.GetMeta().GetName()); err != nil {
				p.logger.Error("unable to delete policy", zap.Error(err))
			}
		}
	}
}
