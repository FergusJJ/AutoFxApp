package main

import (
	"fmt"
	"os"
	"pollo/internal/app"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

//for commands, can pass a function that does the command to another wrapper function
/*

func cmdWithArg(funcToRun func()) tea.Cmd {
    return func() tea.Msg {
				message := funcToRun()
				//check message
				//return message
        return someMsg{id: id}
    }
}

*/

func main() {
	p := tea.NewProgram(app.NewModel("fergus"))
	go func() {
		totalPositions := []app.PositionMessage{}

		for {
			pause := time.Duration(5000) * time.Millisecond // nolint:gosec
			time.Sleep(pause)
			newMessage := app.PositionMessage{ID: pause.String(), Direction: "SELL"}
			totalPositions = append(totalPositions, newMessage)
			p.Send(app.PositionMessageSlice(totalPositions))
			p.Send(app.FeedUpdate("hello"))
		}
	}()
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
