package main

import (
	"fmt"
	"math"
	"math/rand"
	"pollo/internal/app/ui"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

/*
	type ApiMonitorMessage struct {
		CopyPID     string  `json:"copyPID"`
		SymbolID    int     `json:"symbolID"`
		Price       float64 `json:"price"`
		Volume      int     `json:"volume"`
		Direction   string  `json:"direction"`
		MessageType string  `json:"type"` //close or open
	}
*/

var messages = []string{
	"this is message 1",
	"this is message 2",
	"this is message 3",
	"this is message 4",
	"this is message 5",
	"this is message 6",
	"this is message 7",
}

type displayData struct {
	DID         int
	DataMessage string
	Green       bool
}

func main() {
	app := tview.NewApplication()

	feed := ui.NewFeed(5)

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(feed.Box, 0, 1, false)

	go func() {
		// Simulating messages being added to the feed
		for i := 1; i <= 10; i++ {
			message := fmt.Sprintf("Message %d - %s", i, time.Now().Format("2006-01-02 15:04:05"))
			app.QueueUpdateDraw(func() {
				feed.Log(message)
			})
			time.Sleep(time.Second)
		}
	}()

	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		// Close the application when pressing "ESC"
		if event.Key() == tcell.KeyEscape {
			app.Stop()
		}
		return event
	})

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}

func feedData() []displayData {
	m1 := rand.Intn(7)
	m2 := rand.Intn(7)
	d1 := displayData{
		DID:         rand.Int(),
		DataMessage: messages[m1],
		Green:       (int(math.Round(rand.Float64()))) == 1,
	}
	d2 := displayData{
		DID:         rand.Int(),
		DataMessage: messages[m2],
		Green:       (int(math.Round(rand.Float64()))) == 0,
	}
	return []displayData{d1, d2}
}
