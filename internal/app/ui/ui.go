package ui

import (
	"github.com/rivo/tview"
)

/*
Vertical Layout:
1/10 - Need a line saying how to quit the app / The user that is currently copying / maybe other info
6/10 - Need main table area which will show positons.
3/10 - A feed displaying what is happening under the hood

*/

var uiApp = tview.NewApplication()

func InitUi(name string) {
	pages := tview.NewPages()
	MainUILayout := tview.NewFlex().SetDirection(tview.FlexRow)

	infoText := getInfoView(name)
	table := newTable()
	feed := newFeed()
	feed.pushToFeed()
	MainUILayout.SetDirection(tview.FlexRow).
		AddItem(infoText, 0, 1, false).
		AddItem(table.TableView.SetBorder(true), 0, 13, false).
		AddItem(feed.FeedView.SetBorder(true), 0, 6, false)

	pages.AddPage("Main UI", MainUILayout, true, true)
	if err := uiApp.SetRoot(pages, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
