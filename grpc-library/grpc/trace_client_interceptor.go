package grpc

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TraceClientInterceptor() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{},
		cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx, _ = metadataFromContext(ctx)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// metadataCarrier implements the TextMapCarrier interface for gRPC metadata
type metadataCarrier struct {
	md metadata.MD
}

// Ensure it implements propagation.TextMapCarrier
var _ propagation.TextMapCarrier = (*metadataCarrier)(nil)

// Get returns the value associated with the passed key.
func (c metadataCarrier) Get(key string) string {
	if c.md == nil {
		return ""
	}
	values := c.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

// Set stores the key-value pair.
func (c metadataCarrier) Set(key, value string) {
	if c.md == nil {
		c.md = metadata.Pairs()
	}
	c.md.Set(key, value)
}

// Keys lists all the keys stored in this carrier.
func (c metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for k := range c.md {
		keys = append(keys, k)
	}
	return keys
}

// Helper to inject context into outgoing gRPC metadata
func metadataFromContext(ctx context.Context) (context.Context, metadata.MD) {
	md, ok := metadata.FromOutgoingContext(ctx)
	if !ok {
		md = metadata.New(nil)
	}
	carrier := metadataCarrier{md: md}
	otel.GetTextMapPropagator().Inject(ctx, carrier)
	return metadata.NewOutgoingContext(ctx, carrier.md), carrier.md
}
