package api

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/valyala/fasthttp"
)

func errorValidatingLicense(status int, message string) error {
	return fmt.Errorf("%d - %s", status, message)
}

func httpConnError(err error) (string, bool) {
	var (
		errName string
		known   = true
	)

	switch {
	case errors.Is(err, fasthttp.ErrTimeout):
		errName = "timeout"
	case errors.Is(err, fasthttp.ErrNoFreeConns):
		errName = "conn_limit"
	case errors.Is(err, fasthttp.ErrConnectionClosed):
		errName = "conn_close"
	case reflect.TypeOf(err).String() == "*net.OpError":
		errName = "timeout"
	default:
		known = false
	}

	return errName, known
}
