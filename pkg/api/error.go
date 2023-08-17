package api

import (
	"fmt"
)

type ApiErrorType int

const (
	ApiConnectionError    ApiErrorType = 0
	ApiAuthorizationError ApiErrorType = 1
	ApiServerError        ApiErrorType = 2
	ApiCredentialsError   ApiErrorType = 3
	ApiResponseError      ApiErrorType = 4
	ApiResponseCodeError  ApiErrorType = 5
)

type ApiError struct {
	ShouldExit   bool
	ErrorType    ApiErrorType
	UserMessage  string
	ErrorMessage error
}

func errorValidatingLicense(status int, message string) error {
	return fmt.Errorf("%d - %s", status, message)
}
