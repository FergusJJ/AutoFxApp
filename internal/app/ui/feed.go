package ui

import (
	"fmt"
	"strings"

	"github.com/rivo/tview"
)

type Feed struct {
	Box       *tview.TextView
	maxLines  int
	messages  []string
	lineCount int
}

func NewFeed(title string, maxLines int) *Feed {
	box := tview.NewTextView()
	box.SetDynamicColors(true)
	box.SetBorder(true)
	box.SetTitle(title)
	box.SetBorderPadding(0, 0, 1, 1) // Remove vertical padding

	return &Feed{
		Box:       box,
		maxLines:  maxLines,
		messages:  make([]string, 0),
		lineCount: 0,
	}
}

func (f *Feed) Log(message string, color string) {
	coloredMessage := fmt.Sprintf("[%s]%s[white]", color, message)
	f.messages = append(f.messages, coloredMessage)
	f.lineCount++
	if f.maxLines == -1 {
		f.Box.SetText(strings.Join(f.messages, "\n"))
		return
	}
	if f.lineCount > f.maxLines {
		f.lineCount--
		f.messages = f.messages[1:]
	}

	f.Box.SetText(strings.Join(f.messages, "\n"))
}
