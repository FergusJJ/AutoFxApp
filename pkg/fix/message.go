package fix

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (user *FxUser) constructSecurityList(session *FxSession, securityRequestID string) (string, error) {
	var securityListBody string
	var securityListParams []string
	securityListParams = append(securityListParams, formatMessageSlice(SecurityReqID, securityRequestID, true))
	securityListParams = append(securityListParams, formatMessageSlice(SecurityListRequestType, "0", true))
	securityListBody = strings.Join(securityListParams, "|")
	securityListBody = fmt.Sprintf("%s|", securityListBody)
	header := user.constructHeader(securityListBody, SecurityListRequest, session, QUOTE)
	headerWithBody := fmt.Sprintf("%s%s", header, securityListBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructLogin(session *FxSession, channel CtraderMessageChannel) (string, error) {
	var loginBody string
	var loginParams []string
	compIdSlice := strings.Split(user.SenderCompID, ".")
	if len(compIdSlice) == 1 {
		return "", errors.New(`"senderCompID" is incorrect`)
	}
	username := compIdSlice[len(compIdSlice)-1]

	loginParams = append(loginParams, formatMessageSlice(EncryptionMethod, "encryptionDisabled", false))

	loginParams = append(loginParams, formatMessageSlice(HeartbeatInterval, "0", true))

	loginParams = append(loginParams, formatMessageSlice(ResetSequence, "resetDisabled", false))

	loginParams = append(loginParams, formatMessageSlice(Username, username, true))
	loginParams = append(loginParams, formatMessageSlice(Password, user.Password, true))
	loginBody = strings.Join(loginParams, "|")
	loginBody = fmt.Sprintf("%s|", loginBody)

	header := user.constructHeader(loginBody, Logon, session, channel)
	headerWithBody := fmt.Sprintf("%s%s", header, loginBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructNewOrderSingle(session *FxSession, orderData OrderData) (string, error) {
	var newOrderSingleBody string
	var newOrderSingleParams []string
	orderData.OrderType = "market"

	volAsString := fmt.Sprintf("%g", orderData.Volume)
	transactTime := time.Now().UTC().Format(YYYYMMDDhhmmss)
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(ClOrdID, uuid.New().String(), true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSSymbol, orderData.Symbol, true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(Side, orderData.Direction, false))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSTransactTime, transactTime, true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSOrderQty, volAsString, true))
	newOrderSingleParams = append(newOrderSingleParams, formatMessageSlice(NOSOrdType, orderData.OrderType, false))
	newOrderSingleBody = strings.Join(newOrderSingleParams, "|")
	newOrderSingleBody = fmt.Sprintf("%s|", newOrderSingleBody)
	header := user.constructHeader(newOrderSingleBody, NewOrderSingle, session, TRADE)
	headerWithBody := fmt.Sprintf("%s%s", header, newOrderSingleBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

// func (user *FxUser) constructOrderMassStatusRequest(session *FxSession) (string, error) {
// 	var orderMassStatusRequestBody string
// 	var orderMassStatusRequestParams []string
// 	orderMassStatusRequestParams = append(orderMassStatusRequestParams, formatMessageSlice(MassStatusReqID, uuid.New().String(), true))
// 	orderMassStatusRequestParams = append(orderMassStatusRequestParams, formatMessageSlice(MassStatusReqType, "7", true))
// 	orderMassStatusRequestBody = strings.Join(orderMassStatusRequestParams, "|")
// 	orderMassStatusRequestBody = fmt.Sprintf("%s|", orderMassStatusRequestBody)
// 	header := user.constructHeader(orderMassStatusRequestBody, OrderMassStatusRequest, session)
// 	headerWithBody := fmt.Sprintf("%s%s", header, orderMassStatusRequestBody)
// 	trailer := constructTrailer(headerWithBody)
// 	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
// 	return message, nil
// }

func (user *FxUser) constructOrderStatusRequest(session *FxSession, clOrdId string) (string, error) {
	var orderStatusRequestBody string
	var orderStatusRequestParams []string
	orderStatusRequestParams = append(orderStatusRequestParams, formatMessageSlice(ClOrdID, clOrdId, true))
	orderStatusRequestBody = strings.Join(orderStatusRequestParams, "|")
	orderStatusRequestBody = fmt.Sprintf("%s|", orderStatusRequestBody)
	header := user.constructHeader(orderStatusRequestBody, OrderStatusRequest, session, TRADE)
	headerWithBody := fmt.Sprintf("%s%s", header, orderStatusRequestBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructPositionsRequest(session *FxSession) (string, error) {
	var constructPositionsRequestBody string
	var constructPositionsRequestParams []string
	constructPositionsRequestParams = append(constructPositionsRequestParams, formatMessageSlice(PosReqID, uuid.New().String(), true))
	constructPositionsRequestBody = strings.Join(constructPositionsRequestParams, "|")
	constructPositionsRequestBody = fmt.Sprintf("%s|", constructPositionsRequestBody)
	header := user.constructHeader(constructPositionsRequestBody, RequestForPositions, session, TRADE)
	headerWithBody := fmt.Sprintf("%s%s", header, constructPositionsRequestBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	return message, nil
}

func (user *FxUser) constructMarketDataRequest(session *FxSession, subscription MarketDataSubscription) (string, error) {
	var constructMarketDataRequestBody string
	var constructMarketDataRequestParams []string
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(MDReqID, subscription.MDReqID, true))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(SubscriptionRequestType, subscription.Action, false))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(MarketDepth, subscription.MarketDepth, false))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(MDUpdateType, subscription.MDUpdateType, false))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(NoMDEntryType, fmt.Sprint(subscription.NoMDEntryTypes), true))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(MDEntryType, "0", true))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(MDEntryType, "1", true))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(NoRelatedSym, fmt.Sprint(subscription.NoRelatedSym), true))
	constructMarketDataRequestParams = append(constructMarketDataRequestParams, formatMessageSlice(NOSSymbol, subscription.Symbol, true))

	constructMarketDataRequestBody = strings.Join(constructMarketDataRequestParams, "|")
	constructMarketDataRequestBody = fmt.Sprintf("%s|", constructMarketDataRequestBody)
	header := user.constructHeader(constructMarketDataRequestBody, MarketDataRequest, session, QUOTE)
	headerWithBody := fmt.Sprintf("%s%s", header, constructMarketDataRequestBody)
	trailer := constructTrailer(headerWithBody)
	message := strings.ReplaceAll(fmt.Sprintf("%s%s", headerWithBody, trailer), "|", "\u0001")
	// log.Print(headerWithBody)
	return message, nil
}

