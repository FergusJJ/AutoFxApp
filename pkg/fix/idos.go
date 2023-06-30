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

type MessageBodyAndTag struct {
	Tag         string            `json:"tag"`
	MessageBody map[string]string `json:"body"`
}

type OrderData struct {
	Symbol    string  `json:"symbol"`
	Volume    float64 `json:"volume"`
	Direction string  `json:"direction"` //"buy"/"sell"
	OrderType string  `json:"orderType"` //"market"
}

type ExecutionReport struct {
	OrderID         string `json:"OrderID"`
	ClOrdID         string `json:"ClOrdID,omitempty"`
	TotNumReports   string `json:"TotNumReports,omitempty"`
	ExecType        string `json:"ExecType"`
	OrdStatus       string `json:"OrdStatus"`
	Symbol          string `json:"Symbol,omitempty"`
	Side            string `json:"Side,omitempty"`
	TransactTime    string `json:"TransactTime,omitempty"`
	AvgPx           string `json:"AvgPx,omitempty"`
	OrderQty        string `json:"OrderQty,omitempty"`
	LeavesQty       string `json:"LeavesQty,omitempty"`
	CumQty          string `json:"CumQty,omitempty"`
	LastQty         string `json:"LastQty,omitempty"`
	OrdType         string `json:"OrdType,omitempty"`
	Price           string `json:"Price,omitempty"`
	StopPx          string `json:"StopPx,omitempty"`
	TimeInForce     string `json:"TimeInForce,omitempty"`
	ExpireTime      string `json:"ExpireTime,omitempty"`
	Text            string `json:"Text,omitempty"`
	OrdRejReason    string `json:"OrdRejReason,omitempty"`
	PosMaintRptID   string `json:"PosMaintRptID,omitempty"`
	Designation     string `json:"Designation,omitempty"`
	MassStatusReqID string `json:"MassStatusReqID,omitempty"`
	AbsoluteTP      string `json:"AbsoluteTP,omitempty"`
	RelativeTP      string `json:"RelativeTP,omitempty"`
	AbsoluteSL      string `json:"AbsoluteSL,omitempty"`
	RelativeSL      string `json:"RelativeSL,omitempty"`
	TrailingSL      string `json:"TrailingSL,omitempty"`
	TriggerMethodSL string `json:"TriggerMethodSL,omitempty"`
	GuaranteedSL    string `json:"GuaranteedSL,omitempty"`
}

var executionReportTagMapping = map[string]string{
	"37":   "OrderID",
	"11":   "ClOrdID",
	"911":  "TotNumReports",
	"150":  "ExecType",
	"39":   "OrdStatus",
	"55":   "Symbol",
	"54":   "Side",
	"60":   "TransactTime",
	"6":    "AvgPx",
	"38":   "OrderQty",
	"151":  "LeavesQty",
	"14":   "CumQty",
	"32":   "LastQty",
	"40":   "OrdType",
	"44":   "Price",
	"99":   "StopPx",
	"59":   "TimeInForce",
	"126":  "ExpireTime",
	"58":   "Text",
	"103":  "OrdRejReason",
	"721":  "PosMaintRptID",
	"494":  "Designation",
	"584":  "MassStatusReqID",
	"1000": "AbsoluteTP",
	"1001": "RelativeTP",
	"1002": "AbsoluteSL",
	"1003": "RelativeSL",
	"1004": "TrailingSL",
	"1005": "TriggerMethodSL",
	"1006": "GuaranteedSL",
}

type OrderCancelReject struct {
	OrderID          string `json:"OrderID"`
	ClOrdID          string `json:"ClOrdID"`
	OrigClOrdID      string `json:"OrigClOrdID"`
	OrdStatus        string `json:"OrdStatus"`
	CxlRejResponseTo string `json:"CxlRejResponseTo"`
	Text             string `json:"Text"`
}

