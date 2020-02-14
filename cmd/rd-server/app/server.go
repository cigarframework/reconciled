package app

import (
	"net"
	"net/http"

	"github.com/cigarframework/reconciled/pkg/api"
	"github.com/cigarframework/reconciled/pkg/plugins"

	"github.com/cigarframework/reconciled/pkg/common"
	"github.com/cigarframework/reconciled/pkg/grpcserver"
	"github.com/cigarframework/reconciled/pkg/httpserver"
	"github.com/cigarframework/reconciled/pkg/proto"
	"github.com/cigarframework/reconciled/pkg/server"
	etcdv3 "go.etcd.io/etcd/clientv3"
	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

type Server struct {
	config       *Config
	logger       *zap.Logger
	stateServer  *server.Server
	httpServer   *http.Server
	grpcServer   *grpc.Server
	grpcListener net.Listener
	plugins      []api.Plugin
}

func New(logger *zap.Logger, config *Config) (*Server, error) {
	etcdTLS, err := common.LoadTLSConfig(config.ETCD.TLS)
	if err != nil {
		return nil, err
	}

	httpTLS, err := common.LoadTLSConfig(config.HTTP.TLS)
	if err != nil {
		return nil, err
	}

	grpcTLS, err := common.LoadTLSConfig(config.GRPC.TLS)
	if err != nil {
		return nil, err
	}

	etcdClient, err := etcdv3.New(etcdv3.Config{
		Endpoints: config.ETCD.Endpoints,
		TLS:       etcdTLS,
		Username:  config.ETCD.Username,
		Password:  config.ETCD.Password,
	})
	if err != nil {
		return nil, err
	}

	stateServer := server.New(server.WithETCD(etcdClient), server.WithLogger(logger), server.WithETCDPrefix(config.ETCD.Prefix))
	httpHandler := httpserver.New(stateServer, logger, httpserver.WithStreamBufferSize(config.BufferSize))
	grpcHandler := grpcserver.New(stateServer, logger, grpcserver.WithStreamBufferSize(config.BufferSize))

	httpServer := &http.Server{
		Addr:      config.HTTP.Addr,
		Handler:   httpHandler,
		TLSConfig: httpTLS,
		ErrorLog:  zap.NewStdLog(logger),
	}

	var grpcOptions []grpc.ServerOption
	if grpcTLS != nil {
		grpcOptions = append(grpcOptions, grpc.Creds(credentials.NewTLS(grpcTLS)))
	}
	grpcServer := grpc.NewServer(grpcOptions...)
	proto.RegisterStateServiceServer(grpcServer, grpcHandler)

	var apiPlugins []api.Plugin
	for _, c := range config.Plugins {
		plugin, err := plugins.Construct(c, logger, stateServer)
		if err != nil {
			return nil, err
		}
		apiPlugins = append(apiPlugins, plugin)
	}

	return &Server{
		config:      config,
		logger:      logger,
		stateServer: stateServer,
		httpServer:  httpServer,
		grpcServer:  grpcServer,
		plugins:     apiPlugins,
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	if err := s.stateServer.Start(s.plugins, ctx); err != nil {
		return err
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			s.logger.Fatal(err.Error())
		}
	}()

	l, err := net.Listen("tcp", s.config.GRPC.Addr)
	if err != nil {
		return err
	}
	s.grpcListener = l

	go func() {
		if err := s.grpcServer.Serve(l); err != nil {
			s.logger.Fatal(err.Error())
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	s.grpcServer.GracefulStop()
	return common.ReduceErrors(
		s.grpcListener.Close(),
		s.httpServer.Close(),
		s.stateServer.Stop(ctx),
	)

}