func (user *FxUser) constructHeader(bodyMessage string, messageType CtraderSessionMessageType, session *FxSession, channel CtraderMessageChannel) string {
	var messageTypeStr = fmt.Sprintf("%d", messageType)
	var header string
	var headerParams []string
	var messageSequenceString string
	if channel == QUOTE {
		messageSequenceString = strconv.Itoa(session.PriceMessageSequenceNumber)
	} else {
		messageSequenceString = strconv.Itoa(session.TradeMessageSequenceNumber)
	}
	messageTs := time.Now().UTC().Format(YYYYMMDDhhmmss)

	headerParams = append(headerParams, formatMessageSlice(HeaderMessageType, messageTypeStr, false))
	headerParams = append(headerParams, formatMessageSlice(HeaderSenderCompId, user.SenderCompID, true))
	headerParams = append(headerParams, formatMessageSlice(HeaderTargetCompId, user.TargetCompID, true))
	if channel == QUOTE {
		headerParams = append(headerParams, formatMessageSlice(HeaderTargetSubId, "quote", false))
	}
	if channel == TRADE {
		headerParams = append(headerParams, formatMessageSlice(HeaderTargetSubId, "trade", false))
	}
	headerParams = append(headerParams, formatMessageSlice(HeaderMessageSequenceNumber, messageSequenceString, true))
	headerParams = append(headerParams, formatMessageSlice(HeaderMessageTimestamp, messageTs, true))
	header = strings.Join(headerParams, "|")
	messageLength := strconv.Itoa(len(bodyMessage) + len(header) + 1) // +1 is to account for the missing "|"

	headerParams = append([]string{formatMessageSlice(HeaderBeginString, "begin", false), formatMessageSlice(HeaderMessageLength, messageLength, true)}, headerParams...)
	header = fmt.Sprintf("%s|%s", strings.Join(headerParams, "|"), "")
	return header
}

func constructTrailer(message string) (trailer string) {
	checksumInput := strings.ReplaceAll(message, "|", "\u0001")
	checksum := strconv.Itoa(calculateChecksum(checksumInput))
	checksum = func(checksum string) string {
		if len(checksum) == 0 {
			return "000"
		}
		if len(checksum) == 1 {
			return fmt.Sprintf("00%s", checksum)
		}
		if len(checksum) == 2 {
			return fmt.Sprintf("0%s", checksum)
		}
		return checksum
	}(checksum)
	trailer = fmt.Sprintf("10=%s\u0001", checksum)
	return trailer
}

func calculateChecksum(dataToCalculate string) int {
	byteToCalculate := []byte(dataToCalculate)
	checksum := 0
	for _, chData := range byteToCalculate {
		checksum += int(chData)
	}
	return checksum % 256
}

func formatMessageSlice(ids CtraderParamIds, value string, useValueAsValue bool) string {
	if useValueAsValue {
		return fmt.Sprintf("%d=%s", ids, value)
	}
	return fmt.Sprintf("%d=%s", ids, MessageKeyValuePairs[ids][value])
}

func GetUUID() string {
	return uuid.New().String()
}

