package fix

import (
	"encoding/json"
	"strconv"
	"strings"

	"log"
)

// this should really return an interface as well
func ParseFIXResponse(body []byte, messageType CtraderSessionMessageType) (interface{}, error) {

	bodyString := strings.ReplaceAll(string(body), "\u0001", "|")

	messageBodyAndTag := stripHeaderAndTrailer(bodyString)
	if messageBodyAndTag.Tag == "5" {
		err := FixAPIError{ErrorMessage: ErrorInvalidLogon, ShouldRetry: false, AdditionalContext: messageBodyAndTag.MessageBody["58"]}
		return nil, err
	}

	switch messageType {
	case Logon, Logout:

	case Heartbeat:

	case TestRequest:

	case Resend:

	case SequenceReset:

	case RequestForPositions:
		switch messageBodyAndTag.Tag {
		case "AP":
			log.Println(messageBodyAndTag.MessageBody)
			log.Println("requestForPositions")
		default:
			log.Fatalf("case %s not handled for RequestForPositions in ParseFIXResponse", messageBodyAndTag.Tag)
		}

	case NewOrderSingle:
		switch messageBodyAndTag.Tag {
		case "8":
			var executionReportMapping = map[string]string{}
			var executionReport ExecutionReport
			//execution report, order has gone through
			for tag, val := range messageBodyAndTag.MessageBody {
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
		case "9":
			var orderCancelRejectMapping = map[string]string{}
			var orderCancelReject OrderCancelReject
			for tag, val := range messageBodyAndTag.MessageBody {
				orderCancelRejectMapping[orderCancelRejectTagMapping[tag]] = val
			}
			jsonData, err := json.Marshal(orderCancelRejectMapping)
			if err != nil {
				log.Fatal(err)
			}
			err = json.Unmarshal(jsonData, &orderCancelReject)
			if err != nil {
				log.Fatal(err)
			}
			return orderCancelReject, nil

		case "j":
			//BusinessMessageReject, Indicates issues with the FIX message.
			var businessMessageRejectMapping = map[string]string{}
			var businessMessageReject BusinessMessageReject
			for tag, val := range messageBodyAndTag.MessageBody {
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
			return businessMessageReject, nil
		default:
			log.Fatalf("case %s not handled for NewOrderSingle in ParseFIXResponse", messageBodyAndTag.Tag)
		}
	case SecurityListRequest:

	case OrderStatusRequest:
		switch messageBodyAndTag.Tag {
		case "8":
		default:

		}
	//for when no response data is needed
	default:
		log.Println("got good response:")
		log.Println(messageBodyAndTag.MessageBody)

		return nil, nil
	}
	return nil, nil
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

// want to return {messageType; body}
func stripHeaderAndTrailer(message string) *MessageBodyAndTag {
	message = strings.Trim(message, "\u0002")
	message = strings.Trim(message, "\u0001")
	// message = strings.Trim(message, "\x00") //bunch of these
	var strippedMessage = &MessageBodyAndTag{
		Tag:         "",
		MessageBody: map[string]string{},
	}
	tagSlice := strings.Split(message, "|")
	if len(tagSlice) == 1 {
		return nil
	}
	//always has trailing "|" so last index will be ""

	for _, tagAndVal := range tagSlice {
		// if tagAndVal == ""{
		// 	continue
		// }
		tagVal := strings.Split(tagAndVal, "=")
		if len(tagVal) == 1 {
			continue
		}
		switch tagVal[0] {
		case "8", "9", "49", "56", "57", "50", "34", "52", "10":
			continue
		default:
			if tagVal[0] == "35" {
				strippedMessage.Tag = tagVal[1]
				continue
			}
			strippedMessage.MessageBody[tagVal[0]] = tagVal[1]
		}
	}
	return strippedMessage

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
