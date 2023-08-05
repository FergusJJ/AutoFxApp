package fix

import (
	"fmt"
	"log"
)

func (session *FxSession) CtraderLogin(user FxUser) *ErrorWithCause {
	loginMessage, err := user.constructLogin(session)
	if err != nil {

		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}
	resp := session.sendMessage(loginMessage, user)
	if resp.err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ConnectionError,
		}
	}
	bodyStringSlice := preparseBody(resp.body)
	_, _, err = ParseFIXResponse(bodyStringSlice[0], Logon)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}

	session.MessageSequenceNumber++
	return nil
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

/* NEED EXTRA PARAMS*/

// newOrderData := idos.OrderSingleData{OrderQty: "1000", Symbol: "1", Side: "buy", OrdType: "market"}
func (session *FxSession) CtraderNewOrderSingle(user FxUser, orderData OrderData) (*ExecutionReport, *ErrorWithCause) {

	orderMessage, err := user.constructNewOrderSingle(session, orderData)
	if err != nil {
		return nil, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}
	resp := session.sendMessage(orderMessage, user)
	if resp.err != nil {
		return nil, &ErrorWithCause{
			ErrorMessage: resp.err.Error(),
			ErrorCause:   ConnectionError,
		}
	}
	bodyStringSlice := preparseBody(resp.body)
	respType, parsedResp, err := ParseFIXResponse(bodyStringSlice[0], NewOrderSingle)
	if err != nil {
		return nil, &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	switch respType {
	case "8":
		executionReport, ok := parsedResp.(ExecutionReport)
		if !ok {
			log.Fatal(`error casting struct with respType of "8" to ExecutionReport`)
		}
		return &executionReport, nil

	case "9":
		orderCancelReject, ok := parsedResp.(OrderCancelReject)
		if !ok {
			log.Fatal(`error casting struct with respType of "9" to OrderCancelReject`)
		}
		log.Printf("%+v", orderCancelReject)
		return nil, &ErrorWithCause{
			ErrorMessage: "Order Cancel Reject",
			ErrorCause:   MarketError,
		}
	case "j":
		businessMessageReject, ok := parsedResp.(BusinessMessageReject)
		if !ok {
			log.Fatal(`error casting struct with respType of "j" to BusinessMessageReject`)
		}
		log.Printf("%+v", businessMessageReject)
		return nil, &ErrorWithCause{
			ErrorMessage: "Business Message Reject",
			ErrorCause:   ProgramError,
		}

	}
	session.MessageSequenceNumber++
	return nil, &ErrorWithCause{
		ErrorMessage: fmt.Sprintf("respType fallthrough %s", respType),
		ErrorCause:   ProgramError,
	}
}

// Needs clordID
func (session *FxSession) CtraderOrderStatus(user *FxUser) {

}

func (session *FxSession) CtraderMassStatus(user FxUser) *ErrorWithCause {
	statusMessage, err := user.constructOrderMassStatusRequest(session)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}

	}
	resp := session.sendMessage(statusMessage, user)
	if resp.err != nil {
		return &ErrorWithCause{
			ErrorMessage: resp.err.Error(),
			ErrorCause:   ConnectionError,
		}
	}

	bodyStringSlice := preparseBody(resp.body)

	_, _, err = ParseFIXResponse(bodyStringSlice[0], OrderMassStatusRequest)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	session.MessageSequenceNumber++
	return nil
}

func (session *FxSession) CtraderRequestForPositions(user FxUser) *ErrorWithCause {

	orderMessage, err := user.constructPositionsRequest(session) //constructOrderMassStatusRequest(session)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}

	}
	resp := session.sendMessage(orderMessage, user)
	if resp.err != nil {
		return &ErrorWithCause{
			ErrorMessage: resp.err.Error(),
			ErrorCause:   ConnectionError,
		}
	}
	//want to pre parse in case more than 1 position is present
	positionReports := []PositionReport{}
	bodyStringSlice := preparseBody(resp.body)
	for _, v := range bodyStringSlice {
		_, positionReport, err := ParseFIXResponse(v, RequestForPositions)
		if err != nil {
			return &ErrorWithCause{
				ErrorMessage: err.Error(),
				ErrorCause:   UserDataError,
			}
		}
		positionReportCast, ok := positionReport.(PositionReport)
		if !ok {
			log.Fatalf("couldnt cast: \n%+v\n to positionReport{}", positionReport)
		}
		positionReports = append(positionReports, positionReportCast)
	}
	session.MessageSequenceNumber++
	return nil
}

/*

8=FIX.4.4|9=106|35=A|34=1|49=CSERVER|50=TRADE|52=20170117-08:03:04.509|56=live.theBroker.12345|57=any_string|98=0|108=30|141=Y|10=066|
8=FIX.4.4|9=109|35=5|34=1|49=CSERVER|50=TRADE|52=20170117-08:03:04.509|56=live.theBroker.12345|58=InternalError: RET_INVALID_DATA|10=033|
*/
