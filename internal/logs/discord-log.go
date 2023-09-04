package logs

import (
	"fmt"
	"log"
	"pollo/config"
	"reflect"
	"time"

	"github.com/valyala/fasthttp"
)

const maxRetries = 3
const retryDelay = 5 * time.Second

func SendApplicationLog(errToSend error, license string) {

	webhookUri, err := config.Config("LOGGING_URL")
	if err != nil {
		log.Fatal(err)
	}

	errorMessage := errToSend.Error()
	errorType := reflect.TypeOf(errToSend).String()
	webhookBody := fmt.Sprintf(`{
		"content": null,
		"embeds": [
			{
				"title": "Error Information",
				"color": 16711680,
				"fields": [
					{
						"name": "Error Message",
						"value": "%s",
						"inline": true
					},
					{
						"name": "Error Type",
						"value": "%s",
						"inline": true
					}
				]
			},
			{
				"title": "User Information",
				"color": 16711680,
				"fields": [
					{
						"name": "license",
						"value": "%s",
						"inline": true
					}
				]
			}
		],
		"attachments": []
	}`, errorMessage, errorType, license)
	for i := 0; i < maxRetries; i++ {
		req := fasthttp.AcquireRequest()
		req.SetRequestURI(webhookUri)
		req.Header.SetMethod(fasthttp.MethodPost)
		req.Header.SetContentTypeBytes([]byte("application/json"))
		req.SetBodyRaw([]byte(webhookBody))
		resp := fasthttp.AcquireResponse()
		err = fasthttp.Do(req, resp)
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)
		if err == nil && resp.StatusCode() == fasthttp.StatusOK {
			return // Successfully sent webhook

		}

		// If we reached here, it means the webhook failed. Wait for a while and retry.
		time.Sleep(retryDelay)
	}

}
