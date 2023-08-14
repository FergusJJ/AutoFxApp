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
	var messageFails int = 0
	var tradeLoginFinished bool = false
	var quoteLoginFinished bool = false
	var fetchedStoragePositions bool = false
	var gotInitialCtraderPositions bool = false

	for !quoteLoginFinished {
		fxErr := app.FxSession.CtraderLogin(app.FxUser, fix.QUOTE)
		if fxErr != nil {
			errorSource := "MainLoop, app.FxSession.CtraderLogin() (QUOTE)"
			if messageFails > 3 {
				app.Program.SendColor("unexpected error occurred, exiting", "red")
				logs.SendApplicationLog(fmt.Errorf(fxErr.ErrorMessage), errorSource, app.LicenseKey)
				return
			}
			fxErr = app.HandleError(fxErr, errorSource)
			if fxErr != nil {
				return
			}
			messageFails++
			continue
		}
		quoteLoginFinished = true
	}
	messageFails = 0
	for !tradeLoginFinished {
		fxErr := app.FxSession.CtraderLogin(app.FxUser, fix.TRADE)
		if fxErr != nil {
			errorSource := "MainLoop, app.FxSession.CtraderLogin() (TRADE)"
			if messageFails > 3 {
				app.Program.SendColor("unexpected error occurred, exiting", "red")
				logs.SendApplicationLog(fmt.Errorf(fxErr.ErrorMessage), errorSource, app.LicenseKey)
				return
			}
			fxErr = app.HandleError(fxErr, errorSource)
			if fxErr != nil {
				return
			}
			messageFails++
			continue
		}
		tradeLoginFinished = true
	}
	app.Program.SendColor("Logged in to ctrader", "green")
	app.FxSession.LoggedIn = true

	//might want to re-implement security list, as will need the symbol strings
	// app.FxSession.GotSecurityList = true

	var storagePositions = []api.ApiStoredPositionsResponse{}
	for !fetchedStoragePositions {
		tmpStoragePositions, err := app.ApiSession.FetchPositions()
		if err != nil {
			errorSource := "MainLoop, app.ApiSession.FetchPositions()"
			if err.Error() == "unauthorized request" {
				app.Program.SendColor(fmt.Sprintf("%s: reauthorizing session", err.Error()), "yellow")
				refreshError := app.ApiSession.RefreshApiAuth()
				app.Program.SendColor(fmt.Sprintf("%s: error reauthorizing session", refreshError.Error()), "red")
				return
			} else if err.Error() == "internal server error" {
				app.Program.SendColor(fmt.Sprintf("%s: please try again later", err.Error()), "red")
				return
			}
			fxErr := &fix.ErrorWithCause{
				ErrorMessage: err.Error(),
				ErrorCause:   fix.ProgramError,
			}
			fxErr = app.HandleError(fxErr, errorSource)
			if fxErr != nil {
				return
			}
		}
		fetchedStoragePositions = true
		storagePositions = tmpStoragePositions
	}
	app.UiPositionsDataMap = make(map[string]uiPositionData, 0)
	app.FxSession.Positions = make(map[string]fix.Position, 0)
	app.Program.SendColor(fmt.Sprintf("retrieved %d positions from storage", len(storagePositions)), "yellow")
	var tmpSymbols []string = make([]string, 0)
	var positionInfoMap = make(map[string]fix.PositionReport, 0)
	var tmpPIDCopyPIDMapping = make(map[string]string, 0)

	for !gotInitialCtraderPositions {
		positionInfo, fxErr := app.FxSession.CtraderRequestForPositions(app.FxUser)
		if fxErr != nil {
			errSource := "MainLoop, CtraderRequestForPositions()"
			fxErr = app.HandleError(fxErr, errSource)
			if fxErr != nil {
				return
			}
		}
		if len(positionInfo) == 0 {
			app.Program.SendColor("no positions open to be fetched", "green")
			gotInitialCtraderPositions = true
		} else {
			totalExpectedReports, err := strconv.Atoi(positionInfo[0].TotalNumPosReports)
			if err != nil {
				errSource := "MainLoop, TotalNumPosReports type casting"
				fxErr := &fix.ErrorWithCause{
					ErrorMessage: err.Error(),
					ErrorCause:   fix.ProgramError,
				}
				app.HandleError(fxErr, errSource) //only one case here, so return every time
				return
			}
			for _, positionReport := range positionInfo {
				if _, exists := positionInfoMap[positionReport.PosMaintRptID]; !exists {
					positionInfoMap[positionReport.PosMaintRptID] = positionReport
					if !contains(tmpSymbols, positionReport.Symbol) {
						tmpSymbols = append(tmpSymbols, positionReport.Symbol)
					}
				}
			}
			if totalExpectedReports > len(positionInfoMap) {
				app.Program.SendColor(fmt.Sprintf("got %d/%d positions", len(positionInfoMap), totalExpectedReports), "yellow")
				time.Sleep(1 * time.Second)
				continue
			}
			app.Program.SendColor("got all positions positions", "green")
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
					errSource := "MainLoop, ParseInt(vol)"
					fxErr := &fix.ErrorWithCause{
						ErrorMessage: err.Error(),
						ErrorCause:   fix.ProgramError,
					}
					app.HandleError(fxErr, errSource)
					return
				}
				avgPx, err := strconv.ParseFloat(position.SettlPrice, 64)
				if err != nil {
					errSource := "MainLoop, ParseFloat(position.Settl)"
					fxErr := &fix.ErrorWithCause{
						ErrorMessage: err.Error(),
						ErrorCause:   fix.ProgramError,
					}
					app.HandleError(fxErr, errSource)
					return
				}
				if _, exists := tmpPIDCopyPIDMapping[position.PosMaintRptID]; !exists {
					app.Program.SendColor(fmt.Sprintf("position %s is not copied, skipping", position.PosMaintRptID), "yellow")
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
	}

	// app.Program.SendColor(fmt.Sprintf("Got %d from ctrader", len(positionInfo)), "green")
	// }

	// tmpPositions, ok := positionsAsMapSlice.([]fix.Position)
	// if !ok {
	// 	return &fix.ErrorWithCause{
	// 		ErrorMessage: "failed to convert to positions",
	// 		ErrorCause:   fix.ProgramError,
	// 	}
	// }
	// for _, v := range tmpPositions {
	// 	app.FxSession.Positions[v.CopyPID] = v
	// }
	// app.FxSession.Positions = tmpPositions
	//do this from database, ignore positions that aren't stored there as they are not being copied

	//need to start function that will monitor here:
	go app.ApiSession.ListenForMessages()
	//need to start function that will display open positions here:

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
				success, newPos := app.OpenPosition(newMessage)
				if !success {
					continue
				}
				app.FxSession.Positions[newPos.CopyPID] = *newPos
				//need to send PID:CopyPID pair to DB

				//
				symbol := fmt.Sprint(newMessage.SymbolID)
				app.Program.SendColor("creating new marketData subscription", "yellow")
				app.FxSession.NewMarketDataSubscription(symbol)
			case "CLOSE":
				app.Program.SendColor(fmt.Sprintf("closing position %s", newMessage.CopyPID), "green")

				// pid, err := app.ClosePosition(newMessage)
				// if err != nil {
				// 	app.Program.SendColor(fmt.Sprint("close position error: ", err)))
				// 	continue
				// }

				//if successful, check remove subscription
				//but remove position from Positions first

				symbol := fmt.Sprint(newMessage.SymbolID)
				app.FxSession.CheckRemoveMarketDataSubscription(symbol)

			default:
				log.Fatalln("uknown message type sent to the ", err)
			}
		default:
			//send marketDataRequests, and then update ui
			app.Program.SendColor("updating position data", "yellow")
			app.Program.SendColor(fmt.Sprintf("amount of subscriptions: %d", len(app.FxSession.MarketDataSubscriptions)), "yellow")

			var marketDataSnapshots []fix.MarketDataSnapshot
			for _, subscription := range app.FxSession.MarketDataSubscriptions {
				marketDataSnapshot, fxErr := app.FxSession.CtraderMarketDataRequest(app.FxUser, *subscription)
				if fxErr != nil {
					if strings.Contains(fxErr.ErrorMessage, "ALREADY_SUBSCRIBED") {
						continue
					}
					errSource := "MainLoop, switch default. CtraderMarketDataRequest"
					fxErr = app.HandleError(fxErr, errSource)
					if fxErr != nil {
						return
					}
				}
				// if err != nil {
				// 	app.Program.SendColor(fmt.Sprintf("error getting symbol data: %s", err.ErrorMessage), "red")
				// 	continue
				// }
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
				}
				app.UiPositionsDataMap[newUiPosition.copyPositionId] = newUiPosition
			}
			app.Program.Program.Send(PositionMessageSlice(app.UiPositionsDataMap))

			time.Sleep(10 * time.Second)
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
