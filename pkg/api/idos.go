package api

import (
	"github.com/fasthttp/websocket"
)

type ApiSession struct {
	Client       ApiClient
	accountId    int
	Pools        string
	Cid          string
	LicenseKey   string
	refreshToken string
	accessToken  string
}

type ApiClient struct {
	Connection     *websocket.Conn
	CurrentMessage chan []byte
}

type ApiMonitorMessage struct {
	CopyPID         string  `json:"copyPID"`
	SymbolID        int     `json:"symbolID"`
	Symbol          string  `json:"symbol"`
	Price           float64 `json:"price"`
	Volume          int     `json:"volume"`
	Direction       string  `json:"direction"`
	MessageType     string  `json:"type"` //close or open
	OpenedTimestamp int     `json:"openedTimestamp"`
}

type ApiStoredPosition struct {
	CopyPositionID  string `json:"copyPositionID"`
	PositionID      string `json:"positionID"`
	OpenedTimestamp string `json:"openedTimestamp"`
	Symbol          string `json:"symbol"`
	SymbolID        string `json:"symbolID"`
	Volume          string `json:"volume"`
	Side            string `json:"Side"`
	AveragePrice    string `json:"averagePrice"`
}

type apiErrorResponse struct {
	ResponseCode int    `json:"responseCode"`
	Message      string `json:"message"`
}

type validLicenseKeyResponse struct {
	ResponseCode int    `json:"responseCode"`
	Cid          string `json:"cid"`
}

type apiAuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}
