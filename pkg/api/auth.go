package api

import (
	"encoding/json"
	"fmt"
	"pollo/config"

	"github.com/valyala/fasthttp"
)

// get access & refresh tokens
func (s *ApiSession) FetchApiAuth() *ApiError {

	var apiAuthResponseBody = &apiAuthResponse{}
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return &ApiError{
			ErrorType:    ApiCredentialsError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: err,
		}
	}
	requestUri := fmt.Sprintf("http://%s/api/auth/new?license=%s", apiAddress, s.LicenseKey)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("accept", "application/json")
	resp := fasthttp.AcquireResponse()
	err = fasthttp.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return &ApiError{
			ErrorType:    ApiConnectionError,
			ShouldExit:   false,
			UserMessage:  "An error ocurred whilst getting api auth",
			ErrorMessage: err,
		}
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		return &ApiError{
			ErrorType:    ApiServerError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("internal server error"),
		}
	case fasthttp.StatusBadRequest:
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   true,
			UserMessage:  "An error ocurred whilst authorizing session",
			ErrorMessage: fmt.Errorf("got status 400 authorizing session"),
		}
	case fasthttp.StatusOK:
		if err := json.Unmarshal(resp.Body(), apiAuthResponseBody); err != nil {
			return &ApiError{
				ErrorType:    ApiResponseError,
				ShouldExit:   true,
				UserMessage:  "An error ocurred whilst authorizing session",
				ErrorMessage: err,
			}
		}
		s.accessToken = apiAuthResponseBody.AccessToken
		s.refreshToken = apiAuthResponseBody.RefreshToken
		return nil
	default:
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("status code fallthrough: %d", resp.StatusCode()),
		}
	}
}

// if err != nil && err == "timeout" then retry, otherwise exit?
func (s *ApiSession) RefreshApiAuth() *ApiError {
	type refreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	var apiAuthResponseBody = &apiAuthResponse{}
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return &ApiError{
			ErrorType:    ApiCredentialsError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: err,
		}
	}
	requestBody := &refreshRequest{
		RefreshToken: s.refreshToken,
	}
	requestBodyyBytes, _ := json.Marshal(requestBody)

	requestUri := fmt.Sprintf("http://%s/api/auth/new?license=%s", apiAddress, s.LicenseKey)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.SetContentTypeBytes([]byte("application/json"))
	req.SetBodyRaw(requestBodyyBytes)
	resp := fasthttp.AcquireResponse()
	err = fasthttp.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return &ApiError{
			ErrorType:    ApiConnectionError,
			ShouldExit:   false,
			UserMessage:  "An error ocurred whilst refreshing api auth",
			ErrorMessage: err,
		}
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		return &ApiError{
			ErrorType:    ApiServerError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("internal server error"),
		}
	case fasthttp.StatusBadRequest:
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   true,
			UserMessage:  "An error ocurred whilst authorizing session",
			ErrorMessage: fmt.Errorf("got status 400 authorizing session"),
		}
	case fasthttp.StatusUnauthorized: //shouldn't happen
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   true,
			UserMessage:  "An error ocurred whilst authorizing session",
			ErrorMessage: fmt.Errorf("got status 401 authorizing session"),
		}
	case fasthttp.StatusOK:
		if err := json.Unmarshal(resp.Body(), apiAuthResponseBody); err != nil {
			return &ApiError{
				ErrorType:    ApiResponseError,
				ShouldExit:   true,
				UserMessage:  "An unexpected error occurred",
				ErrorMessage: err,
			}
		}
		s.accessToken = apiAuthResponseBody.AccessToken
		s.refreshToken = apiAuthResponseBody.RefreshToken
	default:
		return &ApiError{
			ErrorType:    ApiResponseError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred",
			ErrorMessage: fmt.Errorf("RefreshApiAuth status code fallthrough: %d", resp.StatusCode()),
		}
	}
	return nil
}