var orderCancelRejectTagMapping = map[string]string{
	"37":  "OrderID",
	"11":  "ClOrdID",
	"41":  "OrigClOrdID",
	"39":  "OrdStatus",
	"434": "CxlRejResponseTo",
	"58":  "Text",
}

type BusinessMessageReject struct {
	RefSeqNum            string `json:"RefSeqNum"`
	RefMsgType           string `json:"RefMsgType"`
	BusinessRejectRefID  string `json:"BusinessRejectRefID"`
	BusinessRejectReason string `json:"BusinessRejectReason"`
	Text                 string `json:"Text"`
}

var businessMessageRejectTagMapping = map[string]string{
	"45":  "RefSeqNum",
	"372": "RefMsgType",
	"379": "BusinessRejectRefID",
	"380": "BusinessRejectReason",
	"58":  "Text",
}

type PositionReport struct {
	PosReqID           string `json:"PosReqID"`
	PosMaintRptID      string `json:"PosMaintRptID,omitempty"`
	TotalNumPosReports string `json:"TotalNumPosReports"`
	PosReqResult       string `json:"PosReqResult"`
	Symbol             string `json:"Symbol,omitempty"`
	NoPositions        string `json:"NoPositions,omitempty"`
	LongQty            string `json:"LongQty,omitempty"`
	ShortQty           string `json:"ShortQty,omitempty"`
	SettlPrice         string `json:"SettlPrice,omitempty"`
	AbsoluteTP         string `json:"AbsoluteTP,omitempty"`
	AbsoluteSL         string `json:"AbsoluteSL,omitempty"`
	TrailingSL         string `json:"TrailingSL,omitempty"`
	TriggerMethodSL    string `json:"TriggerMethodSL,omitempty"`
	GuaranteedSL       string `json:"GuaranteedSL,omitempty"`
}

var positionReportTagMapping = map[string]string{
	"710":  "PosReqID",
	"721":  "PosMaintRptID",
	"727":  "TotalNumPosReports",
	"728":  "PosReqResult",
	"55":   "Symbol",
	"702":  "NoPositions",
	"704":  "LongQty",
	"705":  "ShortQty",
	"730":  "SettlPrice",
	"1000": "AbsoluteTP",
	"1002": "AbsoluteSL",
	"1004": "TrailingSL",
	"1005": "TriggerMethodSL",
	"1006": "GuaranteedSL",
}

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
	EncryptionMethod:  {"encryptionEnabled": "1", "encryptionDisabled": "0"},
	HeartbeatInterval: {"noHeartbeat": "0"},
	ResetSequence:     {"resetEnabled": "Y", "resetDisabled": "N"},
	Username:          {},
	Password:          {},

	//New Order Single
	ClOrdID:         {}, //want to use uuid4 for this
	NOSSymbol:       {}, // determined at runtime
	Side:            {"buy": "1", "sell": "2"},
	NOSTransactTime: {}, //determined at runtime
	NOSOrderQty:     {}, //max precision is 0.01
	NOSOrdType:      {"market": "1", "limit": "2", "stop": "3"},
	NOSPrice:        {}, //Limit orders, may want to allow user to specify the slippage as a %
	NOSStopPx:       {}, //Stop orders,
	NOSExpireTime:   {}, //self explanatory, probably not going to use
	NOSDesignation:  {},

	//Standard Headers
	HeaderBeginString:           {"begin": "FIX.4.4"},
	HeaderMessageLength:         {},
	HeaderMessageType:           {"0": "A", "1": "5", "2": "0", "3": "1", "4": "2", "5": "3", "6": "4", "7": "H", "8": "AF", "9": "AN", "10": "D", "11": "x"},
	HeaderSenderCompId:          {},
	HeaderTargetCompId:          {"compId": "CSERVER"},
	HeaderTargetSubId:           {"trade": "TRADE", "quote": "QUOTE"},
	HeaderSenderSubId:           {"trade": "TRADE", "quote": "QUOTE"},
	HeaderMessageSequenceNumber: {},
	HeaderMessageTimestamp:      {},

	//Trailer
	TrailerChecksum: {},
}
