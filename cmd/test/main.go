package main

import (
	"pollo/internal/app/display"
)

func main() {
	var testData = [][]interface{}{
		{"PID332510937", "USDCHF", "BUY", "12000.00", "-4.76", "02 May 2023 17:11:04.983"},
		{"PID332522747", "EURNZD", "SELL", "900.00", "61.7", "04 May 2023 18:07:17.815"},
	}
	display.DrawTable(testData)
}
