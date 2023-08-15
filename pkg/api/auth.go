package api

import (
	"encoding/json"
	"fmt"
	"pollo/config"

	"github.com/valyala/fasthttp"
)

// get access & refresh tokens
func (s *ApiSession) FetchApiAuth() error {

	var apiAuthResponseBody = &apiAuthResponse{}
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return err
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
		return err
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		return fmt.Errorf("internal server error")
	case fasthttp.StatusBadRequest:
		return fmt.Errorf("bad request")
	case fasthttp.StatusOK:
		if err := json.Unmarshal(resp.Body(), apiAuthResponseBody); err != nil {
			return err
		}
		s.accessToken = apiAuthResponseBody.AccessToken
		s.refreshToken = apiAuthResponseBody.RefreshToken
	default:
		return fmt.Errorf("status code fallthrough: %d", resp.StatusCode())
	}
	return nil
}

// if err != nil && err == "timeout" then retry, otherwise exit?
func (s *ApiSession) RefreshApiAuth() error {
	type refreshRequest struct {
		RefreshToken string `json:"refreshToken"`
	}
	var apiAuthResponseBody = &apiAuthResponse{}
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return err
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
		errName, known := httpConnError(err)
		if known {
			return fmt.Errorf(errName)
		}
		return err
	}
	defer fasthttp.ReleaseResponse(resp)

	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		return fmt.Errorf("internal server error")
	case fasthttp.StatusBadRequest:
		return fmt.Errorf("bad request")
	case fasthttp.StatusUnauthorized: //shouldn't happen
		return fmt.Errorf("unauthorized request")
	case fasthttp.StatusOK:
		if err := json.Unmarshal(resp.Body(), apiAuthResponseBody); err != nil {
			return err
		}
		s.accessToken = apiAuthResponseBody.AccessToken
		s.refreshToken = apiAuthResponseBody.RefreshToken
	default:
		return fmt.Errorf("status code fallthrough: %d", resp.StatusCode())
	}

	return nil
}
