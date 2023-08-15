package api

import (
	"encoding/json"
	"fmt"
	"pollo/config"

	"github.com/valyala/fasthttp"
)

func CheckLicense(licenseKey string) (cid string, err error) {
	apiAddress, err := config.Config("API_ADDRESS")
	if err != nil {
		return "", err
	}
	requestUri := fmt.Sprintf("http://%s/whop/validate?license=%s", apiAddress, licenseKey)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Set("accept", "application/json")
	resp := fasthttp.AcquireResponse()
	err = fasthttp.Do(req, resp)
	fasthttp.ReleaseRequest(req)
	if err != nil {
		return "", err
	}
	defer fasthttp.ReleaseResponse(resp)
	jsonResp := validLicenseKeyResponse{}
	if err := json.Unmarshal(resp.Body(), &jsonResp); err != nil {
		return "", err
	}

	if jsonResp.Cid == "" {
		jsonErrorResp := apiErrorResponse{} //where statuscode is the status from WHOP, not the api
		err = json.Unmarshal(resp.Body(), &jsonErrorResp)
		if err != nil {
			return "", err
		}
		err = errorValidatingLicense(jsonErrorResp.ResponseCode, jsonErrorResp.Message)
		return "", err
	}
	cid = jsonResp.Cid
	return cid, nil
}
