package fix

import "fmt"

type CtraderErrorType int

const (
	CtraderLogicError          CtraderErrorType = 0
	CtraderConnectionError     CtraderErrorType = 1
	CtraderUserDataError       CtraderErrorType = 2
	CtraderRetryError          CtraderErrorType = 3
	CtraderBusinessRejectError CtraderErrorType = 4
)

type CtraderError struct {
	UserMessage string
	ErrorType   CtraderErrorType
	ErrorCause  error
	ShouldExit  bool
}

// anything that will quit the program and send webhook (logic errors/bugs) should contain specific data
// otherwise report minimum specific info
func ErrorFromSessionReject(reject SessionRejectMessage) *CtraderError {
	//want to decide whether reason is a bug, or something to do with user
	var err *CtraderError = &CtraderError{}

	switch reject.SessionRejectReason {
	case "7":
		err.UserMessage = "Possible connection error, retrying"
		err.ErrorCause = fmt.Errorf("decryption problem")
		err.ErrorType = CtraderConnectionError
		err.ShouldExit = false

	case "8":
		err.UserMessage = "An error occurred whilst processing a request"
		err.ErrorCause = fmt.Errorf("signature error")
		err.ErrorType = CtraderConnectionError
		err.ShouldExit = true

	case "9":
		err.UserMessage = "An error occurred, check the compID fields in data.json"
		err.ErrorCause = fmt.Errorf("compID error")
		err.ErrorType = CtraderUserDataError
		err.ShouldExit = true

	case "10":
		err.UserMessage = "Possible connection error, retrying"
		err.ErrorCause = fmt.Errorf("sending time accuracy error")
		err.ErrorType = CtraderConnectionError
		err.ShouldExit = false

	default:
		err.UserMessage = "An unexpected error occurred, please try again later"
		err.ErrorCause = fmt.Errorf("reject reason: %s reject text: %s", reject.SessionRejectReason, reject.Text)
		err.ErrorType = CtraderLogicError
		err.ShouldExit = true
	}
	return err

}

var (
	ErrorInvalidLogon         = "failed to login, invalid userdata"
	ErrorInvalidNOSRequest    = "failed to send NewOrderSingle request"
	ErrorFailedToSend         = "failed to send fix message"
	ErrorInvalidLogout        = "failed to logout of fix API"
	ErrorInvalidHeartbeat     = "failed to send FIX Heartbeat message"
	ErrorInvalidTestRequest   = "failed to send TestRequest message"
	ErrorInvalidResend        = "failed to send Resend message"
	ErrorInvalidSequenceReset = "failed to send SequenceReset message"
	ErrorInvalidSecurityList  = "failed to send SecurityList message"
	ErrorInvalidOrderStatus   = "failed to send OrderStatus message"
)
