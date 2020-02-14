package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/common"
	"github.com/cigarframework/reconciled/pkg/storage"
)

type plugin struct {
	config *Config
	client *http.Client
}

func New(config *Config) (api.Plugin, error) {
	tls, err := common.LoadTLSConfig(config.TLS)
	if err != nil {
		return nil, err
	}

	p := &plugin{
		config: config,
		client: &http.Client{},
	}

	if tls != nil {
		p.client.Transport = &http.Transport{
			TLSClientConfig: tls,
		}
	}

	return p, nil
}

func (p *plugin) IsPlugin() {}

func (p *plugin) Audit(ctx context.Context, review *api.ReviewRequest) error {
	if !p.config.Audit {
		return nil
	}

	_, err := p.request(ctx, "/audit", review, false)
	return err
}

func (p *plugin) Admission(ctx context.Context, review *api.ReviewRequest) (*api.ReviewRequest, error) {
	if !p.config.Admission {
		return review, nil
	}

	return p.request(ctx, "/admission", review, true)
}

func (p *plugin) Start(ctx context.Context) error {
	return nil
}

func (p *plugin) Stop(ctx context.Context) error {
	return nil
}

func (p *plugin) request(ctx context.Context, path string, review *api.ReviewRequest, hasResponse bool) (*api.ReviewRequest, error) {
	addr := p.config.Addr[rand.Intn(len(p.config.Addr))]
	u, err := url.Parse(addr)
	if err != nil {
		return nil, err
	}
	u.Path = path

	b, err := json.Marshal(review)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode%100 != 2 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, errors.New(string(b))
	}

	if !hasResponse {
		return nil, nil
	}

	b, err = ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	result := &api.ReviewRequest{
		State: &storage.JSON{},
	}
	if err := json.Unmarshal(b, result); err != nil {
		return nil, err
	}
	return result, nil
}
