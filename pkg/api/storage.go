package api

import (
	"encoding/json"
	"fmt"
	"pollo/config"

	"github.com/valyala/fasthttp"
)

func (s *ApiSession) FetchPositions() ([]ApiStoredPositionsResponse, error) {
	var positions []ApiStoredPositionsResponse
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return []ApiStoredPositionsResponse{}, err
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
		return []ApiStoredPositionsResponse{}, err
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		return []ApiStoredPositionsResponse{}, fmt.Errorf("internal server error")
	case fasthttp.StatusUnauthorized: //in this case will need to use refresh token to reauth
		return []ApiStoredPositionsResponse{}, fmt.Errorf("unauthorized request")
	case fasthttp.StatusOK:
		err := json.Unmarshal([]byte(resp.Body()), &positions)
		if err != nil {
			return positions, err
		}
		return positions, nil
	default:
		return []ApiStoredPositionsResponse{}, fmt.Errorf("status code fallthrough: %d", resp.StatusCode())
	}
}
