package app

import (
	"encoding/json"
	"fmt"
	"log"
	"pollo/internal/logs"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"strconv"
	"strings"
	"time"

	"github.com/fasthttp/websocket"
)

const retryDelay = 5 * time.Second

// going to have to log errors here
func (app *FxApp) MainLoop() {
	app.Program.SendColor("Logging in to ctrader", "yellow")
	var gotInitialCtraderPositions bool = false
	var maxRetries = 3

	var fxErr *fix.CtraderError
	for tries := 0; tries < maxRetries; tries++ {
		fxErr = app.FxSession.CtraderLogin(app.FxUser, fix.QUOTE)
		if fxErr == nil {
			break
		}
		if fxErr.ShouldExit {
			app.Program.SendColor(fxErr.UserMessage, "red")
			logs.SendApplicationLog(fxErr.ErrorCause, app.LicenseKey)
			return
		}
		app.Program.SendColor(fxErr.UserMessage, "yellow")
		time.Sleep(retryDelay)
	}
	if fxErr != nil {
		fxErr.ErrorCause = fmt.Errorf("failed after 3 retries: %v", fxErr.ErrorCause)
		logs.SendApplicationLog(fxErr.ErrorCause, app.LicenseKey)
		return
	}

	for tries := 0; tries < maxRetries; tries++ {
		fxErr = app.FxSession.CtraderLogin(app.FxUser, fix.TRADE)
		if fxErr == nil {
			break
		}
		if fxErr.ShouldExit {
			app.Program.SendColor(fxErr.UserMessage, "red")
			logs.SendApplicationLog(fxErr.ErrorCause, app.LicenseKey)
			return
		}
		app.Program.SendColor(fxErr.UserMessage, "yellow")
		time.Sleep(retryDelay)
	}
	if fxErr != nil {
		fxErr.ErrorCause = fmt.Errorf("failed after 3 retries: %v", fxErr.ErrorCause)
		logs.SendApplicationLog(fxErr.ErrorCause, app.LicenseKey)
		return
	}

	app.Program.SendColor("Logged in to ctrader", "green")

	var apiErr *api.ApiError
	var storagePositions = []api.ApiStoredPosition{}
	for tries := 0; tries < maxRetries; tries++ {
		tmpStoragePositions, apiErr := app.ApiSession.FetchPositions()
		if apiErr == nil {
			storagePositions = tmpStoragePositions
			break
		}
		if apiErr.ShouldExit {
			app.Program.SendColor(apiErr.UserMessage, "red")
			logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
			return
		}
		if apiErr.ErrorType == api.ApiConnectionError {
			app.Program.SendColor(apiErr.UserMessage, "yellow")
			time.Sleep(retryDelay)
			continue
		}
		if apiErr.ErrorType == api.ApiAuthorizationError {
			app.Program.SendColor(apiErr.UserMessage, "yellow")
			apiErr = app.ApiSession.RefreshApiAuth()
			if apiErr != nil {
				if apiErr.ShouldExit {
					app.Program.SendColor(apiErr.UserMessage, "red")
					logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
					return
				}
				//in case that gets connection error on reauth, just retry whole thing
				app.Program.SendColor(apiErr.UserMessage, "yellow")
			}
			time.Sleep(retryDelay)
			continue
		}
	}
	if apiErr != nil {
		apiErr.ErrorMessage = fmt.Errorf("failed after 3 retries: %v", fxErr.ErrorCause)
		logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
		return
	}

	app.UiPositionsDataMap = make(map[string]uiPositionData, 0)
	app.FxSession.Positions = make(map[string]fix.Position, 0)
	app.Program.SendColor(fmt.Sprintf("retrieved %d positions from storage", len(storagePositions)), "yellow")
	var tmpSymbols []string = make([]string, 0)
	var positionInfoMap = make(map[string]fix.PositionReport, 0)
	var tmpPIDCopyPIDMapping = make(map[string]string, 0)

	for !gotInitialCtraderPositions {
		positionInfo, ctErr := app.FxSession.CtraderRequestForPositions(app.FxUser)
		if ctErr != nil {
			if ctErr.ShouldExit {
				app.Program.SendColor(ctErr.UserMessage, "red")
				logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
				return
			}
			if ctErr.ErrorType == fix.CtraderConnectionError {
				app.Program.Program.Send(ctErr.UserMessage)
				time.Sleep(retryDelay)
				continue
			}
		}
		totalExpectedReports, err := strconv.Atoi(positionInfo[0].TotalNumPosReports)
		if err != nil {
			app.Program.SendColor("An unexpected error occurred", "red")
			logs.SendApplicationLog(err, app.LicenseKey)
			return
		}

		if totalExpectedReports == 0 {
			app.Program.SendColor("no positions open to be fetched", "green")
			gotInitialCtraderPositions = true
			continue
		}

		for _, positionReport := range positionInfo {
			if _, exists := positionInfoMap[positionReport.PosMaintRptID]; !exists {
				positionInfoMap[positionReport.PosMaintRptID] = positionReport
				if !contains(tmpSymbols, positionReport.Symbol) {
					tmpSymbols = append(tmpSymbols, positionReport.Symbol)
				}
			}
		}

		//takes time for a large number of  positions to come through so keep re-reading
		//since when fix sends a message comprised of multiple messages, they are not wrapped in one message with a checksum
		//they are just send as individual messages
		for totalExpectedReports > len(positionInfoMap) {
			rereadMessages, err := app.FxSession.TradeClient.ReRead()
			if err != nil {
				app.Program.SendColor("An unexpected error occurred", "red")
				logs.SendApplicationLog(err, app.LicenseKey)
				return
			}
			app.Program.SendColor(fmt.Sprintf("got %d/%d", len(positionInfoMap), totalExpectedReports), "green")
			time.Sleep(100 * time.Millisecond)
			for _, message := range rereadMessages {

				fxRes, ctErr := fix.ParseFixResponse(message, fix.RequestForPositions)
				if ctErr != nil {
					if ctErr.ShouldExit {
						app.Program.SendColor(ctErr.UserMessage, "red")
						logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
						return
					}
				}
				positionReport, ok := fxRes.(fix.PositionReport)
				if !ok {

					rejectMsg, ok := fxRes.(fix.SessionRejectMessage)
					if !ok {
						ctErr := &fix.CtraderError{
							UserMessage: "An unexpected error occurred",
							ErrorType:   fix.CtraderLogicError,
							ErrorCause:  fmt.Errorf("unable to convert interface to SessionRejectMessage"),
							ShouldExit:  true,
						}
						app.Program.SendColor(ctErr.UserMessage, "red")
						logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
						return
					}
					ctErr := fix.ErrorFromSessionReject(rejectMsg)
					app.Program.SendColor(ctErr.UserMessage, "yellow")
					logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
					return
				}
				if _, exists := positionInfoMap[positionReport.PosMaintRptID]; !exists {
					positionInfoMap[positionReport.PosMaintRptID] = positionReport
					if !contains(tmpSymbols, positionReport.Symbol) {
						tmpSymbols = append(tmpSymbols, positionReport.Symbol)
					}
				}
			}
		}

		app.Program.SendColor("got all positions", "green")
		for _, symbol := range tmpSymbols {
			app.FxSession.NewMarketDataSubscription(symbol)
		}
		for _, v := range storagePositions {
			tmpPIDCopyPIDMapping[v.PositionID] = v.CopyPositionID
		}
		for _, position := range positionInfoMap {

			side := "buy"
			vol := position.LongQty
			if position.LongQty == "0" {
				vol = position.ShortQty
				side = "sell"
			}
			volInt, err := strconv.ParseInt(vol, 10, 64)
			if err != nil {
				ctErr := &fix.CtraderError{
					UserMessage: "An unexpected error occurred",
					ErrorType:   fix.CtraderLogicError,
					ErrorCause:  err,
					ShouldExit:  true,
				}
				app.Program.SendColor(ctErr.UserMessage, "red")
				logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
				return
			}
			avgPx, err := strconv.ParseFloat(position.SettlPrice, 64)
			if err != nil {
				ctErr := &fix.CtraderError{
					UserMessage: "An unexpected error occurred",
					ErrorType:   fix.CtraderLogicError,
					ErrorCause:  err,
					ShouldExit:  true,
				}
				app.Program.SendColor(ctErr.UserMessage, "red")
				logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
				return
			}
			if _, exists := tmpPIDCopyPIDMapping[position.PosMaintRptID]; !exists {
				continue
			}
			app.FxSession.Positions[tmpPIDCopyPIDMapping[position.PosMaintRptID]] = fix.Position{
				PID:       position.PosMaintRptID,
				CopyPID:   tmpPIDCopyPIDMapping[position.PosMaintRptID],
				Side:      side,
				Symbol:    position.Symbol,
				AvgPx:     avgPx,
				Volume:    volInt,
				Timestamp: "",
			}
		}
		gotInitialCtraderPositions = true

	}

	go app.ApiSession.ListenForMessages()
	//Needs to add modify event, which will require additional data to be stored
	for {
		select {
		case currentMessage := <-app.ApiSession.Client.CurrentMessage:
			newMessage := &api.ApiMonitorMessage{}
			err := json.Unmarshal(currentMessage, newMessage)
			if err != nil {
				log.Fatalf("%+v", err)
			}
			switch newMessage.MessageType {
			case "OPEN":
				app.Program.SendColor(fmt.Sprintf("%+v", newMessage), "green")
				//any non-fatal errors should be handled within the function, all errors at this point should quit
				newPos, ctErr := app.OpenPosition(newMessage)
				if ctErr != nil {
					if ctErr.ShouldExit {
						app.Program.SendColor(ctErr.UserMessage, "red")
						logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
						return
					}
					logs.SendApplicationLog(fmt.Errorf("app.OpenPosition returned not-fatal error: %w", ctErr.ErrorCause), app.LicenseKey)
					return
				}
				app.FxSession.Positions[newPos.CopyPID] = *newPos
				apiPosition := api.ApiStoredPosition{
					CopyPositionID: newPos.CopyPID,
					PositionID:     newPos.PID,
				}
				var apiErr *api.ApiError
				for tries := 0; tries < maxRetries; tries++ {
					apiErr = app.ApiSession.StorePosition(apiPosition)
					if apiErr == nil {
						break
					}
					if apiErr.ShouldExit {
						app.Program.SendColor(apiErr.UserMessage, "red")
						logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
						return
					}
					if apiErr.ErrorType == api.ApiConnectionError {
						app.Program.SendColor(apiErr.UserMessage, "yellow")
						time.Sleep(retryDelay)
						continue
					}
					if apiErr.ErrorType == api.ApiAuthorizationError {
						app.Program.SendColor(apiErr.UserMessage, "yellow")
						apiErr = app.ApiSession.RefreshApiAuth()
						if apiErr != nil {
							if apiErr.ShouldExit {
								app.Program.SendColor(apiErr.UserMessage, "red")
								logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
								return
							}
							//in case that gets connection error on reauth, just retry whole thing
							app.Program.SendColor(apiErr.UserMessage, "yellow")
						}
						time.Sleep(retryDelay)
						continue
					}
				}
				if apiErr != nil {
					apiErr.ErrorMessage = fmt.Errorf("failed after 3 retries: %v", fxErr.ErrorCause)
					logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
					return
				}
				symbol := fmt.Sprint(newMessage.SymbolID)
				app.Program.SendColor("creating new marketData subscription", "yellow")
				app.FxSession.NewMarketDataSubscription(symbol)
			case "CLOSE":
				app.Program.SendColor(fmt.Sprintf("closing position %s", newMessage.CopyPID), "green")
				// newMessage.CopyPID
				_, exists := app.FxSession.Positions[newMessage.CopyPID]
				if !exists {
					app.Program.SendColor("user doesn't have position", "green")
					continue
				}
				positionToClose := app.FxSession.Positions[newMessage.CopyPID]
				ctErr := app.ClosePosition(positionToClose)
				if ctErr != nil {
					if ctErr.ShouldExit {
						app.Program.SendColor(ctErr.UserMessage, "red")
						logs.SendApplicationLog(ctErr.ErrorCause, app.LicenseKey)
						return
					}
					logs.SendApplicationLog(fmt.Errorf("app.OpenPosition returned not-fatal error: %w", ctErr.ErrorCause), app.LicenseKey)
					return
				}
				var apiErr *api.ApiError
				for tries := 0; tries < maxRetries; tries++ {

					apiErr = app.ApiSession.RemovePosition(positionToClose.PID)
					if apiErr == nil {
						break
					}
					if apiErr.ShouldExit {
						app.Program.SendColor(apiErr.UserMessage, "red")
						logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
						return
					}
					if apiErr.ErrorType == api.ApiConnectionError {
						app.Program.SendColor(apiErr.UserMessage, "yellow")
						time.Sleep(retryDelay)
						continue
					}
					if apiErr.ErrorType == api.ApiAuthorizationError {
						app.Program.SendColor(apiErr.UserMessage, "yellow")
						apiErr = app.ApiSession.RefreshApiAuth()
						if apiErr != nil {
							if apiErr.ShouldExit {
								app.Program.SendColor(apiErr.UserMessage, "red")
								logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
								return
							}
							//in case that gets connection error on reauth, just retry whole thing
							app.Program.SendColor(apiErr.UserMessage, "yellow")
						}
						time.Sleep(retryDelay)
						continue
					}
				}
				if apiErr != nil {
					apiErr.ErrorMessage = fmt.Errorf("failed after 3 retries: %v", fxErr.ErrorCause)
					logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
					return
				}

				symbol := fmt.Sprint(newMessage.SymbolID)
				app.FxSession.CheckRemoveMarketDataSubscription(symbol)

			default:
				log.Fatalln("uknown message type sent to the ", err)
			}
		default:
			//send marketDataRequests, and then update ui
			app.Program.SendColor(fmt.Sprintf("updating %d positions", len(app.FxSession.Positions)), "yellow")

			var marketDataSnapshots []fix.MarketDataSnapshot
			for _, subscription := range app.FxSession.MarketDataSubscriptions {
				marketDataSnapshot, fxErr := app.FxSession.CtraderMarketDataRequest(app.FxUser, *subscription)
				if fxErr != nil {
					if strings.Contains(fxErr.ErrorCause.Error(), "ALREADY_SUBSCRIBED") {
						continue
					}
					if fxErr.ShouldExit {
						app.Program.SendColor(apiErr.UserMessage, "red")
						logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
						return
					}
					app.Program.SendColor(apiErr.UserMessage, "yellow")

				}
				marketDataSnapshots = append(marketDataSnapshots, marketDataSnapshot...)
			}

			var symbolPricePairs = map[string]float64{}
			for _, v := range marketDataSnapshots {
				if v.MDEntryPx == "" {
					continue
				}
				pxFloat, err := strconv.ParseFloat(v.MDEntryPx, 64)
				if err != nil {
					app.Program.SendColor(fmt.Sprintf("error parsing price for symbol, %s", err.Error()), "red")
				}
				symbolPricePairs[v.Symbol] = pxFloat
			}

			for _, v := range app.FxSession.Positions {
				entry := v.AvgPx
				// app.Program.SendColor(fmt.Sprintf("volume: %s current: %f, side: %s", v.Volume, symbolPricePairs[v.Symbol], v.Side), "red")

				grossProfit := calculateProfits(entry, symbolPricePairs[v.Symbol], float64(v.Volume), v.Side)

				volumeStr := strconv.Itoa(int(v.Volume))
				grossProfitStr := roundFloat(grossProfit)
				if v.Timestamp == "" {
					v.Timestamp = "Unavailable"
				}
				entryStr := roundFloat(entry)
				currentPriceStr := roundFloat(symbolPricePairs[v.Symbol])

				newUiPosition := uiPositionData{
					entryPrice:     entryStr,
					currentPrice:   currentPriceStr,
					copyPositionId: v.CopyPID,
					positionId:     v.PID,
					volume:         volumeStr,
					grossProfit:    grossProfitStr,
					side:           v.Side,
					timestamp:      v.Timestamp,
					symbol:         v.Symbol,
					isProfit:       grossProfit > 0,
				}
				app.UiPositionsDataMap[newUiPosition.copyPositionId] = newUiPosition
			}
			app.Program.Program.Send(PositionMessageSlice(app.UiPositionsDataMap))

			time.Sleep(2 * time.Second)
			continue
		}
	}
}

func (app *FxApp) CloseExistingConnections() {
	if app.FxSession.TradeClient != nil {
		//this will never panic, so logging error is fine
		err := app.FxSession.TradeClient.Close()
		if err != nil {
			app.Program.SendColor(fmt.Sprintf("error closing connection to FIX: %s", err.Error()), "red")
		}
		app.FxSession.TradeClient = nil
	}

	if app.FxSession.PriceClient != nil {
		err := app.FxSession.PriceClient.Close()
		if err != nil {
			app.Program.SendColor(fmt.Sprintf("error closing connection to FIX: %s", err.Error()), "red")
		}
		app.FxSession.PriceClient = nil

	}

	if app.ApiSession.Client.Connection != nil {
		err := app.ApiSession.Client.Connection.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		if err != nil {
			app.Program.SendColor(fmt.Sprintf("error writing close to API: %s", err.Error()), "red")
		}
		//this will never panic, so logging error is fine
		err = app.ApiSession.Client.Connection.Close()
		if err != nil {
			app.Program.SendColor(fmt.Sprintf("error closing connection to API: %s", err.Error()), "red")
		}
		app.ApiSession.Client.Connection = nil
	}
}
