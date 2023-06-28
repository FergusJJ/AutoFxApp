package app

import (
	"encoding/json"
	"log"
	"pollo/pkg/api"
	"pollo/pkg/fix"

	"github.com/fasthttp/websocket"
)

func (app *FxApp) MainLoop() (err *fix.ErrorWithCause) {
	var messageFails int = 0

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
				log.Printf("error sending message to FIX, retrying")
				messageFails++
				if messageFails > 3 {
					return err
				}
				//should never happen
			default:
				log.Fatalf("%+v", err)
			}
		}
		log.Println("logged in")
		loginFinished = true
	}
	app.FxSession.LoggedIn = true

	// var fetchedSecurityList = false
	// for !fetchedSecurityList {
	// 	// log.Println(strings.Split(app.ApiSession.Cid, "_")[1])
	// 	err = app.FxSession.CtraderSecurityList(app.FxUser)
	// 	if err != nil {
	// 		switch err.ErrorCause {
	// 		case fix.ProgramError:
	// 			log.Println("program error")
	// 			return err
	// 		case fix.UserDataError:
	// 			log.Println("user data error")
	// 			return err
	// 		case fix.ConnectionError:
	// 			log.Printf("error sending message to FIX, retrying")
	// 			messageFails++
	// 			if messageFails > 3 {
	// 				return err
	// 			}
	// 			//should never happen
	// 		default:
	// 			log.Fatalf("%+v", err)
	// 		}
	// 	}
	// 	log.Println("got security list")

	// 	fetchedSecurityList = true
	// }
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

			case "CLOSE":
				log.Printf("Got CLOSE:%+v\n", newMessage)
			default:
				log.Fatalln("uknown message type sent to the ", err)
			}
			//unmarshal into json
			//need to check what the message is here, execute message
		default:
			//could maybe just poll current orders here?
			//poll position
			//check for updates, if none continue, else update the table
			fxErr := app.FxSession.CtraderMassStatus(app.FxUser)
			if fxErr != nil {
				log.Panicf("%+v", fxErr)
			}
			log.Fatal()
			continue
		}
	}

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