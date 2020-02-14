package authn

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/tidwall/buntdb"
	"go.uber.org/zap"
)

type plugin struct {
	server api.Server
	auth   AuthMethod
	db     *buntdb.DB
	logger *zap.Logger
	config *Config
}

func New(server api.Server, logger *zap.Logger, config *Config) (api.Plugin, error) {
	p := &plugin{
		server: server,
		logger: logger,
		config: config,
	}
	if config.JWT != nil {
		p.auth = &jwtAuth{key: config.JWT.Key}
	} else {
		return nil, errors.New("no valid auth method")
	}
	return p, nil
}

func (p *plugin) IsPlugin() {}

func (p *plugin) Audit(context.Context, *api.ReviewRequest) error {
	return nil
}

func (p *plugin) Admission(ctx context.Context, review *api.ReviewRequest) (*api.ReviewRequest, error) {
	meta := review.State.GetMeta()
	if meta == nil {
		return nil, api.ErrUnAuthenticated
	}

	user := review.User
	isAuthenticated := user != nil
	if isAuthenticated && (user.Group == "" || user.Name == "") {
		isAuthenticated = false
	}

	if isAuthenticated {
		return review, nil
	}

	if user.Token == "" {
		return nil, api.ErrUnAuthenticated
	}

	if err := p.db.View(func(tx *buntdb.Tx) error {
		res, err := tx.Get(user.Token)
		if err == buntdb.ErrNotFound {
			return nil
		}
		u := &api.User{}
		if err := json.Unmarshal([]byte(res), u); err != nil {
			return err
		}

		review.User = u
		return nil
	}); err == nil {
		return review, nil
	}

	user, err := p.auth.Auth(user.Token)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", api.ErrUnAuthenticated, err.Error())
	}

	b, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}

	p.db.Update(func(tx *buntdb.Tx) error {
		_, _, err := tx.Set(user.Token, string(b), &buntdb.SetOptions{TTL: p.config.CacheDuration})
		return err
	})

	review.User = user
	return review, nil
}

func (p *plugin) Start(ctx context.Context) error {
	var err error
	p.db, err = buntdb.Open(":memory:")
	return err
}

func (p *plugin) Stop(ctx context.Context) error {
	return p.db.Close()
}
