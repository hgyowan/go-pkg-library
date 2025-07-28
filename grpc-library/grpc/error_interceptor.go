package grpc

import (
	"context"
	"encoding/json"
	pkgError "github.com/hgyowan/go-pkg-library/error"
	pkgLogger "github.com/hgyowan/go-pkg-library/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net/http"
)

func CustomErrorInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	resp, err = handler(ctx, req)
	if err != nil {
		pkgLogger.ZapLogger.Logger.Error(err.Error())
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
