package webhook

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/common"
	"github.com/cigarframework/reconciled/pkg/optional"
	"github.com/cigarframework/reconciled/pkg/plugins/grpc"
	"github.com/cigarframework/reconciled/pkg/plugins/http"
	"github.com/cigarframework/reconciled/pkg/storage"
	"go.uber.org/zap"
)

const Kind = "WebhookPlugin"

func init() {
	storage.RegisterSystemKind(Kind)
}

//go:generate syncmap -name pluginMap -pkg webhook -o pluginmap_gen.go map[string]*subPlugin

type subPlugin struct {
	name   string
	plugin api.Plugin
	spec   *Spec
}

type subPluginMatchResult struct {
	plugin *subPlugin
	scope  *ScopeSpec
}

type subPluginMatchResultSlice []*subPluginMatchResult

func (s subPluginMatchResultSlice) Len() int {
	return len(s)
}

func (s subPluginMatchResultSlice) Less(i, j int) bool {
	return s[i].plugin.name < s[j].plugin.name
}

func (s subPluginMatchResultSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type plugin struct {
	server             api.Server
	plugins            *pluginMap
	logger             *zap.Logger
	cancelSubscription context.CancelFunc
}

func New(server api.Server, logger *zap.Logger, config *Config) (api.Plugin, error) {
	p := &plugin{
		server:  server,
		logger:  logger,
		plugins: &pluginMap{},
	}
	return p, nil
}

func (p *plugin) IsPlugin() {}

func (p *plugin) Audit(ctx context.Context, review *api.ReviewRequest) error {
	plugins := p.matchPlugin(ctx, review)

	var errs []error
	for _, plugin := range plugins {
		if !plugin.plugin.spec.Audit {
			continue
		}
		if err := plugin.plugin.plugin.Audit(ctx, review); err != nil {
			errs = append(errs, err)
		}
	}

	if len(errs) > 0 {
		fields := make([]zap.Field, len(errs))
		for i, e := range errs {
			fields[i] = zap.Error(e)
		}
		p.logger.Warn("errors occurred running plugins.Audit", fields...)
	}

	return nil
}

func (p *plugin) Admission(ctx context.Context, review *api.ReviewRequest) (*api.ReviewRequest, error) {
	plugins := p.matchPlugin(ctx, review)

	var errs []error
	for _, plugin := range plugins {
		if !plugin.plugin.spec.Admission {
			continue
		}
		newRecord, err := plugin.plugin.plugin.Admission(ctx, review)
		if err != nil {
			if plugin.plugin.spec.Required {
				return nil, err
			}
			errs = append(errs, err)
		} else {
			review = newRecord
		}
	}

	if len(errs) > 0 {
		fields := make([]zap.Field, len(errs))
		for i, e := range errs {
			fields[i] = zap.Error(e)
		}
		p.logger.Warn("errors occurred running plugins.Admission", fields...)
	}

	return review, nil
}

func (p *plugin) matchPlugin(ctx context.Context, review *api.ReviewRequest) []*subPluginMatchResult {
	var result []*subPluginMatchResult

	p.plugins.Range(func(name string, value *subPlugin) bool {
		meta := review.State.GetMeta()

		var scopeFound *ScopeSpec
		for _, scope := range value.spec.Scopes {
			if scope.Group != "" && scope.Group != meta.GetGroup() {
				continue
			}

			if scope.Kind != "" && scope.Kind != meta.GetKind() {
				continue
			}

			if scope.Name != "" && scope.Name != meta.GetName() {
				continue
			}
			actionFound := false

			for _, action := range scope.Actions {
				if action == review.Action {
					actionFound = true
					break
				}
			}

			if !actionFound {
				continue
			}
			scopeFound = scope
			break
		}

		if scopeFound != nil {
			result = append(result, &subPluginMatchResult{
				plugin: value,
				scope:  scopeFound,
			})
		}
		return true
	})

	sort.Sort(subPluginMatchResultSlice(result))
	return result
}

func (p *plugin) Start(ctx context.Context) error {
	c, cancel := context.WithCancel(context.Background())

	list, ch, err := p.server.List(c, &api.ListOptions{Group: storage.SystemGroup, Kind: Kind}, &api.WatchOptions{BufferSize: 64})
	if err != nil {
		return err
	}
	p.cancelSubscription = cancel

	for _, item := range list {
		if err := p.putPlugin(item); err != nil {
			return err
		}
	}

	go p.watchPlugin(ch)
	return nil
}

func (p *plugin) Stop(ctx context.Context) error {
	p.cancelSubscription()
	return nil
}

func (p *plugin) putPlugin(state storage.State) error {
	p.plugins.Delete(state.GetMeta().GetName())
	j := state.(*storage.JSON)
	b, _ := j.Spec.MarshalJSON()

	spec := &Spec{}
	if err := json.Unmarshal(b, spec); err != nil {
		return err
	}

	if err := spec.Validate(); err != nil {
		return fmt.Errorf("%w: %s", api.ErrBadData, err.Error())
	}

	var splugin api.Plugin

	if spec.GRPC != nil {
		c := &grpc.Config{
			Audit:     spec.Audit,
			Admission: spec.Admission,
			Addr:      spec.GRPC.Addr,
		}
		if spec.GRPC.TLS != nil {
			c.TLS = &common.TLSConfig{
				Key:  spec.GRPC.TLS.Key,
				Cert: spec.GRPC.TLS.Cert,
			}
		}
		grpcPlugin, err := grpc.New(c)
		if err != nil {
			return err
		}
		splugin = grpcPlugin
	} else if spec.HTTP != nil {
		c := &http.Config{
			Audit:     spec.Audit,
			Admission: spec.Admission,
			TLS:       nil,
			Addr:      spec.HTTP.Addr,
		}
		if spec.HTTP.TLS != nil {
			c.TLS = &common.TLSConfig{
				Key:  spec.HTTP.TLS.Key,
				Cert: spec.HTTP.TLS.Cert,
			}
		}
		httpPlugin, err := http.New(c)
		if err != nil {
			return err
		}
		splugin = httpPlugin
	}

	p.plugins.Store(
		state.GetMeta().GetName(),
		&subPlugin{
			name:   state.GetMeta().GetName(),
			plugin: splugin,
			spec:   spec,
		})
	return nil
}

func (p *plugin) watchPlugin(ch <-chan *api.Notification) {
	for n := range ch {
		if n.Error != nil {
			p.logger.Error(optional.UseString(n.Error))
		} else if n.State != nil {
			if err := p.putPlugin(n.State); err != nil {
				p.logger.Error("unable to put plugin", zap.Error(err))
			}
		} else if optional.UseBool(n.IsDelete) {
			p.plugins.Delete(n.State.GetMeta().GetName())
		}
	}
}
