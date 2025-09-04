package grpc

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/oklog/run"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"os/signal"
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
			CustomErrorUnaryInterceptor,
			recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			CustomErrorStreamInterceptor,
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
			return pkgError.Wrap(errors.New("failed to listen: " + err.Error()))
		}
		return s.srv.Serve(l)

	}, func(err error) {
		done := make(chan struct{})
		go func() {
			s.srv.GracefulStop()
			close(done)
		}()

		select {
		case <-done:
			pkgLogger.ZapLogger.Logger.Info("gRPC server stopped gracefully")
		case <-time.After(10 * time.Second):
			pkgLogger.ZapLogger.Logger.Warn("Graceful stop timed out, forcing stop")
			s.srv.Stop()
		}
	})

	g.Add(func() error {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case sig := <-sigCh:
			pkgLogger.ZapLogger.Logger.Info("caught signal", zap.String("signal", sig.String()))
			return nil
		}
	}, func(err error) {
	})

	if err := g.Run(); err != nil {
		pkgLogger.ZapLogger.Logger.Error("gRPC server exited with error", zap.Error(err))
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
