package app

import (
	"errors"
	"fmt"
	"pollo/internal/logs"
	"pollo/pkg/fix"
	"strings"
	"time"
)

func calculateProfits(boughtAt, currentPrice, volume float64, side string) (netProfit float64) {
	switch side {
	case "buy":
		netProfit = (currentPrice - boughtAt) * volume
	case "sell":
		netProfit = (boughtAt - currentPrice) * volume
	}

	return netProfit
}

// errors that are (hopefully) the users fault should be logged to the users screen, the program should then wait for input then exit
// errors that are bugs should output a generic message, send the error to webhook, and then wait for the user to input to exit
func (app *FxApp) HandleError(cErr *fix.ErrorWithCause, errorSource string) bool {
	time.Sleep(10 * time.Second)
	switch cErr.ErrorCause {
	case fix.MarketError: //carry on with execution
		app.Program.SendColor(cErr.ErrorMessage, "red")
		logs.SendApplicationLog(errors.New(cErr.ErrorMessage), errorSource, app.LicenseKey)
		time.Sleep(retryDelay)
		return false
	case fix.ConnectionError: //carry on with execution
		app.Program.SendColor("error sending message to API, retrying", "red")
		time.Sleep(retryDelay)
		return false
	case fix.CtraderConnectionError:
		app.Program.SendColor("error sending message to FIX, retrying", "red")
		time.Sleep(retryDelay)
		return false
	case fix.UserDataError: //exit
		app.Program.SendColor(fmt.Sprintf("unexpected error occurred, exiting: %s", cErr.ErrorMessage), "red")
		return true
	case fix.ProgramError: //exit
		app.Program.SendColor("unexpected error occurred, exiting", "red")
		logs.SendApplicationLog(errors.New(cErr.ErrorMessage), errorSource, app.LicenseKey)
		return true
	default: //exit
		app.Program.SendColor("unexpected error occurred, exiting", "red")
		logs.SendApplicationLog(errors.New(cErr.ErrorMessage), errorSource, app.LicenseKey)
		return true
	}
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