func parseFixResponse(resp *FixResponse, reqType CtraderSessionMessageType) (interface{}, error) {
	if resp.MsgType == "3" {
		//session level violation, just want the reason in this case.
		var sessionRejectMessageMapping = map[string]string{}
		var sessionRejectMessage SessionRejectMessage
		for tag, val := range resp.Body {
			sessionRejectMessageMapping[sessionRejectMessageTagMapping[tag]] = val
		}
		jsonData, err := json.Marshal(sessionRejectMessageMapping)
		if err != nil {
			log.Fatal(err)
		}
		err = json.Unmarshal(jsonData, &sessionRejectMessage)
		if err != nil {
			log.Fatal(err)
		}
		log.Fatalf("%+v", sessionRejectMessage)
		return sessionRejectMessage, fmt.Errorf("session reject message: %s", sessionRejectMessage.Text)
	}
	switch reqType {
	case Logon, Logout:
		switch resp.MsgType {
		case "A": //success
			return true, nil
		case "5": //session is logged out
			if reqType == Logon {
				reason, ok := resp.Body[58]
				if !ok {
					log.Fatal("no reason tag for failed login")
				}
				if strings.Contains(reason, "RET_INVALID_DATA") {
					reason = "invalid user data"
				}
				return false, fmt.Errorf("unable to login, %s", reason)
			}
			return false, nil
		default:
			log.Fatalf("unhandled resp.MsgType %s. reqType: %d", resp.MsgType, reqType)
		}
	case MarketDataRequest:
		switch resp.MsgType {
		case "W": // for spots, snapshot, full refresh
			var marketDataSnapshotMapping = map[string]string{}
			var marketDataSnapshot MarketDataSnapshot
			for tag, val := range resp.Body {
				marketDataSnapshotMapping[marketDataSnapshotTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(marketDataSnapshotMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &marketDataSnapshot)
			if err != nil {
				log.Fatal(err)
			}
			return marketDataSnapshot, nil

		case "X": // for depths inc refresh
			var marketDataIncrementalRefreshMapping = map[string]string{}
			var marketDataIncrementalRefresh MarketDataIncrementalRefresh
			for tag, val := range resp.Body {
				marketDataIncrementalRefreshMapping[marketDataIncrementalRefreshTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(marketDataIncrementalRefreshMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &marketDataIncrementalRefresh)
			if err != nil {
				log.Fatal(err)
			}
			return marketDataIncrementalRefresh, nil
		case "Y": //market data request reject (just bad symbol, or message)
			var marketDataRequestRejectMapping = map[string]string{}
			var marketDataRequestReject MarketDataSnapshot
			for tag, val := range resp.Body {
				marketDataRequestRejectMapping[marketDataRequestRejectTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(marketDataRequestRejectMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &marketDataRequestReject)
			if err != nil {
				log.Fatal(err)
			}
			return marketDataRequestReject, fmt.Errorf("market data request reject")
		case "j":
			var businessMessageRejectMapping = map[string]string{}
			var businessMessageReject BusinessMessageReject
			for tag, val := range resp.Body {
				businessMessageRejectMapping[businessMessageRejectTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(businessMessageRejectMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &businessMessageReject)
			if err != nil {
				log.Fatal(err)
			}
			return businessMessageReject, fmt.Errorf("business reject: %s", businessMessageReject.Text)
		default:
			log.Fatalf("unhandled resp.MsgType %s. reqType: %d", resp.MsgType, reqType)
		}
	case NewOrderSingle, OrderStatusRequest, OrderMassStatusRequest, RequestForCancelOrder:
		switch resp.MsgType {
		case "8": //success
			var executionReportMapping = map[string]string{}
			var executionReport ExecutionReport
			//execution report, order has gone through
			for tag, val := range resp.Body {
				executionReportMapping[executionReportTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(executionReportMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &executionReport)
			if err != nil {
				log.Fatal(err)
			}
			return executionReport, nil
		case "j": //fail
			var businessMessageRejectMapping = map[string]string{}
			var businessMessageReject BusinessMessageReject
			for tag, val := range resp.Body {
				businessMessageRejectMapping[businessMessageRejectTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(businessMessageRejectMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &businessMessageReject)
			if err != nil {
				log.Fatal(err)
			}
			return businessMessageReject, fmt.Errorf("business reject") // will have to check for this error when "j" possible and print specific info
		default:
			log.Fatalf("unhandled resp.MsgType %s. reqType: %d", resp.MsgType, reqType)
		}
	case RequestForPositions:
		switch resp.MsgType {
		case "AP":
			var positionReportMapping = map[string]string{}
			var positionReport PositionReport
			for tag, val := range resp.Body {
				positionReportMapping[positionReportTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(positionReportMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &positionReport)
			if err != nil {
				log.Fatal(err)
			}
			return positionReport, nil
		default:
			log.Fatalf("unhandled resp.MsgType %s. reqType: %d", resp.MsgType, reqType)
		}
	default:
		log.Fatalf("no case specified for reqType: %+v\nfix response:\n%+v", reqType, resp)
	}
	return nil, nil
}

/*
8=FIX.4.4|9=166|35=V|49=demo.ctrader.3697899|56=CServer|34=2|52=20230808-19:26:37|262=3d00357c-0ea2-4379-b82a-04387f880071|263=1|264=1|265=1|267=2|269=1|146=1|55=1|

8=FIX.4.4|9=131|35=V|49=live.theBroker.12345|56=CSERVER|34=3|52=20170117-10:26:54|262=876316403|263=1|264=1|265=1|146=1|55=1|267=2|269=0|269=1|10=094|


*/
