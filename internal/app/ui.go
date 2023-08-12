package app

import (
	"fmt"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var appPadding = lipgloss.NewStyle().Padding(2, 1)

var titleStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#FFFDF5")).
	Background(lipgloss.Color("#25A065")).
	Padding(0, 1)

var tableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(0, 1)

type Model struct {
	headerMsg      string
	backgroundFeed []FeedUpdate
	positions      []PositionMessage
	table          table.Model
}

// start event loop
func (m Model) Init() tea.Cmd { return nil }

func NewModel(name string) Model {
	m := Model{
		headerMsg:      getHeader(name),
		backgroundFeed: make([]FeedUpdate, 0),
		positions:      make([]PositionMessage, 0),
		table:          initialiseTable(),
	}
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case PositionMessageSlice:
		// return m, tea.Quit
		tmpPositions := []PositionMessage{}
		for _, v := range msg {
			tmpPositions = append(tmpPositions, v)
		}
		m.positions = tmpPositions
		return m, nil
	case FeedUpdate:
		m.updateMessages(FeedUpdate(msg))
		return m, nil
	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
		if msg.String() == "q" {
			return m, tea.Quit
		}

	default:
	}
	return m, nil
}

func (m Model) View() string {
	//initialise view with header

	// s := m.headerMsg
	s := titleStyle.Render(m.headerMsg)
	s += "\n"

	//append the table on to the message
	rows := []table.Row{}
	for _, v := range m.positions {
		rows = append(rows, []string{v.ID, v.Direction})
	}
	m.table.SetRows(rows)
	s += tableStyle.Render(m.table.View())
	s += "\n"

	//append the feed messages
	for _, msg := range m.backgroundFeed {
		s += fmt.Sprintf("%s\n", msg)
	}

	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\nPress q to quit.\n")
	return appPadding.Render(s)
}

func (m *Model) updateMessages(messages ...FeedUpdate) {
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
