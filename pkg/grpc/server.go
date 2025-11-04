package grpc

import (
	"context"
	"fmt"
	"net"

	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Config struct {
	Address string
}

type Server struct {
	config Config
	logger *zap.Logger
	server *grpc.Server
}

func (s *Server) RegisterService(sd *grpc.ServiceDesc, ss any) {
	s.server.RegisterService(sd, ss)
}

func (s *Server) OnStart(_ context.Context) error {
	lis, err := net.Listen("tcp4", s.config.Address)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	go func(listener net.Listener) {
		s.logger.Info("grpc server listening on", zap.String("address", s.config.Address))

		err := s.server.Serve(listener)
		if err != nil {
			s.logger.Error("failed to serve", zap.Error(err))
		}
	}(lis)

	return nil
}

func (s *Server) OnStop(_ context.Context) error {
	s.server.GracefulStop()

	return nil
}

func NewServer(c Config, log *zap.Logger) *Server {
	return &Server{
		config: c,
		server: grpc.NewServer(grpc.ChainUnaryInterceptor(grpc_zap.UnaryServerInterceptor(log))),
		logger: log,
	}
}
