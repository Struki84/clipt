package clipt

import (
	"fmt"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/struki84/clipt/tui"
)

// make debug mode optional and allow custom debug.log location !
func Render(providers []tui.ChatProvider, session tui.SessionStorage) {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	app := tea.NewProgram(
		tui.NewChatModel(providers, session),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := app.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}

}
