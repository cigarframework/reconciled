package plugins

import (
	"errors"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/plugins/authn"
	"github.com/cigarframework/reconciled/pkg/plugins/authz"
	"github.com/cigarframework/reconciled/pkg/plugins/grpc"
	"github.com/cigarframework/reconciled/pkg/plugins/http"
	"github.com/cigarframework/reconciled/pkg/plugins/webhook"
	"go.uber.org/zap"
)

type Config struct {
	Authn   *authn.Config   `yaml:"authn"`
	Authz   *authz.Config   `yaml:"authz"`
	GRPC    *grpc.Config    `yaml:"grpc"`
	HTTP    *http.Config    `yaml:"http"`
	Webhook *webhook.Config `yaml:"webhook"`
}

func Construct(config *Config, logger *zap.Logger, server api.Server) (api.Plugin, error) {
	if c := config.Authn; c != nil {
		return authn.New(server, logger, c)
	}

	if c := config.Authz; c != nil {
		return authz.New(server, logger, c)
	}

	if c := config.GRPC; c != nil {
		return grpc.New(c)
	}

	if c := config.HTTP; c != nil {
		return http.New(c)
	}

	if c := config.Webhook; c != nil {
		return webhook.New(server, logger, c)
	}

	return nil, errors.New("no plugin specified")
}
