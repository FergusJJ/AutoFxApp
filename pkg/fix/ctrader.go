package fix

import (
	"fmt"
	"log"
	"strings"
)

func (session *FxSession) CtraderLogin(user FxUser, channel CtraderMessageChannel) *ErrorWithCause {
	var fxResponseMap []*FixResponse

	loginMessage, err := user.constructLogin(session, channel)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	if channel == QUOTE {
		fxResponseMap, err = session.PriceClient.RoundTrip(loginMessage)

	}
	if channel == TRADE {
		fxResponseMap, err = session.TradeClient.RoundTrip(loginMessage)

	}

	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   CtraderConnectionError,
		}
	}
	if channel == QUOTE {
		session.PriceMessageSequenceNumber++
	} else {
		session.TradeMessageSequenceNumber++
	}

	if len(fxResponseMap) != 1 {
		return &ErrorWithCause{
			ErrorMessage: fmt.Sprintf("ctrader login, unexpected response length: %d", len(fxResponseMap)),
			ErrorCause:   ProgramError,
		}
	}
	success_, err := ParseFixResponse(fxResponseMap[0], Logon)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	success, ok := success_.(bool)
	if !ok {
		log.Fatalf("cannot convert interface: %v to bool", success_)
	}
	if success {
		return nil
	}
	return &ErrorWithCause{
		ErrorMessage: "unable to login",
		ErrorCause:   UserDataError,
	}

}

// func (session *FxSession) CtraderSecurityList(user FxUser) *ErrorWithCause {
// 	// securityRequestID :=  //"Sxo2Xlb1jzJB" //idk whether this has to be different between users?
// 	securityMessage, err := user.constructSecurityList(session, "Sxo2Xlb1jzJC")
// 	fmt.Println(securityMessage)
// 	if err != nil {
// 		return &ErrorWithCause{
// 			ErrorMessage: err.Error(),
// 			ErrorCause:   ProgramError,
// 		}
// 	}

// 	resp := session.sendMessage(securityMessage, user)
// 	if resp.err != nil {
// 		return &ErrorWithCause{
// 			ErrorMessage: err.Error(),
// 			ErrorCause:   ConnectionError,
// 		}
// 	}

// 	_, err = ParseFIXResponse(resp.body, SecurityListRequest)
// 	if err != nil {
// 		return &ErrorWithCause{
// 			ErrorMessage: err.Error(),
// 			ErrorCause:   ProgramError,
// 		}
// 	}
// 	bodyStringSlice := preparseBody(resp.body)

// 	secListPairs, err := parseSecurityList(string(resp.body))
// 	if err != nil {
// 		return &ErrorWithCause{
// 			ErrorMessage: err.Error(),
// 			ErrorCause:   ProgramError,
// 		}
// 	}
// 	session.SecListPairs = secListPairs
// 	session.MessageSequenceNumber++
// 	return nil

// }

// newOrderData := idos.OrderSingleData{OrderQty: "1000", Symbol: "1", Side: "buy", OrdType: "market"}
func (session *FxSession) CtraderNewOrderSingle(user FxUser, orderData OrderData) (ExecutionReport, *ErrorWithCause) {

	orderMessage, err := user.constructNewOrderSingle(session, orderData)
	if err != nil {
		return ExecutionReport{}, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}
	fxResponseMap, err := session.TradeClient.RoundTrip(orderMessage)
	if err != nil {
		return ExecutionReport{}, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   CtraderConnectionError,
		}
	}
	session.TradeMessageSequenceNumber++
	executionReports := []ExecutionReport{}
	for _, v := range fxResponseMap {
		fxRes, err := ParseFixResponse(v, NewOrderSingle)
		if err != nil {
			if strings.Contains(err.Error(), "business reject") {
				businessRejectRes, ok := fxRes.(BusinessMessageReject)
				if !ok {
					log.Fatalf("cannot cast %v to businessMessageReject", fxRes)
				}
				return ExecutionReport{}, &ErrorWithCause{
					ErrorCause:   MarketError,
					ErrorMessage: businessRejectRes.Text,
				}
			}
		}

		executionReport, ok := fxRes.(ExecutionReport)
		if !ok {
			log.Fatalf("cannot cast %v to executionReport", fxRes)
		}
		executionReports = append(executionReports, executionReport)
	}

	return executionReports[0], nil
}

