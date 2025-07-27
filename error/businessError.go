package error

import (
	"encoding/json"
	"errors"
	"google.golang.org/grpc/status"
)

type BusinessError struct {
	error
	Status *Status
}

func EmptyBusinessError() *BusinessError {
	return &BusinessError{
		error:  errors.New("empty error"),
		Status: &Status{},
	}
}

// WrapWithCode
// 새로운 code 의 business error 를 생성하고자 하는 경우 첫번째 인자에 EmptyBusinessError() 를 넣어주세요.
func WrapWithCode(err error, code Code, details ...string) error {
	if err == nil {
		return nil
	}

	var status Status
	if s, ok := businessCodeMap[code]; ok {
		status = s
	} else {
		status = businessCodeMap[None]
	}

	if errors.Is(err, EmptyBusinessError()) {
		err = &status
	}

	if len(details) > 0 {
		status.Detail = details
	}

	return &BusinessError{
		error:  WrapBusiness(err, status.Message),
		Status: &status,
	}
}

// WrapWithCodeAndData
// 새로운 code 의 business error 를 생성하고자 하는 경우 첫번째 인자에 EmptyBusinessError() 를 넣어주세요.
func WrapWithCodeAndData(err error, code Code, data interface{}, details ...string) error {
	if err == nil {
		return nil
	}

	var status Status
	if s, ok := businessCodeMap[code]; ok {
		status = s
	} else {
		status = businessCodeMap[None]
	}

	if errors.Is(err, EmptyBusinessError()) {
		err = &status
	}

	if len(details) > 0 {
		status.Detail = details
	}

	if data != nil {
		status.Data = data
	}

	return &BusinessError{
		error:  WrapBusiness(err, status.Message),
		Status: &status,
	}
}

func WrapWithCustomStatus(err error, status Status, details ...string) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, EmptyBusinessError()) {
		err = &status
	}

	if len(details) > 0 {
		status.Detail = details
	}

	return &BusinessError{
		error:  WrapBusiness(err, status.Message),
		Status: &status,
	}
}

func CastBusinessError(err error) (*BusinessError, bool) {
	type grpcstatus interface{ GRPCStatus() *status.Status }
	for err != nil {
		switch err.(type) {
		case *BusinessError:
			return err.(*BusinessError), true
		case grpcstatus:
			st, ok := status.FromError(err)
			if ok {
				var sts Status
				_ = json.Unmarshal([]byte(st.Message()), &sts)
				return &BusinessError{
					error:  err,
					Status: &sts,
				}, true
			}
		}

		err = errors.Unwrap(err)
	}

	return nil, false
}

func CompareBusinessError(err error, code Code) bool {
	if t, ok := CastBusinessError(err); ok {
		if t.Status.Code == int(code) {
			return true
		}
	}

	return false
}
