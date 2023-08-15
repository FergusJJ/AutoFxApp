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
	Padding(0, 1)

// var lossStyle = lipgloss.NewStyle().
// 	Foreground(lipgloss.Color("#ff0000"))
// var profitStyle = lipgloss.NewStyle().
// 	Foreground(lipgloss.Color("#00ff00"))

var tableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.RoundedBorder()).
	BorderForeground(lipgloss.Color("240")).
	Padding(0, 1)

type Model struct {
	headerMsg      string
	backgroundFeed []FeedUpdate
	positions      []uiPositionData
	table          table.Model
}

// start event loop
func (m Model) Init() tea.Cmd { return nil }

func NewModel(name string) Model {
	m := Model{
		headerMsg:      getHeader(name),
		backgroundFeed: make([]FeedUpdate, 0),
		positions:      make([]uiPositionData, 0),
		table:          initialiseTable(),
	}
	return m
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//maybe need an event listener for scrren resize which will re-init the table
	switch msg := msg.(type) {
	case PositionMessageSlice:
		// return m, tea.Quit
		tmpPositions := []uiPositionData{}
		for copyPID, v := range msg {
			v.copyPositionId = copyPID
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
	// case tickMsg:
	// 	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	// 	if w != m.w || h != m.h {
	// 		m.updateSize(w, h)
	// 	}
	// 	return m, tea.Batch(tick, func() tea.Msg { return tea.WindowSizeMsg{Width: w, Height: h} })

	case tea.QuitMsg:
		return m, tea.Quit
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
		if v.isProfit {
			rows = append(rows, []string{v.positionId, v.copyPositionId, v.timestamp, v.symbol, v.volume, v.side, v.entryPrice, v.currentPrice, v.grossProfit})
		} else {
			rows = append(rows, []string{v.positionId, v.copyPositionId, v.timestamp, v.symbol, v.volume, v.side, v.entryPrice, v.currentPrice, v.grossProfit})
		}

	}
	m.table = initialiseTable(rows...)
	// m.table.SetRows(rows)
	s += tableStyle.Render(m.table.View())
	s += "\n"

	//append the feed messages
	tmp := ""
	for _, msg := range m.backgroundFeed {
		tmp += fmt.Sprintf("%s\n", msg)

	}
	s += tmp
	s += lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("\nPress q to quit.\n")
	return appPadding.Render(s)
}

func (m *Model) updateMessages(messages ...FeedUpdate) {
	if len(messages) == 0 {
		return
	}
	//if amount of messages < 5; then update messages, if length exceeds 5, then remove oldest message
	for _, msg := range messages {
		if len(m.backgroundFeed) < 10 {
			m.backgroundFeed = append(m.backgroundFeed, msg)
		} else {
			tmpFeed := m.backgroundFeed[1:]
			tmpFeed = append(tmpFeed, msg)
			m.backgroundFeed = tmpFeed
		}

	}
}
