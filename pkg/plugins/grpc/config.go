package grpc

import (
	"github.com/cigarframework/reconciled/pkg/common"
	"google.golang.org/grpc"
)

type Config struct {
	Audit       bool              `yaml:"audit"`
	Admission   bool              `yaml:"admission"`
	TLS         *common.TLSConfig `yaml:"tls"`
	Addr        []string          `yaml:"addr"`
	DialOptions []grpc.DialOption `yaml:"-"`
}
