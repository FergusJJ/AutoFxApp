package ui

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Feed struct {
	Box       *tview.TextView
	maxLines  int
	messages  []string
	lineCount int
}

func NewFeed(maxLines int) *Feed {
	box := tview.NewTextView()
	box.SetBorder(true).SetTitle("Feed")
	box.SetBorderPadding(0, 0, 1, 1) // Remove vertical padding

	return &Feed{
		Box:       box,
		maxLines:  maxLines,
		messages:  make([]string, 0),
		lineCount: 0,
	}
}

func (f *Feed) Log(message string, color ...tcell.Color) {
	if len(color) == 0 {
		//color should just be grey
	}
	f.messages = append(f.messages, message)
	f.lineCount++

	if f.lineCount > f.maxLines {
		f.lineCount--
		f.messages = f.messages[1:]
	}

	f.Box.SetText(strings.Join(f.messages, "\n"))
}
