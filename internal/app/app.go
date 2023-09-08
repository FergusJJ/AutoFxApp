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

	var tmpSymbols []string = make([]string, 0)
	for _, position := range storagePositions {
		avgPx, err := strconv.ParseFloat(position.AveragePrice, 64)
		if err != nil {
			app.Program.SendColor("An unexpected error occurred", "red")
			logs.SendApplicationLog(err, app.LicenseKey)
			return
		}

		volumeInt, err := strconv.Atoi(position.Volume)
		if err != nil {
			app.Program.SendColor("An unexpected error occurred", "red")
			logs.SendApplicationLog(err, app.LicenseKey)
			return
		}

		app.FxSession.Positions[position.CopyPositionID] = fix.Position{
			PID:        position.PositionID,
			CopyPID:    position.CopyPositionID,
			Side:       position.Side,
			Symbol:     fmt.Sprint(position.SymbolID),
			SymbolName: position.Symbol,
			AvgPx:      avgPx,
			Volume:     int64(volumeInt),
			Timestamp:  position.OpenedTimestamp,
		}
		if !contains(tmpSymbols, fmt.Sprint(position.SymbolID)) {
			tmpSymbols = append(tmpSymbols, fmt.Sprint(position.SymbolID))
		}
	}
	for _, symbol := range tmpSymbols {
		app.FxSession.NewMarketDataSubscription(symbol)
	}

	app.Program.SendColor(fmt.Sprintf("retrieved %d positions from storage", len(storagePositions)), "yellow")

	go app.ApiSession.ListenForMessages()
	//Needs to add modify event, which will require additional data to be stored
	for {
		select {
		case currentMessage := <-app.ApiSession.Client.CurrentMessage:
			app.Program.SendColor("incoming message", "yellow")

			newMessage := &api.ApiMonitorMessage{}
			err := json.Unmarshal(currentMessage, newMessage)
			if err != nil {
				log.Fatalf("%+v", err)
			}

			switch newMessage.MessageType {
			case "OPEN":
				thirtySecondsMillis := 30000
				if (int(time.Now().UnixMilli()) - newMessage.OpenedTimestamp) > thirtySecondsMillis {
					break
				}
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
				newPos.SymbolName = newMessage.Symbol
				app.FxSession.Positions[newPos.CopyPID] = *newPos

				ctraderFormat := "20060102-15:04:05.000"
				ts, err := time.Parse(ctraderFormat, newPos.Timestamp)
				if err != nil {
					app.Program.SendColor("An unexpected error occurred", "red")
					logs.SendApplicationLog(err, app.LicenseKey)
					return
				}
				apiPosition := api.ApiStoredPosition{
					CopyPositionID:  newPos.CopyPID,
					PositionID:      newPos.PID,
					OpenedTimestamp: ts.String(),
					Symbol:          newMessage.Symbol,
					SymbolID:        newPos.Symbol,
					Volume:          fmt.Sprint(newPos.Volume),
					Side:            newPos.Side,
					AveragePrice:    fmt.Sprintf("%5f", newPos.AvgPx),
				}
				var apiErr *api.ApiError = &api.ApiError{}
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
					apiErr.ErrorMessage = fmt.Errorf("failed after 3 retries: %v", apiErr.ErrorMessage)
					logs.SendApplicationLog(apiErr.ErrorMessage, app.LicenseKey)
					return
				}
				symbol := fmt.Sprint(newMessage.SymbolID)
				app.Program.SendColor("creating new marketData subscription", "yellow")
				app.FxSession.NewMarketDataSubscription(symbol)
			case "CLOSE":

				_, exists := app.FxSession.Positions[newMessage.CopyPID]
				if !exists {
					continue
				}
				positionToClose := app.FxSession.Positions[newMessage.CopyPID]
				closeSymbolID := app.FxSession.Positions[newMessage.CopyPID].Symbol
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
				// log.Printf("%+v", newMessage)
				// symbol := fmt.Sprint(newMessage.SymbolID)
				delete(app.FxSession.Positions, newMessage.CopyPID)
				delete(app.UiPositionsDataMap, newMessage.CopyPID)
				app.FxSession.CheckRemoveMarketDataSubscription(closeSymbolID)

			default:
				log.Fatalln("uknown message type sent to the ", err)
			}
		default:
			//send marketDataRequests, and then update ui
			// app.Program.SendColor(fmt.Sprintf("updating %d positions", len(app.FxSession.Positions)), "yellow")
			//need to make sure that removed positions are removed from the ui map
			if len(app.FxSession.Positions) == 0 {
				time.Sleep(2 * time.Second)
				app.Program.Program.Send(PositionMessageSlice(app.UiPositionsDataMap))
				continue
			}
			var marketDataSnapshots []fix.MarketDataSnapshot
			for _, subscription := range app.FxSession.MarketDataSubscriptions {
				marketDataSnapshot, fxErr := app.FxSession.CtraderMarketDataRequest(app.FxUser, *subscription)
				if fxErr != nil {
					//need to sort this out
					if strings.Contains(fxErr.ErrorCause.Error(), "ALREADY_SUBSCRIBED") {
						continue
					}
					if fxErr.ShouldExit {
						app.Program.SendColor(fxErr.UserMessage, "red")
						logs.SendApplicationLog(fxErr.ErrorCause, app.LicenseKey)
						return
					}
					app.Program.SendColor(fxErr.UserMessage, "yellow")

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
					symbolName:     v.SymbolName,
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
				if currentPriceStr == "0" {
					continue
				}
				app.UiPositionsDataMap[newUiPosition.copyPositionId] = newUiPosition
			}
			app.Program.Program.Send(PositionMessageSlice(app.UiPositionsDataMap))
			time.Sleep(500 * time.Millisecond)
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
