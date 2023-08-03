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
	fmt.Println("initialising stuff")
	time.Sleep(time.Second * 5)
	p := tea.NewProgram(app.NewModel("fergus"))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
