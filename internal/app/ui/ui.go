package ui

import "github.com/rivo/tview"

type AppUi struct {
	App         *tview.Application
	Pages       *tview.Pages
	FeedPage    *FeedOnlyPage
	MainPage    *MainPage
	CurrentPage string
}

func (ui *AppUi) SwitchPage(newPage string) {
	switch newPage {
	case "feed page":
		ui.FeedPage.Init(ui)
		ui.Pages.SwitchToPage(newPage)
		ui.CurrentPage = newPage
	case "main page":
		ui.MainPage.Init(ui)
		ui.Pages.SwitchToPage(newPage)
		ui.CurrentPage = newPage
	}
}

func InitializeUi() *AppUi {
	ui := &AppUi{}
	ui.App = tview.NewApplication()
	ui.Pages = tview.NewPages()

	ui.FeedPage = NewFeedOnlyPage()
	ui.FeedPage.Init(ui)
	ui.Pages.AddPage("feed page", ui.FeedPage.View, true, false)

	ui.MainPage = NewMainPage()
	ui.MainPage.Init(ui)
	ui.Pages.AddPage("main page", ui.MainPage.View, true, false)

	return ui
}
