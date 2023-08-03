package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var tableStyling = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

type Model struct {
	headerMsg      string
	backgroundFeed []string
	table          table.Model
}

// start event loop
func (m Model) Init() tea.Cmd { return nil }

func NewModel(name string) Model {
	m := Model{
		headerMsg:      getHeader(name),
		backgroundFeed: make([]string, 0),
		table:          initialiseTable(),
	}
	return m
}

//message could be a signal from a channel indicating that there are changes to positions,
//

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	type backendUpdate struct {
		//type of update, new feed message, new position, etc.
		//any data that is sent. if feedMessage, this is just a string
		//if this is a position update then a list of positions, because will want to update all
		//positions if it is a getActivePositions message. If not then the list of positions will just be
		//the previous list with some positions added or subtracted
	}
	switch msg := msg.(type) {
	case positionsUpdate:

	case feedUpdate:

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.String() == "q" {
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m Model) View() string {
	//initialise view with header
	s := m.headerMsg

	//append the table on to the message
	s += tableStyling.Render(m.table.View())
	s += "\n"

	//append the feed messages
	for _, msg := range m.backgroundFeed {
		s += fmt.Sprintf("%s\n", msg)
	}

	s += "\nPress q to quit.\n"
	return s
}

func (m *Model) updateMessages(messages ...string) {
	if len(messages) == 0 {
		return
	}
	//if amount of messages < 5; then update messages, if length exceeds 5, then remove oldest message
	for _, msg := range messages {
		if len(m.backgroundFeed) < 5 {
			m.backgroundFeed = append(m.backgroundFeed, msg)
		} else {
			tmpFeed := m.backgroundFeed[1:]
			tmpFeed = append(tmpFeed, msg)
			m.backgroundFeed = tmpFeed
		}

	}
}
