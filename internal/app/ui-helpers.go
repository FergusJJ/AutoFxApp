package app

import (
	"fmt"
	"time"

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

var (
	textRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E88388"))
	textGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8CC8C"))
	textYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#DBAB79"))
	textFaint  = lipgloss.NewStyle().Faint(true)
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

func (p *AppProgram) SendColor(message, color string) {
	timestamp := time.Now().Format("15:04:05")
	timestamp = textFaint.Render(timestamp)

	switch color {
	case "green":
		message = textGreen.Render(message)
	case "red":
		message = textRed.Render(message)
	case "yellow":
		message = textYellow.Render(message)
	}

	message = fmt.Sprintf("%s - %s", timestamp, message)

	p.Program.Send(FeedUpdate(message))
}
