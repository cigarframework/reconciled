package http

import (
	"github.com/cigarframework/reconciled/pkg/common"
)

type Config struct {
	Audit     bool              `yaml:"audit"`
	Admission bool              `yaml:"admission"`
	TLS       *common.TLSConfig `yaml:"tls"`
	Addr      []string          `yaml:"addr"`
}
