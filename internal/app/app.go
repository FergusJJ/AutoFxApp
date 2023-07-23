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
	// time.Sleep(time.Second * 3)
	messageFails = 0
	var loginFinished bool = false
	for !loginFinished {
		err = app.FxSession.CtraderLogin(app.FxUser)
		if err != nil {
			switch err.ErrorCause {
			case fix.ProgramError:
				return err
			case fix.UserDataError:
				return err
			case fix.ConnectionError:
				app.ScreenWriter.Write("error sending message to FIX, retrying")
				messageFails++
				if messageFails > 3 {
					return err
				}
				//should never happen
			default:
				log.Fatalf("%+v", err)
			}
		}
		// app.UI.MainPage.Log("Logged in to ctrader", "green")
		app.ScreenWriter.Write("logged in")
		loginFinished = true
	}
	app.FxSession.LoggedIn = true
	app.FxSession.GotSecurityList = true
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
				log.Printf("Got OPEN:%+v\n", newMessage)
				pid, err := app.OpenPosition(newMessage)
				if err != nil {
					app.ScreenWriter.Write(fmt.Sprint("open position error: ", err))
					continue
				}
				//send pid to db with the copy pid
				app.ScreenWriter.Write(fmt.Sprint("sending pid:", pid))
			case "CLOSE":
				log.Printf("Got CLOSE:%+v\n", newMessage)
				pid, err := app.ClosePosition(newMessage)
				if err != nil {
					app.ScreenWriter.Write(fmt.Sprint("close position error: ", err))
					continue
				}
				//send pid to db with the copy pid
				app.ScreenWriter.Write(fmt.Sprint("sending pid:", pid))
			default:
				log.Fatalln("uknown message type sent to the ", err)
			}
			//unmarshal into json
			//need to check what the message is here, execute message
		default:
			//could maybe just poll current orders here?
			//poll position
			//check for updates, if none continue, else update the table
			// fxErr := app.FxSession.CtraderRequestForPositions(app.FxUser)
			// if fxErr != nil {
			// 	log.Panicf("%+v", fxErr)
			// }
			time.Sleep(2 * time.Second)
			continue
		}
	}
	return nil
}

func (app *FxApp) CloseExistingConnections() {
	if app.FxSession.Connection != nil {
		//this will never panic, so logging error is fine
		err := app.FxSession.Connection.Close()
		if err != nil {
			log.Printf("error closing connection to FIX: %s", err.Error())
		}
	}

	if app.ApiSession.Client.Connection != nil {
		err := app.ApiSession.Client.Connection.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		if err != nil {
			log.Printf("error writing close to API: %s", err.Error())
		}
		//this will never panic, so logging error is fine
		err = app.ApiSession.Client.Connection.Close()
		if err != nil {
			log.Printf("error closing connection to API: %s", err.Error())
		}
	}
}
