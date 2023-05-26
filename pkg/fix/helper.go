package fix

import (
	"strconv"
	"strings"

	"log"
)

func parseFIXResponse(body []byte, messageType CtraderSessionMessageType) error {
	badRequest, reason := checkForBadRequest(string(body))
	if !badRequest {
		log.Println("req is good")
		return nil
	} else {
		log.Println("req bad")

	}
	switch messageType {
	case Logon:
		err := FixAPIError{ErrorMessage: ErrorInvalidLogon, ShouldRetry: false, AdditionalContext: reason}
		// err = FixError{FixAPIError: {ErrorMessage: ErrorInvalidLogon, ShouldRetry: false, AdditionalContext: reason}}
		return err
	case Logout:
		err := FixAPIError{ErrorMessage: ErrorInvalidLogout, ShouldRetry: false, AdditionalContext: reason}
		return err
	case Heartbeat:
		err := FixAPIError{ErrorMessage: ErrorInvalidHeartbeat, ShouldRetry: false, AdditionalContext: reason}
		return err
	case TestRequest:
		err := FixAPIError{ErrorMessage: ErrorInvalidTestRequest, ShouldRetry: false, AdditionalContext: reason}
		return err
	case Resend:
		err := FixAPIError{ErrorMessage: ErrorInvalidResend, ShouldRetry: false, AdditionalContext: reason}
		return err
	case SequenceReset:
		err := FixAPIError{ErrorMessage: ErrorInvalidSequenceReset, ShouldRetry: false, AdditionalContext: reason}
		return err
	case NewOrderSingle:
		err := FixAPIError{ErrorMessage: ErrorInvalidNOSRequest, ShouldRetry: false, AdditionalContext: reason}
		return err
	case SecurityListRequest:
		err := FixAPIError{ErrorMessage: ErrorInvalidLogout, ShouldRetry: false, AdditionalContext: reason}
		return err
	case OrderStatusRequest:
		err := FixAPIError{ErrorMessage: ErrorInvalidOrderStatus, ShouldRetry: false, AdditionalContext: reason}
		return err

	default:
		log.Fatal("programming error: sessionMessageCode fallthrough")
		return nil
	}
}

func checkForBadRequest(responseBody string) (badRequest bool, reason string) {
	splitAtTagSep := strings.Split(responseBody, "|")
	tagValMap := make(map[int]string)
	//should never happen
	if len(splitAtTagSep) == 0 {
		log.Fatal("programming error at: len(splitAtTagSep) == 0")
	}
	for i := range splitAtTagSep {
		tagValPair := strings.Split(splitAtTagSep[i], "=")
		if len(tagValPair) == 0 {
			continue
		}
		tag, err := strconv.Atoi(tagValPair[0])
		if err != nil {
			//should never happen
			log.Fatal("programming error at: tag, err := strconv.Atoi(tagValPair[0])")
		}
		tagValMap[tag] = tagValPair[1]
	}
	if tagValMap[35] == "3" {
		badRequest = true
		reason = tagValMap[58]
	}
	return badRequest, reason
}

func parseSecurityList(securityBlob string) (map[string]string, error) {

	//removing padding
	securityBlob = strings.Trim(securityBlob, "\u0002")
	securityBlob = strings.Trim(securityBlob, "\u0001")

	securityBlob = strings.ReplaceAll(securityBlob, "\u0001", "|")
	splitAtTagSep := strings.Split(securityBlob, "|")
	tagValMap := make(map[int]string)
	//should never happen
	if len(splitAtTagSep) == 0 {
		log.Fatal("programming error at: len(splitAtTagSep) == 0")
	}

	secList := make(map[string]string)
	currTag := ""
	for i := range splitAtTagSep {
		tagValPair := strings.Split(splitAtTagSep[i], "=")
		if len(tagValPair) == 0 {
			continue
		}
		tag, err := strconv.Atoi(tagValPair[0])

		if err != nil {
			//should never happen
			log.Fatal("programming error at: tag, err := strconv.Atoi(tagValPair[0])")
		}
		if tag == 55 {
			currTag = tagValPair[1]
		}
		if tag == 1007 {
			if currTag == "" {
				log.Fatalf("i don't think this should've happened")
			}
			secList[tagValPair[1]] = currTag
			currTag = ""
		}
		tagValMap[tag] = tagValPair[1]
	}

	return secList, nil
}
