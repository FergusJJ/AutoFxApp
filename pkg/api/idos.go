package api

import "github.com/fasthttp/websocket"

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
	CopyPID     string  `json:"copyPID"`
	SymbolID    int     `json:"symbolID"`
	Price       float64 `json:"price"`
	Volume      int     `json:"volume"`
	Direction   string  `json:"direction"`
	MessageType string  `json:"type"` //close or open
}

type ApiStoredPosition struct {
	CopyPositionID string `json:"copyPositionID"`
	PositionID     string `json:"positionID"`
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
