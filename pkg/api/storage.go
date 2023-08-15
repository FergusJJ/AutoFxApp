package api

import (
	"encoding/json"
	"fmt"
	"pollo/config"

	"github.com/valyala/fasthttp"
)

func (s *ApiSession) FetchPositions() ([]ApiStoredPosition, error) {
	var positions []ApiStoredPosition
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return []ApiStoredPosition{}, err
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
		return []ApiStoredPosition{}, err
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusInternalServerError:
		return []ApiStoredPosition{}, fmt.Errorf("internal server error")
	case fasthttp.StatusUnauthorized: //in this case will need to use refresh token to reauth
		return []ApiStoredPosition{}, fmt.Errorf("unauthorized request")
	case fasthttp.StatusOK:
		err := json.Unmarshal([]byte(resp.Body()), &positions)
		if err != nil {
			return positions, err
		}
		return positions, nil
	default:
		return []ApiStoredPosition{}, fmt.Errorf("status code fallthrough: %d", resp.StatusCode())
	}
}

func (s *ApiSession) StorePosition(position ApiStoredPosition) error {
	type storePositionResponse struct {
		Message string `json:"message"`
	}
	var responseBody storePositionResponse
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return err
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
		return err
	}
	defer fasthttp.ReleaseResponse(resp)
	switch resp.StatusCode() {
	case fasthttp.StatusBadRequest:
		responseStruct := &storePositionResponse{}
		_ = json.Unmarshal(resp.Body(), responseStruct)
		if responseStruct.Message == "unable to store duplicate position" {
			return fmt.Errorf("duplicate position error")
		}

		return fmt.Errorf("bad request payload error: %s", responseStruct.Message)
	case fasthttp.StatusInternalServerError:
		return fmt.Errorf("internal server error")
	case fasthttp.StatusUnauthorized: //in this case will need to use refresh token to reauth
		return fmt.Errorf("unauthorized request")
	case fasthttp.StatusAccepted:
		err := json.Unmarshal([]byte(resp.Body()), &responseBody)
		if err != nil {
			return err
		}
		return nil
	default:
		return fmt.Errorf("status code fallthrough: %d", resp.StatusCode())
	}
}
