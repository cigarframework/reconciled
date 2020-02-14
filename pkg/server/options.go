package server

import (
	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
)

type options struct {
	logger     *zap.Logger
	etcdPrefix string
	etcdClient *etcdv3.Client
}

type optionFunc func(options *options) *options

func WithETCDPrefix(prefix string) optionFunc {
	return func(options *options) *options {
		options.etcdPrefix = prefix
		return options
	}
}

func WithETCD(client *etcdv3.Client) optionFunc {
	return func(options *options) *options {
		options.etcdClient = client
		return options
	}
}

func WithLogger(logger *zap.Logger) optionFunc {
	return func(options *options) *options {
		options.logger = logger
		return options
	}
}
