package fix

import (
	"strconv"
	"strings"

	"log"
)

// this should really return an interface as well

func parseSecurityList(securityBlob string) (map[string]string, error) {

	//removing padding

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

// check whether subscription to symbol exists already, if not add MarketDataSubscription
func (session *FxSession) NewMarketDataSubscription(symbol string) {
	var md MarketDataSubscription
	session.mtx.Lock()
	defer session.mtx.Unlock()
	if session.MarketDataSubscriptions == nil {
		session.MarketDataSubscriptions = make(map[string]*MarketDataSubscription)
	}
	if _, exists := session.MarketDataSubscriptions[symbol]; exists {
		return
	}

	md = MarketDataSubscription{
		MDReqID:        GetUUID(),
		Action:         "subscribe",
		MarketDepth:    "spot",
		MDUpdateType:   "incrementalRefresh",
		NoMDEntryTypes: 2,
		MDEntryType:    []int{0, 1},
		NoRelatedSym:   1,
		Symbol:         symbol,
	}

	session.MarketDataSubscriptions[symbol] = &md
}

// will just change to unsubscribe, remove from mapping once unsub message has been sent
func (session *FxSession) RemoveMarketDataSubscription(symbol string) {
	session.MarketDataSubscriptions[symbol].Action = "unsubscribe"
}

// if there are no positions with the symbol anymore, then remove market subscription
func (session *FxSession) CheckRemoveMarketDataSubscription(symbol string) {
	for _, v := range session.Positions {
		if v.Symbol == symbol {
			return
		}
	}
	session.RemoveMarketDataSubscription(symbol)
	//
}
