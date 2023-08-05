package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type PositionMessage struct {
	ID        string
	Direction string
}

type (
	QuitApp              struct{}
	FeedUpdate           string
	PositionMessageSlice []PositionMessage
)

func getHeader(name string) string {

	var headerStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4"))

	header := headerStyle.Render(fmt.Sprintf("Welcome %s | CopyFX - Version 0.1", name))
	return header
}

func initialiseTable() table.Model {

	var colStyling = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#0000ff"))
	thisPID := colStyling.Render("PID")
	cols := []table.Column{
		{Title: thisPID, Width: 12},
		{Title: "DIRECTION", Width: 10},
	}
	rows := []table.Row{}
	tableOpts := []table.Option{
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithColumns(cols),
		table.WithRows(rows),
	}
	t := table.New(
		tableOpts...,
	)
	style := table.DefaultStyles()
	style.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	t.SetStyles(style)
	return t
}
