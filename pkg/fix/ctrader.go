package fix

import (
	"fmt"
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
	_, err = ParseFIXResponse(resp.body, Logon)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}

	session.MessageSequenceNumber++
	return nil
}

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

	_, err = ParseFIXResponse(resp.body, SecurityListRequest)
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
	resp := session.sendMessage(orderMessage, user)
	if resp.err != nil {
		return &ErrorWithCause{
			ErrorMessage: resp.err.Error(),
			ErrorCause:   ConnectionError,
		}
	}
	_, err = ParseFIXResponse(resp.body, NewOrderSingle)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
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
	_, err = ParseFIXResponse(resp.body, OrderMassStatusRequest)
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
	_, err = ParseFIXResponse(resp.body, RequestForPositions)
	if err != nil {
		return &ErrorWithCause{
			ErrorMessage: err.Error(),
			ErrorCause:   UserDataError,
		}
	}
	session.MessageSequenceNumber++
	return nil
}

/*

8=FIX.4.4|9=106|35=A|34=1|49=CSERVER|50=TRADE|52=20170117-08:03:04.509|56=live.theBroker.12345|57=any_string|98=0|108=30|141=Y|10=066|
8=FIX.4.4|9=109|35=5|34=1|49=CSERVER|50=TRADE|52=20170117-08:03:04.509|56=live.theBroker.12345|58=InternalError: RET_INVALID_DATA|10=033|
*/
