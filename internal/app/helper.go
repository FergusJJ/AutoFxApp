package app

import (
	"fmt"
	"pollo/internal/logs"
	"strings"
)

func calculateProfits(boughtAt, currentPrice, volume float64, side string) (netProfit float64) {
	//need comission & swap to calculate gross profit.Going to have to store more in database for position, or could just store orderID and get the execution report?
	switch side {
	case "buy":
		netProfit = (currentPrice - boughtAt) * volume
	case "sell":
		netProfit = (boughtAt - currentPrice) * volume
	}

	return netProfit
}

func SendError(err error, license string) {
	logs.SendApplicationLog(err, license)
}

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}

func roundFloat(val float64) string {
	// Start by rounding to 5 decimal places
	str := fmt.Sprintf("%.5f", val)

	// Trim unnecessary trailing zeroes
	str = strings.TrimRight(str, "0")

	// If the last character is a dot, trim that too
	if strings.HasSuffix(str, ".") {
		str = strings.TrimRight(str, ".")
	}

	return str
}

/*

{
    "positionId": 55742252,
    "positionStatus": 1,
    "tradeSide": 2,
    "symbol": {
        "symbolName": "GBPJPY",
        "digits": 3,
        "pipPosition": 2,
        "symbolId": 7,
        "description": "Great Britain Pound vs Japanese Yen"
    },
    "volume": 2000000,
    "entryPrice": 186.347,
    "openTimestamp": 1692295096822,
    "utcLastUpdateTimestamp": 1692305940442,
    "commission": -106,
    "swap": -375,
    "marginRate": 1.1729654913552443,
    "profit": 6700,
    "profitInPips": 53.1,
    "currentPrice": 185.816,
    "channel": "FIX",
    "mirroringCommission": 0,
    "usedMargin": 23459,
    "introducingBrokerCommission": 0,
    "moneyDigits": 2,
    "pnlConversionFee": 0
}

*/
