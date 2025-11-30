package grpc

import (
	"context"
	"encoding/json"
	"net/http"

	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	otelCodes "go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func CustomErrorUnaryInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		var traceID string
		var spanID string
		span := trace.SpanFromContext(ctx)
		if span != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Code(codes.Internal), err.Error())

			sc := span.SpanContext()
			traceID = sc.TraceID().String()
			spanID = sc.SpanID().String()
		}

		pkgLogger.ZapLogger.Logger.Error(err.Error(), zap.String("trace_id", traceID), zap.String("span_id", spanID))
		castedErr, ok := pkgError.CastBusinessError(err)
		if ok {
			b, _ := json.Marshal(castedErr.Status)
			return nil, status.Errorf(codes.Internal, string(b))
		}

		b, _ := json.Marshal(pkgError.Status{
			Code:           int(pkgError.None),
			HttpStatusCode: http.StatusInternalServerError,
			Message:        err.Error(),
		})
		return nil, status.Errorf(codes.Internal, string(b))
	}
	return resp, nil
}

func CustomErrorStreamInterceptor(
	srv any, req grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler,
) (err error) {
	err = handler(srv, req)
	if err != nil {
		var traceID string
		var spanID string
		span := trace.SpanFromContext(req.Context())
		if span != nil {
			span.RecordError(err)
			span.SetStatus(otelCodes.Code(codes.Internal), err.Error())

			sc := span.SpanContext()
			traceID = sc.TraceID().String()
			spanID = sc.SpanID().String()
		}

		pkgLogger.ZapLogger.Logger.Error(err.Error(), zap.String("trace_id", traceID), zap.String("span_id", spanID))
		castedErr, ok := pkgError.CastBusinessError(err)
		if ok {
			b, _ := json.Marshal(castedErr.Status)
			return status.Errorf(codes.Internal, string(b))
		}

		b, _ := json.Marshal(pkgError.Status{
			Code:           int(pkgError.None),
			HttpStatusCode: http.StatusInternalServerError,
			Message:        err.Error(),
		})
		return status.Errorf(codes.Internal, string(b))
	}
	return nil
}
