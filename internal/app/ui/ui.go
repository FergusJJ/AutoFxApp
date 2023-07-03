package ui

import "github.com/rivo/tview"

func InitUi() *tview.Application {
	var app = tview.NewApplication()
	pages := tview.NewPages()
	pages.AddPage("Loading", newLoadingScreen(), false, true)
	return app
}

func newLoadingScreen() *tview.Flex {
	flexLayout := tview.NewFlex()
	flexLayout.SetDirection(tview.FlexRow)
	flexLayout.SetBorder(false)
	return flexLayout
}

func newMainScreen() {

}
