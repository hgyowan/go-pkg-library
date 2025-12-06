package grpc

import (
	"context"
	"errors"
	"net"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/hgyowan/go-pkg-library/envs"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"github.com/oklog/run"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	ghealth "google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

const (
	DefaultMaxMsgSize = 100 * 1024 * 1024
)

type GrpcServer interface {
	RegisterService(desc *grpc.ServiceDesc, impl any)
	Serve(ctx context.Context, port string)
	ShutdownHealth()
}

type server struct {
	srv    *grpc.Server
	health *ghealth.Server
}

func MustNewGRPCServer() GrpcServer {
	hs := ghealth.NewServer()
	hs.SetServingStatus("*", grpc_health_v1.HealthCheckResponse_SERVING)

	srv := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(),
			CustomErrorUnaryInterceptor,
		),
		grpc.ChainStreamInterceptor(
			recovery.StreamServerInterceptor(),
			CustomErrorStreamInterceptor,
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

	return &server{srv: srv, health: hs}
}

func (s *server) Serve(ctx context.Context, port string) {
	grpc_health_v1.RegisterHealthServer(s.srv, s.health)
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

			s.health.SetServingStatus("*", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
			time.Sleep(300 * time.Millisecond)

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

func (s *server) ShutdownHealth() {
	s.health.SetServingStatus("*", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
}
