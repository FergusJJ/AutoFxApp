package fix

type ErrorType int

const (
	ConnectionError ErrorType = iota
	UserDataError
	ProgramError
)

type ErrorWithCause struct {
	ErrorMessage string
	ErrorCause   ErrorType
}

type FixAPIError struct {
	ErrorMessage      string
	AdditionalContext string // may change to struct or something idk
	ShouldRetry       bool
}

func (e FixAPIError) Error() string {
	return e.ErrorMessage
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
