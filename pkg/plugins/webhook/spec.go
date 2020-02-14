package webhook

import (
	"errors"
	"fmt"

	"github.com/cigarframework/reconciled/pkg/api"
)

type TLSSpec struct {
	Key  string `json:"key,omitempty"`
	Cert string `json:"cert,omitempty"`
}

type GRPCSpec struct {
	TLS  *TLSSpec `json:"tls,omitempty"`
	Addr []string `yaml:"addr"`
}

type HTTPSpec struct {
	TLS  *TLSSpec `json:"tls,omitempty"`
	Addr []string `json:"addr,omitempty"`
}

type ScopeSpec struct {
	Group   string       `json:"group,omitempty"`
	Kind    string       `json:"kind,omitempty"`
	Name    string       `json:"name,omitempty"`
	Actions []api.Action `json:"actions,omitempty"`
}

type Spec struct {
	Audit     bool         `json:"audit,omitempty"`
	Admission bool         `json:"admission,omitempty"`
	Required  bool         `json:"required,omitempty"`
	Scopes    []*ScopeSpec `json:"scopes,omitempty"`
	GRPC      *GRPCSpec    `json:"grpc,omitempty"`
	HTTP      *HTTPSpec    `json:"http,omitempty"`
}

func (s *Spec) Validate() error {
	if (s.GRPC == nil && s.HTTP == nil) || (s.GRPC != nil && s.HTTP != nil) {
		return errors.New("specify either grpc or http")
	}

	if !s.Admission && !s.Audit {
		return errors.New("not audit nor admission plugin")
	}

	for _, scope := range s.Scopes {
		for _, a := range scope.Actions {
			if !a.Valid() {
				return fmt.Errorf("action \"%s\" is not valid", a.String())
			}
		}
	}

	return nil
}
