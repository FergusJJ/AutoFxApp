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
	header := fmt.Sprintf("Welcome %s | CopyFX - Version 0.1\n\n", name)
	return header
}

func initialiseTable() table.Model {

	cols := []table.Column{
		{Title: "PID", Width: 12},
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
