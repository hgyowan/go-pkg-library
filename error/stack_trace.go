package error

import (
	"fmt"
	"runtime"
)

var wrapFormat = "%s \n %w"
var funcInfoFormat = "{%s:%d} [%s]"

func getFuncInfo(pc uintptr, file string, line int) string {
	f := runtime.FuncForPC(pc)
	if f == nil {
		return fmt.Sprintf(funcInfoFormat, file, line, "unknwon")
	}
	return fmt.Sprintf(funcInfoFormat, file, line, f.Name())
}

func wrap(err error, msg string, skip int) error {
	pc, file, line, ok := runtime.Caller(skip)

	if !ok {
		return fmt.Errorf(wrapFormat, msg, err)
	}

	// {file:line} [funcName] msg
	stack := fmt.Sprintf("%s %s", getFuncInfo(pc, file, line), msg)
	return fmt.Errorf(wrapFormat, stack, err)
}

func WrapWithMessage(err error, msg string) error {
	if err == nil {
		return nil
	}
	return wrap(err, msg, 2)
}

func Wrap(err error) error {
	if err == nil {
		return nil
	}
	return wrap(err, "", 2)
}

func WrapBusiness(err error, msg string) error {
	if err == nil {
		return nil
	}
	return wrap(err, msg, 3)
}
