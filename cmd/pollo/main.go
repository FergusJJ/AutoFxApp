package main

import (
	"errors"
	"log"
	"os"
	"pollo/config"
	"pollo/internal/app"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"pollo/pkg/shutdown"
)

/*
When getting position report, fix sends multiple messages, these  don't always come in at once, so aren't read into []byte on one request so for 4
positions may get 3 the first time then 1 the second time, this will mean that all 4 are fetched eventually but may want to make sure that server
has completely finished sending the messages first
*/

func main() {
	var exitCode int
	defer func() {
		os.Exit(exitCode)
	}()

	cleanup, err := start()
	defer cleanup()
	if err != nil {
		log.Printf("error: %s", err.Error())
		exitCode = 1
		return
	}

	shutdown.Gracefully()
}

func start() (func(), error) {
	//verify with whop
	//load user data, check they have access

	app, cleanup, err := initialiseProgram()
	if err != nil {
		return nil, err
	}
	/*
		Run app, all the app specific variables should be kept within this struct
		May want to add an "isLoggedIn" field to app, in case it is necessary to always logout of FIX
	*/
	// go app.MainLoop()
	// if err := app.UI.App.Run(); err != nil {
	// 	log.Println("ui error:", err)
	// }
	appErr := app.MainLoop()
	if appErr != nil {
		log.Printf("%v: %s\n", appErr.ErrorCause, appErr.ErrorMessage)
		return func() {
			//close app in cleanup
			app.ScreenWriter.Write("running cleanup...")
			cleanup()
		}, nil
	}
	return func() {
		app.ScreenWriter.Write("running cleanup...")
		cleanup()
	}, nil
}

func initialiseProgram() (*app.FxApp, func(), error) {

	App := &app.FxApp{}
	App.ScreenWriter = app.NewScreenWriter(5)

	//FxUser & Lisence Key Start
	fxUser, err := config.LoadDataFromJson()
	if err != nil {
		return nil, nil, err
	}
	App.FxUser = *fxUser

	licenseKey, pools, err := config.LoadSettingsFromJson()
	if err != nil {
		return nil, nil, err
	}
	if licenseKey == "" {
		err = errors.New("licenseKey is empty, update settings.json")
		return nil, nil, err
	}
	App.LicenseKey = licenseKey
	App.ApiSession.Pools = pools
	//FxUser & Lisence Key Done

	//FxSession Start

	fxConn, err := fix.CreateConnection(App.FxUser.HostName, fix.TradePort)
	if err != nil {
		//cleanup should involve closing fx connection
		return nil, func() {
			App.CloseExistingConnections()
		}, err
	}
	App.FxSession.Connection = fxConn
	App.FxSession.MessageSequenceNumber = 1
	App.FxSession.LoggedIn = false
	App.ScreenWriter.Write("connected to fix api")
	//FxSesion Done

	//ApiSession Start
	cid, err := api.CheckLicense(App.LicenseKey)
	if err != nil {
		return nil, nil, err
	}
	App.ScreenWriter.Write("license verified")
	App.ApiSession.Cid = cid

	apiConn, err := api.CreateApiConnection(App.ApiSession.Cid, pools)
	if err != nil {
		return nil, func() {
			App.CloseExistingConnections()
			App.ScreenWriter.Write("closed existing connections")
		}, err
	}
	App.ApiSession.Client.Connection = apiConn
	App.ApiSession.Client.CurrentMessage = make(chan []byte)
	App.ScreenWriter.Write("connected to internal api")

	//ApiSesion Done

	//start the actual program, initilse monitoring client via ws,
	//start the function that will be responsible for sending fix api requests
	return App, func() {
		//cleanup operations, i.e. close api ws connection, close fix api session
		App.CloseExistingConnections()
		App.ScreenWriter.Write("closed existing connections")
	}, nil
}
