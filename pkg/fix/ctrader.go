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
	err = parseFIXResponse(resp.body, Logon)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	session.MessageSequenceNumber++
	return nil
}

// func (session *FxSession) CtraderLogout(user *FxUser) {

// }
// XAGUSD=42
// GBPUSD=2
func (session *FxSession) CtraderSecurityList(user FxUser) *ErrorWithCause {
	// securityRequestID :=  //"Sxo2Xlb1jzJB" //idk whether this has to be different between users?
	securityMessage, err := user.constructSecurityList(session, "Sxo2Xlb1jzJC")
	fmt.Println(securityMessage)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}

	resp := session.sendMessage(securityMessage, user)
	if resp.err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ConnectionError,
		}
	}

	err = parseFIXResponse(resp.body, SecurityListRequest)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}

	secListPairs, err := parseSecurityList(string(resp.body))
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}
	session.SecListPairs = secListPairs
	session.MessageSequenceNumber++
	return nil

}

/* NEED EXTRA PARAMS*/

// newOrderData := idos.OrderSingleData{OrderQty: "1000", Symbol: "1", Side: "buy", OrdType: "market"}
func (session *FxSession) CtraderNewOrderSingle(user FxUser, orderData OrderData) *ErrorWithCause {
	orderData = OrderData{
		Direction: "buy",
		Volume:    12000.00,
		Symbol:    "1",
		OrderType: "market",
	}

	orderMessage, err := user.constructNewOrderSingle(session, orderData)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   ProgramError,
		}
	}
	fmt.Println(orderMessage)
	resp := session.sendMessage(orderMessage, user)
	if resp.err != nil {
		return &ErrorWithCause{
			ErrorMessage: resp.err.Error(),
			ErrorCause:   ConnectionError,
		}
	}
	err = parseFIXResponse(resp.body, NewOrderSingle)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	log.Println(string(resp.body))
	session.MessageSequenceNumber++
	return nil
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
	err = parseFIXResponse(resp.body, OrderMassStatusRequest)
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
	err = parseFIXResponse(resp.body, RequestForPositions)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	parsePositionReport(string(resp.body))
	session.MessageSequenceNumber++
	return nil
}
