package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func getInfoView(name string) *tview.Flex {

	var nameText = tview.NewTextView().
		SetTextColor(tcell.ColorGreen).
		SetText(fmt.Sprintf("Welcome %s!", name))
	var infoText = tview.NewTextView().
		SetTextColor(tcell.ColorGray).
		SetText("(q) to quit")
	nameText.SetBorderPadding(0, 0, 1, 0)
	infoText.SetBorderPadding(0, 0, 1, 0)

	testBox := tview.NewFlex()
	testBox.SetDirection(tview.FlexRow)
	testBox.AddItem(nameText, 0, 1, false).AddItem(infoText, 0, 1, false)
	return testBox
}
