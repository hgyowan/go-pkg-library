package context

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/grpc/metadata"
	"strconv"
)

type DefaultContext struct {
	context.Context
	// context metadata 에 키가 존재하는지 검사하는 Option
	valid       ContextValidOption
	contextData ContextData
}

type ContextData struct {
	UserID      uint
	RequestID   string
	AccessToken string
	IP          string
	UserAgent   string
	error       error
}

type ContextValidOption struct {
	ValidUserID      bool
	ValidRequestID   bool
	ValidAccessToken bool
	ValidIP          bool
	ValidUserAgent   bool
}

func (cc *ContextValidOption) Apply(defaultContext *DefaultContext) {
	if defaultContext != nil {
		defaultContext.valid.ValidUserID = cc.ValidUserID
		defaultContext.valid.ValidRequestID = cc.ValidRequestID
		defaultContext.valid.ValidAccessToken = cc.ValidAccessToken
		defaultContext.valid.ValidIP = cc.ValidIP
		defaultContext.valid.ValidUserAgent = cc.ValidUserAgent
	}
}

type Option interface {
	Apply(defaultContext *DefaultContext)
}

// Deprecated: NewContext Use ParseContent

func IncomingContext(ctx context.Context, option ...Option) *DefaultContext {
	// default 설정
	dc := &DefaultContext{
		Context: ctx,
		valid:   ContextValidOption{},
	}

	for _, opt := range option {
		opt.Apply(dc)
	}

	return dc
}

func OutgoingContext(ctx context.Context, option ...Option) *DefaultContext {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ctx = metadata.NewOutgoingContext(context.Background(), md)
	} else {
		ctx = context.Background()
	}
	// default 설정
	dc := &DefaultContext{
		Context: ctx,
		valid:   ContextValidOption{},
	}

	for _, opt := range option {
		opt.Apply(dc)
	}

	return dc
}

func (dc *DefaultContext) UserID() *DefaultContext {
	if md, ok := metadata.FromIncomingContext(dc.Context); ok {
		v := md.Get("user_id")
		if len(v) > 0 {
			userID, err := strconv.ParseUint(v[0], 10, 64)
			if err != nil {
				dc.Error(errors.New("user_id is not a valid uint"))
				return dc
			}

			dc.contextData.UserID = uint(userID)
		} else {
			if dc.valid.ValidUserID {
				dc.Error(errors.New("user_id is not supplied"))
			}
		}
	}

	return dc
}

func (dc *DefaultContext) AddUserID(userID string) *DefaultContext {
	dc.Context = metadata.AppendToOutgoingContext(dc.Context, "user_id", userID)

	return dc
}

func (dc *DefaultContext) AddRequestId(requestId string) *DefaultContext {
	dc.Context = metadata.AppendToOutgoingContext(dc.Context, "request_id", requestId)

	return dc
}

func (dc *DefaultContext) RequestId() *DefaultContext {
	if md, ok := metadata.FromIncomingContext(dc.Context); ok {
		v := md.Get("request_id")
		if len(v) > 0 {
			dc.contextData.RequestID = v[0]
		} else {
			if dc.valid.ValidRequestID {
				dc.Error(errors.New("request_id is not supplied"))
			}
		}
	}

	return dc
}

func (dc *DefaultContext) AccessToken() *DefaultContext {
	if md, ok := metadata.FromIncomingContext(dc.Context); ok {
		v := md.Get("access_token")
		if len(v) > 0 {
			dc.contextData.AccessToken = v[0]
		} else {
			if dc.valid.ValidUserID {
				dc.Error(errors.New("access_token is not supplied"))
			}
		}
	}

	return dc
}

func (dc *DefaultContext) AddAccessToken(accessToken string) *DefaultContext {
	dc.Context = metadata.AppendToOutgoingContext(dc.Context, "access_token", accessToken)

	return dc
}

func (dc *DefaultContext) IP() *DefaultContext {
	if md, ok := metadata.FromIncomingContext(dc.Context); ok {
		v := md.Get("ip")
		if len(v) > 0 {
			dc.contextData.IP = v[0]
		} else {
			if dc.valid.ValidIP {
				dc.Error(errors.New("ip is not supplied"))
			}
		}
	}

	return dc
}

func (dc *DefaultContext) UserAgent() *DefaultContext {
	if md, ok := metadata.FromIncomingContext(dc.Context); ok {
		v := md.Get("user_agent")
		if len(v) > 0 {
			dc.contextData.UserAgent = v[0]
		} else {
			if dc.valid.ValidUserAgent {
				dc.Error(errors.New("user_agent is not supplied"))
			}
		}
	}

	return dc
}

func (dc *DefaultContext) Scan() (ContextData, error) {
	if dc.contextData.error != nil {
		return ContextData{}, errors.New(dc.contextData.error.Error())
	}

	return dc.contextData, nil
}

func (dc *DefaultContext) Error(err error) {
	if dc.contextData.error == nil {
		dc.contextData.error = errors.New(err.Error())
	} else {
		dc.contextData.error = fmt.Errorf("%w and %v", err, dc.contextData.error)
	}
}
