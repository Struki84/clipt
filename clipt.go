package clipt

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/struki84/clipt/tui"
)

func SetProviders(providers []tui.ChatProvider) {
	tui.Providers = providers
}

func Render(provider tui.ChatProvider) {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	app := tea.NewProgram(
		tui.NewChatModel(tui.Providers),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := app.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
