package fix

import "sync"

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
	mtx                     sync.Mutex
	GotSecurityList         bool
	LoggedIn                bool
	SecListPairs            map[string]string
	MarketDataSubscriptions map[string]*MarketDataSubscription //key should be symbolID
	Positions               map[string]Position                //key should be copy position id
	MessageSequenceNumber   int
	PriceClient             *FixClient
	TradeClient             *FixClient
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

var executionReportTagMapping = map[int]string{
	37:   "OrderID",
	11:   "ClOrdID",
	911:  "TotNumReports",
	150:  "ExecType",
	39:   "OrdStatus",
	55:   "Symbol",
	54:   "Side",
	60:   "TransactTime",
	6:    "AvgPx",
	38:   "OrderQty",
	151:  "LeavesQty",
	14:   "CumQty",
	32:   "LastQty",
	40:   "OrdType",
	44:   "Price",
	99:   "StopPx",
	59:   "TimeInForce",
	126:  "ExpireTime",
	58:   "Text",
	103:  "OrdRejReason",
	721:  "PosMaintRptID",
	494:  "Designation",
	584:  "MassStatusReqID",
	1000: "AbsoluteTP",
	1001: "RelativeTP",
	1002: "AbsoluteSL",
	1003: "RelativeSL",
	1004: "TrailingSL",
	1005: "TriggerMethodSL",
	1006: "GuaranteedSL",
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

var businessMessageRejectTagMapping = map[int]string{
	45:  "RefSeqNum",
	372: "RefMsgType",
	379: "BusinessRejectRefID",
	380: "BusinessRejectReason",
	58:  "Text",
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

var positionReportTagMapping = map[int]string{
	710:  "PosReqID",
	721:  "PosMaintRptID",
	727:  "TotalNumPosReports",
	728:  "PosReqResult",
	55:   "Symbol",
	702:  "NoPositions",
	704:  "LongQty",
	705:  "ShortQty",
	730:  "SettlPrice",
	1000: "AbsoluteTP",
	1002: "AbsoluteSL",
	1004: "TrailingSL",
	1005: "TriggerMethodSL",
	1006: "GuaranteedSL",
}

type MarketDataSnapshot struct {
	MDReqID      string `json:"MDReqID,omitempty"`      // Tag 262: An ID of the market data request sent previously.
	Symbol       string `json:"Symbol"`                 // Tag 55: Instrument identificators are provided by Spotware.
	NoMDEntries  string `json:"NoMDEntries"`            // Tag 268: The number of entries following.
	MDEntryType  string `json:"MDEntryType,omitempty"`  // Tag 269: 0 = Bid; 1 = Offer.
	QuoteEntryID string `json:"QuoteEntryID,omitempty"` // Tag 299: A unique identification of the quote as a part of QuoteSet.
	MDEntryPx    string `json:"MDEntryPx,omitempty"`    // Tag 270: A price of the Market Data Entry.
	MDEntrySize  string `json:"MDEntrySize,omitempty"`  // Tag 271: Volume of the Market Data Entry.
	MDEntryID    string `json:"MDEntryID,omitempty"`    // Tag 278: A unique Market Data Entry identifier.
}

var marketDataSnapshotTagMapping = map[int]string{
	262: "MDReqID",
	55:  "Symbol",
	268: "NoMDEntries",
	269: "MDEntryType",
	299: "QuoteEntryID",
	270: "MDEntryPx",
	271: "MDEntrySize",
	278: "MDEntryID",
}

type MarketDataIncrementalRefresh struct {
	MDReqID        string  `json:"MDReqID,omitempty"`     // Tag 262
	NoMDEntries    int     `json:"NoMDEntries"`           // Tag 268
	MDUpdateAction string  `json:"MDUpdateAction"`        // Tag 279: 0 = New; 2 = Delete.
	MDEntryType    string  `json:"MDEntryType,omitempty"` // Tag 269
	MDEntryID      string  `json:"MDEntryID"`             // Tag 278
	Symbol         string  `json:"Symbol"`                // Tag 55
	MDEntryPx      float64 `json:"MDEntryPx,omitempty"`   // Tag 270
	MDEntrySize    float64 `json:"MDEntrySize,omitempty"` // Tag 271
}

var marketDataIncrementalRefreshTagMapping = map[int]string{
	262: "MDReqID",        // An ID of the market data request sent previously.
	268: "NoMDEntries",    // The number of entries following.
	279: "MDUpdateAction", // A type of the Market Data update action.
	269: "MDEntryType",    // 0 = Bid; 1 = Offer.
	278: "MDEntryID",      // An ID of the Market Data Entry.
	55:  "Symbol",         // Instrument identifiers provided by Spotware.
	270: "MDEntryPx",      // Required only when MDUpdateAction (tag=279) = 0.
	271: "MDEntrySize",    // Required only when MDUpdateAction (tag=279) = 0.
}

var marketDataRequestRejectTagMapping = map[int]string{
	262: "MDReqID",        // Must refer to MDReqID of the request.
	281: "MDReqRejReason", // 0 = Unknown symbol; 4 = Unsupported SubscriptionRequestType; 5 = Unsupported MarketDepth.
}

type MarketDataRequestReject struct {
	MDReqID        string `json:"MDReqID"`                  // Tag 262
	MDReqRejReason int    `json:"MDReqRejReason,omitempty"` // Tag 281: 0 = Unknown symbol; 4 = Unsupported SubscriptionRequestType; 5 = Unsupported MarketDepth.
}

var sessionRejectMessageTagMapping = map[int]string{
	45:  "RefSeqNum",
	58:  "Text",
	354: "EncodedTextLen",
	355: "EncodedText",
	371: "RefTagID",
	372: "RefMsgType",
	373: "SessionRejectReason",
}

type SessionRejectMessage struct {
	RefSeqNum           string `json:"RefSeqNum"`                     // Tag 45
	Text                string `json:"Text,omitempty"`                // Tag 58
	EncodedTextLen      string `json:"EncodedTextLen,omitempty"`      // Tag 354
	EncodedText         string `json:"EncodedText,omitempty"`         // Tag 355
	RefTagID            string `json:"RefTagID,omitempty"`            // Tag 371
	RefMsgType          string `json:"RefMsgType,omitempty"`          // Tag 372
	SessionRejectReason string `json:"SessionRejectReason,omitempty"` // Tag 373: Contains coded values for rejection reasons
}

type CtraderMessageChannel int
type CtraderSessionMessageType int
type CtraderQualifier string
type CtraderParamIds int

const (
	QUOTE CtraderMessageChannel = iota
	TRADE
)

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
	MarketDataRequest
	RequestForCancelOrder
)

const (
	EncryptionMethod  CtraderParamIds = 98
	HeartbeatInterval CtraderParamIds = 108
	ResetSequence     CtraderParamIds = 141
	Username          CtraderParamIds = 553
	Password          CtraderParamIds = 554

	//35 = V
	MDReqID                 CtraderParamIds = 262
	SubscriptionRequestType CtraderParamIds = 263
	MarketDepth             CtraderParamIds = 264
	MDUpdateType            CtraderParamIds = 265
	NoMDEntryType           CtraderParamIds = 267
	MDEntryType             CtraderParamIds = 269
	NoRelatedSym            CtraderParamIds = 146

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

	SubscriptionRequestType: {"subscribe": "1", "unsubscribe": "2"},
	MarketDepth:             {"depth": "0", "spot": "1"},
	MDUpdateType:            {"incrementalRefresh": "1"}, // don't know what other values can go here, just grabbed from example
	NoRelatedSym:            {},
	MDEntryType:             {"bid": "0", "offer": "1"},

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
	HeaderMessageType:           {"0": "A", "1": "5", "2": "0", "3": "1", "4": "2", "5": "3", "6": "4", "7": "H", "8": "AF", "9": "AN", "10": "D", "11": "x", "12": "V"},
	HeaderSenderCompId:          {},
	HeaderTargetCompId:          {"compId": "CSERVER"},
	HeaderTargetSubId:           {"trade": "TRADE", "quote": "QUOTE"},
	HeaderSenderSubId:           {"trade": "TRADE", "quote": "QUOTE"},
	HeaderMessageSequenceNumber: {},
	HeaderMessageTimestamp:      {},

	//Trailer
	TrailerChecksum: {},
}

type ResponseType struct {
	Name     string
	IsError  bool
	IsReject bool
}

type MarketDataSubscription struct {
	MDReqID        string // uuid
	Action         string // subscribe, unsubscribe
	MarketDepth    string // depth, spot
	MDUpdateType   string // always "incrementalRefresh"
	NoMDEntryTypes int    //always 2, will get bids and asks
	MDEntryType    []int  //always {0,1} //both sent with 269 tag
	NoRelatedSym   int    //always 1
	Symbol         string // Could have multiple symbols per request
}

type Position struct {
	PID     string
	CopyPID string
	Side    string
	Symbol  string
	AvgPx   float64
}
