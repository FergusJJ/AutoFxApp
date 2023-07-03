package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MainPage struct {
	View  *tview.Flex
	Feed  *Feed
	Table *PositionTable
}

func NewMainPage() *MainPage {
	flex := tview.NewFlex()
	feed := NewFeed("main page", 5)
	table := NewTable([]string{
		"OrderID",
		"ClOrdID",
		"ExecType",
		"OrdStatus",
		"Symbol",
		"Side",
		"TransactTime",
		"OrderQty",
		"LeavesQty",
		"PosMaintRptID"})
	table.Init()
	flex.SetDirection(tview.FlexRow)

	return &MainPage{
		View:  flex,
		Feed:  feed,
		Table: table,
	}
}

func (p *MainPage) Init(app *AppUi) {
	p.View.Clear()
	p.Feed = NewFeed("main page", 5)
	p.View.AddItem(p.Table.TableView, 0, 1, true)
	p.View.AddItem(p.Feed.Box, 0, 1, true)
	// p.View.AddItem(p.Table.TableView.Box, 0, 1, true)
	p.Feed.Box.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Rune() == rune(tcell.KeyCtrlC) {
			app.App.Stop()
		} else if event.Rune() == '1' {
			app.FeedPage.Init(app)
			app.Pages.SwitchToPage("feed page")
		}
		return event
	})
}

func (p *MainPage) Log(message string, color string) {
	p.Feed.Log(message, color)

}
