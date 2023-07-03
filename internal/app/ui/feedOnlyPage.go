package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type FeedOnlyPage struct {
	View *tview.Flex
	Feed *Feed
}

func NewFeedOnlyPage() *FeedOnlyPage {
	flex := tview.NewFlex()
	feed := NewFeed("feed page", 5)
	flex.AddItem(feed.Box, 0, 1, true)
	return &FeedOnlyPage{
		View: flex,
		Feed: feed,
	}
}

func (p *FeedOnlyPage) Init(app *AppUi) {
	p.View.Clear()
	p.Feed = NewFeed("feed page", 5)
	p.View.AddItem(p.Feed.Box, 0, 1, true)
	p.Feed.Box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune(tcell.KeyCtrlC) {
			app.App.Stop()
		} else if event.Rune() == '2' {
			app.MainPage.Init(app)
			app.Pages.SwitchToPage("main page")
			// app.MainPage.Log("switched to main page", )
		}
		return event
	})
}

func (p *FeedOnlyPage) Log(message string, color string) {
	p.Feed.Log(message, color)

}
