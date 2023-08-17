package main

import (
	"errors"
	"log"
	"os"
	"pollo/config"
	"pollo/internal/app"
	"pollo/internal/logs"
	"pollo/pkg/api"
	"pollo/pkg/fix"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/inancgumus/screen"
)

/*
Make sure that if there is an interruption in network connection, to exit, will probably be detected by
Websocket api monitor thingy, this will mean that it is less likely to have to retry ctrader messages due to  interruptions

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
	errChan := make(chan error, 1)
	app, cleanup, err := initialiseProgram()
	if err != nil {
		return nil, err
	}
	screen.Clear()
	screen.MoveTopLeft()
	go func() {
		defer close(done)
		app.MainLoop()
		time.Sleep(10 * time.Second)
		//need some sort of pause here on the ui, block until the user exits
		app.Program.Program.Send(tea.QuitMsg{})
	}()
	errChan <- func() error {
		_, err := app.Program.Program.Run()
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

	case <-done:
	}
	return func() {
		log.Println("running cleanup...")
		cleanup()
	}, nil
}

func initialiseProgram() (*app.FxApp, func(), error) {

	App := &app.FxApp{}

	App.Program.Program = tea.NewProgram(app.NewModel("fergus"), tea.WithAltScreen())
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
	App.FxSession.TradeMessageSequenceNumber = 1
	App.FxSession.PriceMessageSequenceNumber = 1
	App.FxSession.LoggedIn = false
	App.FxSession.MarketDataSubscriptions = make(map[string]*fix.MarketDataSubscription)
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
	App.ApiSession.LicenseKey = App.LicenseKey
	var apiErr *api.ApiError
	var deadline = time.Now().UnixMilli() + 600000
	for time.Now().UnixMilli() < deadline {
		apiErr = App.ApiSession.FetchApiAuth()
		if apiErr == nil {
			break
		}
		if apiErr.ShouldExit {
			logs.SendApplicationLog(apiErr.ErrorMessage, App.LicenseKey)
			return nil, func() {
				App.CloseExistingConnections()
				log.Println("closed existing connections")
			}, err
		}
		log.Println(apiErr.UserMessage)
		time.Sleep(10 * time.Second)
	}
	if apiErr != nil {
		return nil, func() {
			App.CloseExistingConnections()
			log.Println("closed existing connections")
		}, err
	}

	log.Println("session authorised")

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
