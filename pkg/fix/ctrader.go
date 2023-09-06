package fix

import (
	"fmt"
	"log"
	"reflect"
)

func (session *FxSession) CtraderLogin(user FxUser, channel CtraderMessageChannel) *CtraderError {
	var fxResponseMap []*FixResponse

	loginMessage, err := user.constructLogin(session, channel)
	if err != nil {
		return &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  err,
			ShouldExit:  true,
		}
	}
	if channel == QUOTE {
		fxResponseMap, err = session.PriceClient.RoundTrip(loginMessage)

	}
	if channel == TRADE {
		fxResponseMap, err = session.TradeClient.RoundTrip(loginMessage)

	}

	if err != nil {
		return &CtraderError{
			UserMessage: "An error occurred whilst logging in",
			ErrorType:   CtraderConnectionError,
			ErrorCause:  err,
			ShouldExit:  false,
		}
	}
	if channel == QUOTE {
		session.PriceMessageSequenceNumber++
	} else {
		session.TradeMessageSequenceNumber++
	}

	if len(fxResponseMap) != 1 {
		return &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  fmt.Errorf("ctrader login, unexpected response length: %d", len(fxResponseMap)),
			ShouldExit:  true,
		}
	}
	success_, ctErr := ParseFixResponse(fxResponseMap[0], Logon)
	if ctErr != nil {
		return ctErr
	}
	success, ok := success_.(bool)
	if !ok {
		//then we always have session reject
		rejectMsg, ok := success_.(SessionRejectMessage)
		if !ok {
			return &CtraderError{
				UserMessage: "An unexpected error occurred whilst handling another error",
				ErrorType:   CtraderLogicError,
				ErrorCause:  fmt.Errorf("unable to convert interface to SessionRejectMessage"),
				ShouldExit:  true,
			}
		}
		return ErrorFromSessionReject(rejectMsg)
	}
	if success {
		return nil
	}
	return &CtraderError{
		UserMessage: "Unable to login",
		ErrorType:   CtraderLogicError,
		ErrorCause:  fmt.Errorf("unable to convert interface to SessionRejectMessage"),
		ShouldExit:  true,
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
func (session *FxSession) CtraderNewOrderSingle(user FxUser, orderData OrderData) (*ExecutionReport, *CtraderError) {

	orderMessage, err := user.constructNewOrderSingle(session, orderData)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  err,
			ShouldExit:  true,
		}
	}
	fxResponseMap, err := session.TradeClient.RoundTrip(orderMessage)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An error occurred sending an order",
			ErrorType:   CtraderConnectionError,
			ErrorCause:  err,
			ShouldExit:  false,
		}
	}
	session.TradeMessageSequenceNumber++
	executionReports := []ExecutionReport{}
	for _, v := range fxResponseMap {
		fxRes, err := ParseFixResponse(v, NewOrderSingle)
		if err != nil {
			return nil, err
		}
		executionReport, ok := fxRes.(ExecutionReport)
		if !ok {
			BusinessMessageReject, ok := fxRes.(BusinessMessageReject)
			if !ok {
				rejectMsg, ok := fxRes.(SessionRejectMessage)
				if !ok {
					return nil, &CtraderError{
						UserMessage: "An unexpected error occurred whilst handling another error",
						ErrorType:   CtraderLogicError,
						ErrorCause:  fmt.Errorf("unable to convert interface to SessionRejectMessage"),
						ShouldExit:  true,
					}
				}
				return nil, ErrorFromSessionReject(rejectMsg)
			}
			//handle reject and return error
			return nil, &CtraderError{
				UserMessage: "An order couldn't be processed",
				ErrorType:   CtraderBusinessRejectError,
				ErrorCause:  fmt.Errorf("business reject: new order single: %s", BusinessMessageReject.Text),
				ShouldExit:  true,
			}
		}
		executionReports = append(executionReports, executionReport)
	}
	if len(executionReports) > 1 {
		//maybe just log something to client for now? Is unlikely that will happen with market order I think.
		//or could add to some sort of messageQueue?
		for _, v := range executionReports {
			log.Printf("%+v\n", v)
		}
		log.Fatal("unexpected response length")
	}
	return &executionReports[0], nil
}

// // Needs clordID
// func (session *FxSession) CtraderOrderStatus(user *FxUser) {

// }