// // Needs clordID
// func (session *FxSession) CtraderOrderStatus(user *FxUser) {

// }

// might want to return a mapping here, then can check the 911 tag of the first item to see the number of reports if it is needed
func (session *FxSession) CtraderOrderStatus(user FxUser, clOrdID string) (ExecutionReport, *ErrorWithCause) {
	statusMessage, err := user.constructOrderStatusRequest(session, clOrdID)
	if err != nil {
		return ExecutionReport{}, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}

	}
	fxResponseMap, err := session.TradeClient.RoundTrip(statusMessage)
	if err != nil {
		return ExecutionReport{}, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   CtraderConnectionError,
		}
	}
	session.TradeMessageSequenceNumber++
	if len(fxResponseMap) != 1 {
		log.Fatalf("order status: resp map len %d: %+v", len(fxResponseMap), fxResponseMap)
	}
	fxRes, err := ParseFixResponse(fxResponseMap[0], NewOrderSingle)
	if err != nil {
		if strings.Contains(err.Error(), "business reject") {
			businessRejectRes, ok := fxRes.(BusinessMessageReject)
			if !ok {
				log.Fatalf("cannot cast %v to businessMessageReject", fxRes)
			}
			return ExecutionReport{}, &ErrorWithCause{
				ErrorCause:   MarketError,
				ErrorMessage: businessRejectRes.Text,
			}
		}
	}

	executionReport, ok := fxRes.(ExecutionReport)
	if !ok {
		log.Fatalf("cannot cast %v to executionReport", fxRes)
	}
	return executionReport, nil
}

func (session *FxSession) CtraderRequestForPositions(user FxUser) ([]PositionReport, *ErrorWithCause) {

	positionsMessage, err := user.constructPositionsRequest(session) //constructOrderMassStatusRequest(session)
	if err != nil {
		return []PositionReport{}, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}

	}
	fxResponseMap, err := session.TradeClient.RoundTrip(positionsMessage)
	if err != nil {
		return nil, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   CtraderConnectionError,
		}
	}
	session.TradeMessageSequenceNumber++
	var positions = make([]PositionReport, 0)
	for _, message := range fxResponseMap {
		//does not currently return an error so not gonna have proper handling rn
		fxRes, err := ParseFixResponse(message, RequestForPositions)
		if err != nil {
			log.Fatalf("error getting positions: %+v", err)
		}
		positionReport, ok := fxRes.(PositionReport)
		if !ok {
			log.Fatalf("cannot cast %v to positionReport", fxRes)
		}
		positions = append(positions, positionReport)
	}

	return positions, nil
}

func (session *FxSession) CtraderMarketDataRequest(user FxUser, subscription MarketDataSubscription) ([]MarketDataSnapshot, *ErrorWithCause) {

	marketDataRequestMessage, err := user.constructMarketDataRequest(session, subscription)
	if err != nil {
		return nil, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}
	fxResponseMap, err := session.PriceClient.RoundTrip(marketDataRequestMessage)
	if err != nil {
		return nil, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   CtraderConnectionError,
		}
	}
	session.PriceMessageSequenceNumber++
	var data = make([]MarketDataSnapshot, 0)
	for _, message := range fxResponseMap {
		//does not currently return an error so not gonna have proper handling rn
		fxRes, err := ParseFixResponse(message, MarketDataRequest)
		if err != nil {
			_, ok := fxRes.(BusinessMessageReject)
			if !ok {
				log.Fatalf("error getting marketData: %+v", err)
			}
			return data, &ErrorWithCause{
				ErrorCause:   MarketError,
				ErrorMessage: err.Error(),
			}
		}
		marketDataSnapshot, ok := fxRes.(MarketDataSnapshot)
		if !ok {
			log.Fatalf("cannot cast %v to marketDataSnapshot", fxRes)
		}
		data = append(data, marketDataSnapshot)
	}

	return data, nil

}

/*

doc:
34=3|50=QUOTE|263=1|264=1|265=1|146=1|55=1|267=2|269=0|269=1|10=094|


go
57=TRADE|50=TRADE|34=2|263=1|264=1|265=1|267=2|269=1|146=1|55=1|


doc -
57=TRADE
50=TRADE

doc +
50=QUOTE

*/
