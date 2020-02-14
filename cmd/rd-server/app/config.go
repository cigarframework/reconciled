package app

import (
	"github.com/cigarframework/reconciled/pkg/common"
	"github.com/cigarframework/reconciled/pkg/plugins"
)

type Config struct {
	ETCD       ETCDConfig        `yaml:"etcd"`
	HTTP       HTTPConfig        `yaml:"http"`
	GRPC       GRPCConfig        `yaml:"grpc"`
	BufferSize int               `yaml:"bufferSize"`
	Plugins    []*plugins.Config `yaml:"plugins"`
}

type ETCDConfig struct {
	Prefix    string            `yaml:"prefix"`
	Endpoints []string          `yaml:"endpoints"`
	Username  string            `yaml:"username"`
	Password  string            `yaml:"password"`
	TLS       *common.TLSConfig `yaml:"tls"`
}

type HTTPConfig struct {
	Addr string `yaml:"addr"`
	TLS  *common.TLSConfig
}

type GRPCConfig struct {
	Addr string `yaml:"addr"`
	TLS  *common.TLSConfig
}
