package api

import (
	"encoding/json"
	"fmt"

	"pollo/config"

	"github.com/valyala/fasthttp"
)

func (s *ApiSession) FetchPositions() ([]ApiStoredPosition, *ApiError) {
	var positions []ApiStoredPosition
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		apiErr := &ApiError{
			ErrorType:    ApiCredentialsError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: err,
		}
		return []ApiStoredPosition{}, apiErr
	}
	requestUri := fmt.Sprintf("http://%s/api/user/position/all", apiAddress)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.accessToken))

	resp := fasthttp.AcquireResponse()
	err = fasthttp.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		apiErr := &ApiError{
			ErrorType:    ApiConnectionError,
			ShouldExit:   false,
			UserMessage:  "An error ocurred whilst getting positions",
			ErrorMessage: err,
		}
		return []ApiStoredPosition{}, apiErr
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		apiErr := &ApiError{
			ErrorType:    ApiServerError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("internal server error"),
		}
		return []ApiStoredPosition{}, apiErr
	case fasthttp.StatusUnauthorized: //in this case will need to use refresh token to reauth
		apiErr := &ApiError{
			ErrorType:    ApiAuthorizationError,
			ShouldExit:   false,
			UserMessage:  "Session not authorized",
			ErrorMessage: fmt.Errorf("unauthorized request"),
		}
		return []ApiStoredPosition{}, apiErr
	case fasthttp.StatusOK:

		err := json.Unmarshal([]byte(resp.Body()), &positions)
		if err != nil {
			apiErr := &ApiError{
				ErrorType:    ApiResponseError,
				ShouldExit:   true,
				UserMessage:  "An unexpected error occurred, please try again later",
				ErrorMessage: err,
			}
			return positions, apiErr
		}
		return positions, nil
	default:
		apiErr := &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("status code fallthrough: %d", resp.StatusCode()),
		}
		return []ApiStoredPosition{}, apiErr
	}
}

func (s *ApiSession) RemovePosition(pid string) *ApiError {
	type storePositionResponse struct {
		Message string `json:"message"`
	}
	type deletePositionReqBody struct {
		PositionID string `json:"positionID"`
	}
	var requestBody = &deletePositionReqBody{
		PositionID: pid,
	}
	var responseBody storePositionResponse
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return &ApiError{
			ErrorType:    ApiCredentialsError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: err,
		}
	}
	requestBodyyBytes, _ := json.Marshal(requestBody)
	requestUri := fmt.Sprintf("http://%s/api/user/position/delete", apiAddress)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.accessToken))
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
			UserMessage:  "An error ocurred whilst removing a position",
			ErrorMessage: err,
		}
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusBadRequest:
		responseStruct := &storePositionResponse{}
		_ = json.Unmarshal(resp.Body(), responseStruct)
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ErrorMessage: fmt.Errorf(responseStruct.Message),
			UserMessage:  "An unexpected error occurred whilst removing a position",
			ShouldExit:   true,
		}
	case fasthttp.StatusInternalServerError:
		apiErr := &ApiError{
			ErrorType:    ApiServerError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("internal server error"),
		}
		return apiErr
	case fasthttp.StatusUnauthorized: //in this case will need to use refresh token to reauth
		apiErr := &ApiError{
			ErrorType:    ApiAuthorizationError,
			ShouldExit:   false,
			UserMessage:  "Session not authorized",
			ErrorMessage: fmt.Errorf("unauthorized request"),
		}
		return apiErr
	case fasthttp.StatusAccepted:
		err := json.Unmarshal([]byte(resp.Body()), &responseBody)
		if err != nil {
			return &ApiError{
				ErrorType:    ApiResponseError,
				ShouldExit:   true,
				UserMessage:  "An unexpected error occurred, please try again later",
				ErrorMessage: err,
			}
		}
		return nil
	default:
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   false,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("status code fallthrough: %d", resp.StatusCode()),
		}
	}
}

func (s *ApiSession) StorePosition(position ApiStoredPosition) *ApiError {
	type storePositionResponse struct {
		Message string `json:"message"`
	}
	var responseBody storePositionResponse
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return &ApiError{
			ErrorType:    ApiCredentialsError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: err,
		}
	}
	requestBodyyBytes, _ := json.Marshal(position)
	requestUri := fmt.Sprintf("http://%s/api/user/position/new", apiAddress)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.Set("accept", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", s.accessToken))
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
			UserMessage:  "An error ocurred whilst Storing positions",
			ErrorMessage: err,
		}
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusBadRequest:
		responseStruct := &storePositionResponse{}
		_ = json.Unmarshal(resp.Body(), responseStruct)
		if responseStruct.Message == "unable to store duplicate position" {
			return &ApiError{
				ErrorType:    ApiResponseCodeError,
				ShouldExit:   false,
				UserMessage:  fmt.Sprintf("An error ocurred whilst storing position %s", position.PositionID),
				ErrorMessage: fmt.Errorf("unable to store duplicate position: %s", position.PositionID),
			}
		}
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   true,
			UserMessage:  fmt.Sprintf("An error ocurred whilst storing position %s", position.PositionID),
			ErrorMessage: fmt.Errorf("%s: %s", responseStruct.Message, position.PositionID),
		}
	case fasthttp.StatusInternalServerError:
		apiErr := &ApiError{
			ErrorType:    ApiServerError,
			ShouldExit:   true,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("internal server error"),
		}
		return apiErr
	case fasthttp.StatusUnauthorized: //in this case will need to use refresh token to reauth
		apiErr := &ApiError{
			ErrorType:    ApiAuthorizationError,
			ShouldExit:   false,
			UserMessage:  "Session not authorized",
			ErrorMessage: fmt.Errorf("unauthorized request"),
		}
		return apiErr
	case fasthttp.StatusAccepted:
		err := json.Unmarshal([]byte(resp.Body()), &responseBody)
		if err != nil {
			return &ApiError{
				ErrorType:    ApiResponseError,
				ShouldExit:   true,
				UserMessage:  "An unexpected error occurred, please try again later",
				ErrorMessage: err,
			}
		}
		return nil
	default:
		return &ApiError{
			ErrorType:    ApiResponseCodeError,
			ShouldExit:   false,
			UserMessage:  "An unexpected error occurred, please try again later",
			ErrorMessage: fmt.Errorf("status code fallthrough: %d", resp.StatusCode()),
		}
	}
}
