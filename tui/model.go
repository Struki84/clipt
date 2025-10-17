package tui

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// func SetCommands(cmds []list.Item) {
// 	Commands = cmds
// }

func SetProviders(providers []ChatProvider) {
	Providers = providers
}

func Render(provider ChatProvider) {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	app := tea.NewProgram(
		NewLayoutView(provider),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := app.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
