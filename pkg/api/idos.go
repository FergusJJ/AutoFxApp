package api

import "github.com/fasthttp/websocket"

type ApiSession struct {
	Client ApiClient
	Cid    string
}

type ApiClient struct {
	Connection     *websocket.Conn
	CurrentMessage chan []byte
}

type ApiMonitorMessage struct {
	MessageType        string  `json:"messageType"`
	ID                 string  `json:"id"`
	Created            int     `json:"created"` //should convert to epochTS?
	Symbol             string  `json:"symbol"`
	Volume             float64 `json:"volume"`
	Direction          string  `json:"direction"`
	EntryPrice         float64 `json:"entryPrice"`
	CurrentPrice       float64 `json:"currentPrice"`
	Swap               float64 `json:"swap"`
	Commissons         float64 `json:"commissions"`
	ClosingCommissions float64 `json:"closingCommissions"`
	Pips               float64 `json:"pips"`
	GrossProfitEst     float64 `json:"grossProfitEst"`
	NetProfitEst       float64 `json:"netProfitEst"`
	Channel            string  `json:"channel"`
}

type apiErrorResponse struct {
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message"`
}

type validLicenseKeyResponse struct {
	ResponseCode int    `json:"responseCode"`
	Cid          string `json:"cid"`
}
