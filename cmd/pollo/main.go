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

	apexLog "github.com/apex/log"

	"github.com/apex/log/handlers/cli"
)

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
	appErr := app.MainLoop()
	if appErr != nil {
		log.Printf("%v: %s\n", appErr.ErrorCause, appErr.ErrorMessage)
		return func() {
			log.Println("running cleanup...")
			cleanup()
		}, nil
	}
	return func() {
		log.Println("running cleanup...")
		cleanup()
	}, nil
}

func initialiseProgram() (*app.FxApp, func(), error) {
	app := &app.FxApp{}

	//AppLogger Start
	logger := &apexLog.Logger{
		Handler: cli.New(os.Stdout),
		Level:   1,
	}
	app.AppLogger = logger
	//AppLogger Done

	//FxUser & Lisence Key Start
	fxUser, err := config.LoadDataFromJson()
	if err != nil {
		return nil, nil, err
	}
	app.FxUser = *fxUser

	licenseKey, err := config.LoadLicenseKeyFromJson()
	if err != nil {
		return nil, nil, err
	}
	if licenseKey == "" {
		err = errors.New("licenseKey is empty, update settings.json")
		return nil, nil, err
	}
	app.LicenseKey = licenseKey
	//FxUser & Lisence Key Done

	//FxSession Start

	fxConn, err := fix.CreateConnection(app.FxUser.HostName, fix.TradePort)
	if err != nil {
		//cleanup should involve closing fx connection
		return nil, func() {
			app.CloseExistingConnections()
		}, err
	}
	app.FxSession.Connection = fxConn
	app.FxSession.MessageSequenceNumber = 1
	app.FxSession.LoggedIn = false
	log.Println("connected to fix api")
	//FxSesion Done

	//ApiSession Start
	cid, err := api.CheckLicense(app.LicenseKey)
	if err != nil {
		return nil, nil, err
	}
	log.Println("license verified")
	app.ApiSession.Cid = cid

	apiConn, err := api.CreateApiConnection(app.ApiSession.Cid)
	if err != nil {
		return nil, func() {
			app.CloseExistingConnections()
			log.Println("closed existing connections")
		}, err
	}
	app.ApiSession.Client.Connection = apiConn
	app.ApiSession.Client.CurrentMessage = make(chan []byte)
	log.Println("connected to internal api")

	//ApiSesion Done

	//start the actual program, initilse monitoring client via ws,
	//start the function that will be responsible for sending fix api requests
	return app, func() {
		//cleanup operations, i.e. close api ws connection, close fix api session
		app.CloseExistingConnections()
		log.Println("closed existing connections")
	}, nil
}