package fix_test

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func Test_Parser(t *testing.T) {
	// ctrl := []byte("\u0001")
	stringToParse := "8=FIX.4.4|9=193|35=AP|34=2|49=CServer|50=TRADE|52=20230806-22:06:14.396|56=demo.ctrader.3697899|57=TRADE|55=1|710=48c8a9a9-a7f4-454d-b51e-5629f1094caa|721=47888826|727=6|728=0|730=1.07829|702=1|704=1000|705=0|10=137|8=FIX.4.4|9=194|35=AP|34=3|49=CServer|50=TRADE|52=20230806-22:06:14.396|56=demo.ctrader.3697899|57=TRADE|55=3|710=48c8a9a9-a7f4-454d-b51e-5629f1094caa|721=52942663|727=6|728=0|730=157.742|702=1|704=0|705=12000|10=176|8=FIX.4.4|9=193|35=AP|34=4|49=CServer|50=TRADE|52=20230806-22:06:14.396|56=demo.ctrader.3697899|57=TRADE|55=1|710=48c8a9a9-a7f4-454d-b51e-5629f1094caa|721=47888815|727=6|728=0|730=1.07829|702=1|704=1000|705=0|10=137|"
	dataStr := strings.ReplaceAll(stringToParse, "|", "\u0001")
	data := []byte(dataStr)
	var messages [][]byte

	//start of new stuff
	startIndex := 0
	for {
		if len(data[startIndex:]) <= 7 {
			//10=xxx| is trailer, if less than 7 then have hit end of message
			break
		}
		index := bytes.Index(data[startIndex:], []byte("\u000110="))
		if index >= 0 {
			endIndex := index + startIndex + 7 // + 8 to account for actual trailer
			if endIndex <= len(data) && string(data[endIndex]) == "\u0001" {
				messages = append(messages, data[startIndex:endIndex+1])
				startIndex = endIndex + 1
			} else {
				break
			}
		} else {
			break
		}

	}
	log.Println(messages)
	for _, msg := range messages {
		log.Println(string(msg))
	}

}
