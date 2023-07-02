package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
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

// func main() {
// 	//start standard program init stuff
// 	// log.Println("loading stuff..")
// 	// time.Sleep(3 * time.Second)
// 	//end of normal init stuff

// 	ui.InitUi("Fergus")

// 	//other stuff happening here..
// 	//call goroutines
// 	ticker := time.NewTicker(10 * time.Second)
// 	for {

// 		select {
// 		case <-ticker.C:
// 			//once first bit of data is fetched, load table
// 			_ = feedData()
// 			//update ui here

// 		}
// 	}

// }

func main() {
	app := tview.NewApplication()
	box := tview.NewTextView()
	box.SetBorder(false).SetTitle("Feed")
	box.SetBorderPadding(0, 0, 1, 1) // Remove vertical padding

	go func() {
		messages := []string{}

		// Simulating messages being added to the feed
		for i := 1; i <= 1000; i++ {
			message := fmt.Sprintf("Message %d - %s\n", i, time.Now().Format("2006-01-02 15:04:05"))
			app.QueueUpdateDraw(func() {
				_, _, _, height := box.GetRect()
				messages = append(messages, message)
				if len(messages) >= height-1 {
					messages = messages[1:] // Remove the top message if there is no space
				}
				box.SetText(strings.Join(messages, ""))
			})
			time.Sleep(time.Microsecond * 500)
		}
	}()

	flex := tview.NewFlex().SetDirection(tview.FlexRow)
	flex.AddItem(tview.NewBox(), 0, 4, false)
	flex.AddItem(box, 0, 1, false)

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
