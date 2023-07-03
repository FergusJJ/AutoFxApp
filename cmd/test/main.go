package main

import (
	"fmt"
	"log"
	"math/rand"
	"pollo/internal/app/ui"
	"time"
)

type displayData struct {
	DID         int
	DataMessage string
	Green       bool
}

func main() {

	uiApp := ui.InitializeUi()
	uiApp.SwitchPage("main page")
	uiApp.App.SetRoot(uiApp.Pages, true)
	go func() {
		for {
			// Simulate receiving data
			time.Sleep(2 * time.Second) // Wait for 2 seconds
			dVal := rand.Intn(100)
			// Generate a random message
			message := fmt.Sprintf("Data: %d", dVal)
			color := "red"
			if dVal > 50 {
				color = "green"
			}

			// Log the message in the second page
			uiApp.MainPage.Log(message, color)
			uiApp.App.Draw()
		}
	}()
	go func() {
		entries := []ui.Entry{
			ui.Entry{
				OrderID:      "94878953",
				ClOrdID:      "3f357761-3c81-4a50-b777-b3d38fd6bd92",
				ExecType:     "0",
				OrdStatus:    "0",
				Symbol:       "3",
				Side:         "2",
				TransactTime: "20230630-15:49:07.160",
				OrderQty:     "120002", LeavesQty: "12000",
				PosMaintRptID: "52942663"},
			// Add more entries as needed
		}

		uiApp.MainPage.Table.AddEntry(entries[0])
		// uiApp.MainPage.Table.AddEntry(entries[1])
		// uiApp.MainPage.Table.AddEntry(entries[2])
		// uiApp.MainPage.Table.AddEntry(entries[3])
		// uiApp.MainPage.Table.AddEntry(entries[4])
		// uiApp.MainPage.Table.AddEntry(entries[5])
		// uiApp.MainPage.Table.AddEntry(entries[6])

		// uiApp.MainPage.Table.RemoveEntry(entries[1].OrderID)
		uiApp.App.Draw()

	}()

	if uiErr := uiApp.App.Run(); uiErr != nil {
		log.Fatal("ui err:", uiErr)
	}

}
