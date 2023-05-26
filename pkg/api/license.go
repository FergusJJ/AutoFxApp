package api

import (
	"encoding/json"
	"fmt"
	"pollo/config"

	"github.com/valyala/fasthttp"
)

func CheckLicense(licenseKey string) (cid string, err error) {
	apiEnpoint, err := config.Config("API_ADDRESS")
	if err != nil {
		return "", err
	}
	requestUri := fmt.Sprintf("http://%s/whop/validate?license=%s", apiEnpoint, licenseKey)
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(requestUri)
	req.Header.SetMethod(fasthttp.MethodGet)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	resp := fasthttp.AcquireResponse()
	if err := fasthttp.Do(req, resp); err != nil {
		return "", err
	}
	jsonResp := validLicenseKeyResponse{}
	if err := json.Unmarshal(resp.Body(), &jsonResp); err != nil {
		return "", err
	}
	//then we have an error
	if jsonResp.Cid == "" {
		jsonErrorResp := apiErrorResponse{}
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
