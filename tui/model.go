package tui

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/struki84/clipt/tui/schema"
)

type ChatModel struct {
	Layout    LayoutView
	Storage   schema.SessionStorage
	Providers []schema.ChatProvider
}

func NewChatModel(providers []schema.ChatProvider, storage schema.SessionStorage) ChatModel {
	provider := providers[0]
	model := ChatModel{
		Layout:    NewLayoutView(provider),
		Providers: providers,
		Storage:   storage,
	}

	if model.Storage != nil {
		currentSession, err := model.Storage.LoadRecentSession()
		if err != nil {
			log.Printf("Could not load chat history.")
		}

		model.Layout.Chat.Session = currentSession
		model.Layout.Chat.Msgs = currentSession.Msgs
	}

	return model
}

func (model ChatModel) Init() tea.Cmd {
	return tea.Batch(model.Layout.Init())
}

func (model ChatModel) View() string {
	return model.Layout.View()
}

func (model ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if model.Layout.Menu.Active {
				model.Layout.Menu = model.Layout.Menu.Close()
				model.Layout.Chat.Input.SetValue("")
				return model, nil
			}
		case tea.KeyCtrlC:
			return model, tea.Quit
		}
	case schema.ExecuteCmd:
		return msg.Cmd.Execute(model)
	}

	cmds := []tea.Cmd{}
	l, cmd := model.Layout.Update(msg)
	model.Layout = l.(LayoutView)

	cmds = append(cmds, cmd)
	return model, tea.Batch(cmds...)
}
