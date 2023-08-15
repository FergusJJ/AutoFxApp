package app

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
)

type (
	QuitApp              struct{}
	FeedUpdate           string
	PositionMessageSlice map[string]uiPositionData
	tickMsg              int
)

var (
	textRed    = lipgloss.NewStyle().Foreground(lipgloss.Color("#E88388"))
	textGreen  = lipgloss.NewStyle().Foreground(lipgloss.Color("#A8CC8C"))
	textYellow = lipgloss.NewStyle().Foreground(lipgloss.Color("#DBAB79"))
	textFaint  = lipgloss.NewStyle().Faint(true)
)

func getHeader(name string) string {
	header := titleStyle.Render(fmt.Sprintf("Welcome %s | CopyFX - Version 0.1", name))
	return header
}

func initialiseTable(rows ...table.Row) table.Model {

	var colStyling = lipgloss.NewStyle().
		Bold(true)

	userPID := colStyling.Render("PID")
	copyPID := colStyling.Render("COPY PID")
	symbol := colStyling.Render("SYMBOL")
	timestamp := colStyling.Render("OPENED")
	side := colStyling.Render("SIDE")
	volume := colStyling.Render("VOLUME")
	grossProfit := colStyling.Render("GROSS")
	entryPrice := colStyling.Render("ENTRY")
	currentPrice := colStyling.Render("CURRENT PRICE")
	cols := []table.Column{
		{Title: userPID, Width: 15},
		{Title: copyPID, Width: 15},
		{Title: timestamp, Width: 15},
		{Title: symbol, Width: 15},
		{Title: volume, Width: 15},
		{Title: side, Width: 15},
		{Title: entryPrice, Width: 15},
		{Title: currentPrice, Width: 15},
		{Title: grossProfit, Width: 15},
	}
	// rows := []table.Row{}
	tableOpts := []table.Option{
		table.WithFocused(true),
		table.WithHeight(15),
		table.WithColumns(cols),
		table.WithRows(rows),
	}
	t := table.New(
		tableOpts...,
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)

	s.Cell.AlignHorizontal(lipgloss.Left)
	// s.Cell.Color
	// s.Selected = s.Selected.
	// 	Foreground(lipgloss.Color("229")).
	// 	Background(lipgloss.Color("57")).
	// 	Bold(false)
	t.SetStyles(s)

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

/*
error getting marketData: business reject: ALREADY_SUBSCRIBED:An attempt to subscribe twice
*/
