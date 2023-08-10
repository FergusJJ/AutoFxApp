package app

import (
	"encoding/json"
	"fmt"
	"log"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"time"

	"github.com/fasthttp/websocket"
)

// going to have to log errors here
func (app *FxApp) MainLoop() (err *fix.ErrorWithCause) {
	var messageFails int = 0
	messageFails = 0
	var loginFinished bool = false
	for !loginFinished {
		err = app.FxSession.CtraderLogin(app.FxUser)
		if err != nil {
			app.Progam.Send(FeedUpdate(err.ErrorMessage))
			switch err.ErrorCause {
			case fix.ProgramError:
				return err
			case fix.UserDataError:
				return err
			case fix.ConnectionError:
				app.Progam.Send(FeedUpdate("error sending message to FIX, retrying"))
				messageFails++
				if messageFails > 3 {
					return err
				}

			default:
				log.Fatalf("%+v", err)
			}
			continue
		}
		app.Progam.Send(FeedUpdate("Logged in to ctrader"))
		loginFinished = true
	}
	// app.FxSession

	app.FxSession.LoggedIn = true
	app.FxSession.GotSecurityList = true
	//need to start function that will monitor here:
	go app.ApiSession.ListenForMessages()
	// success, newPos := app.OpenPosition(nil)
	// if !success {
	// 	log.Fatalf("unsuccessful")
	// }
	symbol := fmt.Sprint(1)
	app.FxSession.NewMarketDataSubscription(symbol)
	// app.Progam.Send(FeedUpdate(fmt.Sprint(*newPos)))
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
				app.FxSession.NewMarketDataSubscription(symbol)
			case "CLOSE":
				app.Progam.Send(FeedUpdate(fmt.Sprintf("Got CLOSE:%+v\n", newMessage)))

				// pid, err := app.ClosePosition(newMessage)
				// if err != nil {
				// 	app.Progam.Send(FeedUpdate(fmt.Sprint("close position error: ", err)))
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
			app.Progam.Send(FeedUpdate("updating position data"))
			var marketDataSnapshots []fix.MarketDataSnapshot
			for _, subscription := range app.FxSession.MarketDataSubscriptions {
				marketDataSnapshot, err := app.FxSession.CtraderMarketDataRequest(app.FxUser, *subscription)
				if err != nil {
					app.Progam.Send(FeedUpdate(fmt.Sprintf("error getting symbol data: %s", err.ErrorMessage)))
					continue
				}
				marketDataSnapshots = append(marketDataSnapshots, marketDataSnapshot...)
			}

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
			app.Progam.Send(FeedUpdate(fmt.Sprintf("error closing connection to FIX: %s", err.Error())))
		}
		app.FxSession.TradeClient = nil
	}

	if app.FxSession.PriceClient != nil {
		err := app.FxSession.PriceClient.Close()
		if err != nil {
			app.Progam.Send(FeedUpdate(fmt.Sprintf("error closing connection to FIX: %s", err.Error())))
		}
		app.FxSession.PriceClient = nil

	}

	if app.ApiSession.Client.Connection != nil {
		err := app.ApiSession.Client.Connection.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		if err != nil {
			app.Progam.Send(FeedUpdate(fmt.Sprintf("error writing close to API: %s", err.Error())))
		}
		//this will never panic, so logging error is fine
		err = app.ApiSession.Client.Connection.Close()
		if err != nil {
			app.Progam.Send(FeedUpdate(fmt.Sprintf("error closing connection to API: %s", err.Error())))
		}
		app.ApiSession.Client.Connection = nil
	}
}
