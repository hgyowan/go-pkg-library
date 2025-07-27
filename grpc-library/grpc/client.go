package grpc

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/retry"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/timeout"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

func MustNewGRPCClient(address string) *grpc.ClientConn {
	conn, err := grpc.NewClient(address,
		grpc.WithChainUnaryInterceptor(
			timeout.UnaryClientInterceptor(30*time.Second),
			withGrpcRetryBackoffInterceptor(),
		),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(100*1024*1024),
			grpc.MaxCallSendMsgSize(100*1024*1024),
		),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                2 * time.Minute,
			Timeout:             3 * time.Second,
			PermitWithoutStream: true,
		}),
	)

	if err != nil {
		panic(err)
	}
	return conn
}

func withGrpcRetryBackoffInterceptor() grpc.UnaryClientInterceptor {
	opts := []retry.CallOption{
		retry.WithBackoff(retry.BackoffExponentialWithJitter(100*time.Millisecond, 0.5)),
		retry.WithCodes(codes.Aborted, codes.Unavailable),
		retry.WithMax(9),
	}

	return retry.UnaryClientInterceptor(opts...)
}
