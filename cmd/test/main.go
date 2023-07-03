package main

import (
	"fmt"
	"pollo/internal/app/ui"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type displayData struct {
	DID         int
	DataMessage string
	Green       bool
}

func main() {
	app := tview.NewApplication()

	pages := tview.NewPages()

	// First page content
	firstPage := tview.NewFlex().SetDirection(tview.FlexRow)
	firstPageText := tview.NewTextView().SetText("First Page")
	firstPage.AddItem(firstPageText, 0, 1, true)

	// Second page content
	secondPage := tview.NewFlex().SetDirection(tview.FlexRow)
	feed := ui.NewFeed(5)
	secondPage.AddItem(feed.Box, 0, 1, true)

	pages.AddPage("first", firstPage, true, true)
	pages.AddPage("second", secondPage, true, false)

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
		} else if event.Rune() == '1' {
			pages.SwitchToPage("first")
		} else if event.Rune() == '2' {
			pages.SwitchToPage("second")
		}
		return event
	})

	if err := app.SetRoot(pages, true).Run(); err != nil {
		panic(err)
	}
}