// might want to return a mapping here, then can check the 911 tag of the first item to see the number of reports if it is needed
func (session *FxSession) CtraderOrderStatus(user FxUser, clOrdID string) (*ExecutionReport, *CtraderError) {
	statusMessage, err := user.constructOrderStatusRequest(session, clOrdID)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  err,
			ShouldExit:  true,
		}

	}
	fxResponseMap, err := session.TradeClient.RoundTrip(statusMessage)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An error occurred whilst getting an order",
			ErrorType:   CtraderConnectionError,
			ErrorCause:  err,
			ShouldExit:  false,
		}
	}

	session.TradeMessageSequenceNumber++
	if len(fxResponseMap) != 1 {
		return nil, &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  fmt.Errorf("order status: resp map len %d: %+v", len(fxResponseMap), fxResponseMap),
			ShouldExit:  true,
		}
	}
	fxRes, ctErr := ParseFixResponse(fxResponseMap[0], OrderStatusRequest)
	if ctErr != nil {
		return nil, ctErr
	}
	executionReport, ok := fxRes.(ExecutionReport)
	if ok {
		return &executionReport, nil
	}

	BusinessMessageReject, ok := fxRes.(BusinessMessageReject)
	if ok {
		return nil, &CtraderError{
			UserMessage: "An order couldn't be processed",
			ErrorType:   CtraderBusinessRejectError,
			ErrorCause:  fmt.Errorf("business reject: new order single: %s", BusinessMessageReject.Text),
			ShouldExit:  true,
		}
	}

	rejectMsg, ok := fxRes.(SessionRejectMessage)
	if ok {
		return nil, ErrorFromSessionReject(rejectMsg)
	}

	return nil, &CtraderError{
		UserMessage: "An unexpected error occurred whilst handling another error",
		ErrorType:   CtraderLogicError,
		ErrorCause:  fmt.Errorf("unable to convert interface type: %s", reflect.TypeOf(fxRes)),
		ShouldExit:  true,
	}
	//handle reject and return error

}

func (session *FxSession) CtraderRequestForPositions(user FxUser) ([]PositionReport, *CtraderError) {

	positionsMessage, err := user.constructPositionsRequest(session) //constructOrderMassStatusRequest(session)
	if err != nil {
		return []PositionReport{}, &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  err,
			ShouldExit:  true,
		}

	}
	fxResponseMap, err := session.TradeClient.RoundTrip(positionsMessage)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An error occurred whilst fetching positions",
			ErrorType:   CtraderConnectionError,
			ErrorCause:  err,
			ShouldExit:  false,
		}
	}
	session.TradeMessageSequenceNumber++
	var positions = make([]PositionReport, 0)
	for _, message := range fxResponseMap {
		//does not currently return an error so not gonna have proper handling rn
		fxRes, ctErr := ParseFixResponse(message, RequestForPositions)
		if ctErr != nil {
			return []PositionReport{}, ctErr
		}
		positionReport, ok := fxRes.(PositionReport)
		if !ok {
			rejectMsg, ok := fxRes.(SessionRejectMessage)
			if !ok {
				return nil, &CtraderError{
					UserMessage: "An unexpected error occurred whilst handling another error",
					ErrorType:   CtraderLogicError,
					ErrorCause:  fmt.Errorf("unable to convert interface to SessionRejectMessage"),
					ShouldExit:  true,
				}
			}
			return nil, ErrorFromSessionReject(rejectMsg)
		}
		positions = append(positions, positionReport)
	}

	return positions, nil
}

func (session *FxSession) CtraderMarketDataRequest(user FxUser, subscription MarketDataSubscription) ([]MarketDataSnapshot, *CtraderError) {

	marketDataRequestMessage, err := user.constructMarketDataRequest(session, subscription)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An unexpected error occurred",
			ErrorType:   CtraderLogicError,
			ErrorCause:  err,
			ShouldExit:  true,
		}
	}
	fxResponseMap, err := session.PriceClient.RoundTrip(marketDataRequestMessage)
	if err != nil {
		return nil, &CtraderError{
			UserMessage: "An error occurred whilst fetching market data",
			ErrorType:   CtraderConnectionError,
			ErrorCause:  err,
			ShouldExit:  false,
		}
	}
	session.PriceMessageSequenceNumber++
	var data = make([]MarketDataSnapshot, 0)
	for _, message := range fxResponseMap {
		//does not currently return an error so not gonna have proper handling rn
		fxRes, err := ParseFixResponse(message, MarketDataRequest)
		if err != nil {
			return []MarketDataSnapshot{}, err
		}

		marketDataSnapshot, ok := fxRes.(MarketDataSnapshot)
		if !ok {
			if marketDataRequestReject, ok := fxRes.(MarketDataRequestReject); ok {
				return []MarketDataSnapshot{}, &CtraderError{
					UserMessage: "An error occurred whilst fetching market data for a symbol",
					ErrorType:   CtraderBusinessRejectError,
					ErrorCause:  fmt.Errorf("business reject: market data snapshot: %d", marketDataRequestReject.MDReqRejReason),
					ShouldExit:  false,
				}
			}
			if businessMessageReject, ok := fxRes.(BusinessMessageReject); ok {
				return []MarketDataSnapshot{}, &CtraderError{
					UserMessage: "An error occurred whilst fetching market data",
					ErrorType:   CtraderBusinessRejectError,
					ErrorCause:  fmt.Errorf("business reject: market data snapshot: %s", businessMessageReject.Text),
					ShouldExit:  true,
				}
			}

			if _, ok = fxRes.(MarketDataIncrementalRefresh); ok {
				return []MarketDataSnapshot{}, &CtraderError{
					UserMessage: "An unexpected error occurred",
					ErrorType:   CtraderLogicError,
					ErrorCause:  fmt.Errorf("got MarketDataIncrementalRefresh"),
					ShouldExit:  true,
				}
			}
			rejectMsg, ok := fxRes.(SessionRejectMessage)
			if ok {
				return []MarketDataSnapshot{}, ErrorFromSessionReject(rejectMsg)
			}

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
