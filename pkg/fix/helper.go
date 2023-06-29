package fix

import (
	"strconv"
	"strings"

	"log"
)

// this should really return an interface as well
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

func parsePositionReport(positionReport string) {
	positionReport = strings.Trim(positionReport, "\u0002")
	positionReport = strings.Trim(positionReport, "\u0001")
	positionReport = strings.ReplaceAll(positionReport, "\u0001", "|")

	log.Println(positionReport)

	splitAtTagSep := strings.Split(positionReport, "|")
	// tagValMap := make(map[int]string)
	//should never happen
	if len(splitAtTagSep) == 0 {
		log.Fatal("programming error at: len(splitAtTagSep) == 0")
	}
	for i := range splitAtTagSep {
		tagValPair := strings.Split(splitAtTagSep[i], "=")
		if len(tagValPair) == 0 {
			continue
		}
		//tag 728 contains amount of positions
		log.Printf("tag: %s, val: %s", tagValPair[0], tagValPair[1])

	}

}

/*

#1 order

8=FIX.4.4|9=192|35=AP|34=2|49=CServer|50=TRADE|52=20230629-11:25:27.011|56=demo.ctrader.3697899|57=TRADE|55=42|710=79cba5a2-a91a-4372-831d-e6cef7a570e9|721=47831861|727=3|728=0|730=23.16|702=1|704=1000|705=0|10=095|
8=FIX.4.4|9=193|35=AP|34=3|49=CServer|50=TRADE|52=20230629-11:25:27.011|56=demo.ctrader.3697899|57=TRADE|55=1|710=79cba5a2-a91a-4372-831d-e6cef7a570e9|721=47888826|727=3|728=0|730=1.07829|702=1|704=1000|705=0|10=168|
8=FIX.4.4|9=193|35=AP|34=4|49=CServer|50=TRADE|52=20230629-11:25:27.011|56=demo.ctrader.3697899|57=TRADE|55=1|710=79cba5a2-a91a-4372-831d-e6cef7a570e9|721=47888815|727=3|728=0|730=1.07829|702=1|704=1000|705=0|10=167|

#2 orders

8=FIX.4.4|9=192|35=AP|34=2|49=CServer|50=TRADE|52=20230629-11:26:54.861|56=demo.ctrader.3697899|57=TRADE|55=42|710=4229c759-442f-4ba5-a978-bfed33abdb9e|721=47831861|727=3|728=0|730=23.16|702=1|704=1000|705=0|10=164|
8=FIX.4.4|9=193|35=AP|34=3|49=CServer|50=TRADE|52=20230629-11:26:54.861|56=demo.ctrader.3697899|57=TRADE|55=1|710=4229c759-442f-4ba5-a978-bfed33abdb9e|721=47888826|727=3|728=0|730=1.07829|702=1|704=1000|705=0|10=237|

#0 orders

8=FIX.4.4|9=192|35=AP|34=2|49=CServer|50=TRADE|52=20230629-11:30:59.562|56=demo.ctrader.3697899|57=TRADE|55=42|710=0999387d-c0d8-4410-a663-faa041b5c4fe|721=47831861|727=3|728=0|730=23.16|702=1|704=1000|705=0|10=005|


#order resp
#52897018

#get pos result
#

*/
