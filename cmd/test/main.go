package main

import (
	"log"
	"math"
	"math/rand"
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

var app = tview.NewApplication()
var text = tview.NewTextView().
	SetTextColor(tcell.ColorGreen).
	SetText("(q) to quit")

func main() {
	//start standard program init stuff
	log.Println("loading stuff..")
	time.Sleep(3 * time.Second)
	//end of normal init stuff

	if err := app.SetRoot(text, true).EnableMouse(true).Run(); err != nil {
		log.Fatal(err)
	}

	//other stuff happening here..
	//call goroutines
	ticker := time.NewTicker(10 * time.Second)
	for {

		select {
		case <-ticker.C:
			//once first bit of data is fetched, load table
			_ = feedData()
			//update ui here

		}
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
