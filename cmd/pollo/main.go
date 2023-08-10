package main

import (
	"errors"
	"log"
	"os"
	"pollo/config"
	"pollo/internal/app"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
	//shutdown.Gracefully()
}

func start() (func(), error) {
	done := make(chan struct{})
	backendErrChan := make(chan error, 1)
	errChan := make(chan error, 1)
	app, cleanup, err := initialiseProgram()
	if err != nil {
		return nil, err
	}
	go func() {
		defer close(done)
		app.MainLoop()
	}()
	errChan <- func() error {
		_, err := app.Progam.Run()
		if err != nil {
			return err
		}
		return nil
	}()
	select {
	case err := <-errChan:
		if err != nil {
			log.Println("ui error:", err)
			return cleanup, err
		}
	case err := <-backendErrChan:
		if err != nil {
			if err != nil {
				log.Println("backend error:", err)
				return cleanup, err
			}
		}
	case <-done:
	}
	return func() {
		log.Println("running cleanup...")
		cleanup()
	}, nil
}

func initialiseProgram() (*app.FxApp, func(), error) {

	App := &app.FxApp{}
	App.Progam = tea.NewProgram(app.NewModel("fergus"))
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
	timeout := time.Duration(10 * time.Second)
	fxPriceClient := fix.NewTCPClient(App.FxUser.HostName, fix.PricePort, timeout, 4096)
	if err = fxPriceClient.Dial(); err != nil {
		return nil, func() {
			App.CloseExistingConnections()
			log.Println("closed existing connections")
		}, err
	}
	fxTradeClient := fix.NewTCPClient(App.FxUser.HostName, fix.TradePort, timeout, 4096)
	if err = fxTradeClient.Dial(); err != nil {
		return nil, func() {
			App.CloseExistingConnections()
			log.Println("closed existing connections")
		}, err
	}

	App.FxSession.TradeClient = fxTradeClient
	App.FxSession.PriceClient = fxPriceClient
	App.FxSession.MessageSequenceNumber = 1
	App.FxSession.LoggedIn = false
	log.Println("connected to fix api")
	//FxSesion Done

	//ApiSession Start
	cid, err := api.CheckLicense(App.LicenseKey)
	if err != nil {
		return nil, func() {
			App.CloseExistingConnections()
			log.Println("closed existing connections")
		}, err
	}
	log.Println("license verified")
	App.ApiSession.Cid = cid

	apiConn, err := api.CreateApiConnection(App.ApiSession.Cid, pools)
	if err != nil {
		return nil, func() {
			App.CloseExistingConnections()
			log.Println("closed existing connections")
		}, err
	}
	App.ApiSession.Client.Connection = apiConn
	App.ApiSession.Client.CurrentMessage = make(chan []byte)
	log.Println("connected to internal api")

	//ApiSesion Done

	//start the actual program, initilse monitoring client via ws,
	//start the function that will be responsible for sending fix api requests
	return App, func() {
		//cleanup operations, i.e. close api ws connection, close fix api session
		App.CloseExistingConnections()
		log.Println("closed existing connections")
	}, nil
}
