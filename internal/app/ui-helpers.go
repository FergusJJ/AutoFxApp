package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type (
	feedUpdate      string
	positionsUpdate <-chan []interface{}
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

	rows := []table.Row{
		{"PID332522747", "Buy"},
		{"PID332522747", "Sell"},
	}

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
		BorderForeground(lipgloss.Color("200")).
		BorderBottom(true).
		Bold(false)
	t.SetStyles(style)
	return t
}

func updateMessages() tea.Cmd {
	return nil
}

func updateTable() {

}
