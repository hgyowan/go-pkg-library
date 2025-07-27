package grpc

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/oklog/run"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"net"
	"strings"
	"syscall"
	"time"
)

const (
	DefaultMaxMsgSize = 100 * 1024 * 1024
)

type GrpcServer interface {
	RegisterService(desc *grpc.ServiceDesc, impl any)
	Serve(ctx context.Context, port string)
}

type server struct {
	srv *grpc.Server
}

func MustNewGRPCServer() GrpcServer {
	s := &server{}

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			CustomErrorInterceptor,
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
		),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             time.Minute,
			PermitWithoutStream: false,
		}),
		grpc.MaxRecvMsgSize(DefaultMaxMsgSize),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge:      3 * time.Minute,
			MaxConnectionAgeGrace: 30 * time.Second,
		}),
	)

	s.srv = srv

	return s
}

func (s *server) Serve(ctx context.Context, port string) {
	grpc_health_v1.RegisterHealthServer(s.srv, &server{})
	g := &run.Group{}

	if port == "" {
		port = envs.ServerPort
	}

	if !strings.HasPrefix(port, ":") {
		port = ":" + port
	}

	g.Add(func() error {
		l, err := net.Listen("tcp", port)
		if err != nil {
			return err
		}
		return s.srv.Serve(l)

	}, func(err error) {
		s.srv.GracefulStop()
		s.srv.Stop()
	})

	g.Add(run.SignalHandler(ctx, syscall.SIGINT, syscall.SIGTERM))
	if err := g.Run(); err != nil {
		pkgLogger.ZapLogger.Logger.Info("error serving grpc")
	}
}

func (s *server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.srv.RegisterService(desc, impl)
}

func (s *server) Check(ctx context.Context, req *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (s *server) Watch(*grpc_health_v1.HealthCheckRequest, grpc_health_v1.Health_WatchServer) error {
	return nil
}

func (s *server) List(ctx context.Context, request *grpc_health_v1.HealthListRequest) (*grpc_health_v1.HealthListResponse, error) {
	return nil, nil
}
