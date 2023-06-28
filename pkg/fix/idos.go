package fix

import "crypto/tls"

const TradePort = 5212
const PricePort = 5211

type FxUser struct {
	HostName     string `json:"hostName"`
	Password     string `json:"password"`
	SenderCompID string `json:"senderCompID"`
	TargetCompID string `json:"targetCompID"`
	SenderSubID  string `json:"senderSubID"`
}

type FxSession struct {
	GotSecurityList       bool
	LoggedIn              bool
	SecListPairs          map[string]string
	MessageSequenceNumber int
	Connection            *tls.Conn
}

type FxResponse struct {
	body []byte
	err  error
}

type OrderData struct {
	Symbol    string  `json:"symbol"`
	Volume    float64 `json:"volume"`
	Direction string  `json:"direction"` //"buy"/"sell"
	OrderType string  `json:"orderType"` //"market"
}

type ExecutionRe

type CtraderSessionMessageType int
type CtraderQualifier string
type CtraderParamIds int

const (
	Logon CtraderSessionMessageType = iota
	Logout
	Heartbeat
	TestRequest
	Resend
	Reject
	SequenceReset
	OrderStatusRequest
	OrderMassStatusRequest
	RequestForPositions
	NewOrderSingle
	SecurityListRequest
)

const (
	EncryptionMethod  CtraderParamIds = 98
	HeartbeatInterval CtraderParamIds = 108
	ResetSequence     CtraderParamIds = 141
	Username          CtraderParamIds = 553
	Password          CtraderParamIds = 554

	//35 = (H, D)
	ClOrdID CtraderParamIds = 11
	Side    CtraderParamIds = 54

	//35=D
	NOSSymbol       CtraderParamIds = 55
	NOSTransactTime CtraderParamIds = 60
	NOSOrderQty     CtraderParamIds = 38
	NOSOrdType      CtraderParamIds = 40
	NOSPrice        CtraderParamIds = 44  //Not required
	NOSStopPx       CtraderParamIds = 99  //Not required
	NOSExpireTime   CtraderParamIds = 126 //Not Required
	NOSDesignation  CtraderParamIds = 494 //Not Required

	//35=AF OMSR
	MassStatusReqID   CtraderParamIds = 584
	MassStatusReqType CtraderParamIds = 585
	//35=(AN, AP)
	PosReqID CtraderParamIds = 710

	//35=AP
	TotalNumPosReports CtraderParamIds = 727
	PosReqResult       CtraderParamIds = 728

	//35=x
	SecurityReqID           CtraderParamIds = 320
	SecurityListRequestType CtraderParamIds = 559

	HeaderBeginString           CtraderParamIds = 8
	HeaderMessageLength         CtraderParamIds = 9
	HeaderMessageType           CtraderParamIds = 35
	HeaderSenderCompId          CtraderParamIds = 49
	HeaderTargetCompId          CtraderParamIds = 56
	HeaderTargetSubId           CtraderParamIds = 57
	HeaderSenderSubId           CtraderParamIds = 50
	HeaderMessageSequenceNumber CtraderParamIds = 34
	HeaderMessageTimestamp      CtraderParamIds = 52

	TrailerChecksum CtraderParamIds = 10
)

const (
	YYYYMMDDhhmmss = "20060102-15:04:05"
)

// MessageKeyValuePairs probably should be a map, names in map would give more context to values as well
var MessageKeyValuePairs = map[CtraderParamIds]map[string]string{
	//For Logon
	EncryptionMethod:  map[string]string{"encryptionEnabled": "1", "encryptionDisabled": "0"},
	HeartbeatInterval: map[string]string{"noHeartbeat": "0"},
	ResetSequence:     map[string]string{"resetEnabled": "Y", "resetDisabled": "N"},
	Username:          map[string]string{},
	Password:          map[string]string{},

	//New Order Single
	ClOrdID:         map[string]string{}, //want to use uuid4 for this
	NOSSymbol:       map[string]string{}, // determined at runtime
	Side:            map[string]string{"buy": "1", "sell": "2"},
	NOSTransactTime: map[string]string{}, //determined at runtime
	NOSOrderQty:     map[string]string{}, //max precision is 0.01
	NOSOrdType:      map[string]string{"market": "1", "limit": "2", "stop": "3"},
	NOSPrice:        map[string]string{}, //Limit orders, may want to allow user to specify the slippage as a %
	NOSStopPx:       map[string]string{}, //Stop orders,
	NOSExpireTime:   map[string]string{}, //self explanatory, probably not going to use
	NOSDesignation:  map[string]string{},

	//Standard Headers
	HeaderBeginString:           map[string]string{"begin": "FIX.4.4"},
	HeaderMessageLength:         map[string]string{},
	HeaderMessageType:           map[string]string{"0": "A", "1": "5", "2": "0", "3": "1", "4": "2", "5": "3", "6": "4", "7": "H", "8": "AF", "9": "AN", "10": "D", "11": "x"},
	HeaderSenderCompId:          map[string]string{},
	HeaderTargetCompId:          map[string]string{"compId": "CSERVER"},
	HeaderTargetSubId:           map[string]string{"trade": "TRADE", "quote": "QUOTE"},
	HeaderSenderSubId:           map[string]string{"trade": "TRADE", "quote": "QUOTE"},
	HeaderMessageSequenceNumber: map[string]string{},
	HeaderMessageTimestamp:      map[string]string{},

	//Trailer
	TrailerChecksum: map[string]string{},
}
