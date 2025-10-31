package tui

import (
	"context"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type StreamMsg struct {
	Chunk string
}

type ChatModel struct {
	Layout    LayoutView
	Providers []ChatProvider
	Provider  ChatProvider
	Storage   SessionStorage
	Stream    chan string
}

func NewChatModel(providers []ChatProvider, storage SessionStorage) ChatModel {
	provider := providers[0]
	model := ChatModel{
		Layout:    NewLayoutView(provider),
		Providers: providers,
		Provider:  provider,
		Storage:   storage,
		Stream:    make(chan string),
	}

	if model.Storage != nil {
		currentSession, err := model.Storage.LoadRecentSession()
		if err != nil {
			log.Printf("Could not load chat history.")
		}

		model.Layout.Session = currentSession
		model.Layout.Msgs = currentSession.Msgs
	}

	return model
}

func (model ChatModel) Init() tea.Cmd {
	log.Printf("ChatMode.Init()")

	model.Provider.Stream(context.TODO(), func(ctx context.Context, chunk []byte) error {
		model.Stream <- string(chunk)
		return nil
	})

	return tea.Batch(model.Layout.Init(), model.HandleStream)
}

func (model ChatModel) View() string {
	return model.Layout.View()
}

func (model ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if model.Layout.MenuActive {
				model.Layout.MenuActive = false
				model.Layout.ChatInput.SetValue("")
				model.Layout.CurrentMenuItems = model.Layout.MenuItems
				model.Layout.FilteredMenuItems = model.Layout.MenuItems
				return model, nil
			}
		case tea.KeyCtrlC:
			return model, tea.Quit
		case tea.KeyEnter:
			if model.Layout.MenuActive {
				if len(model.Layout.FilteredMenuItems) > 0 {
					selectedItem := model.Layout.ChatMenu.SelectedItem().(ChatCmd)
					model.Layout.ChatInput.SetValue(selectedItem.title)
					return selectedItem.Execute(model)
				}
			} else {
				if model.Layout.ChatInput.Value() != "" && model.Layout.ChatInput.Focused() {
					input := model.Layout.ChatInput.Value()

					userMsg := ChatMsg{
						Content:   input,
						Role:      "User",
						Timestamp: time.Now().Unix(),
					}

					model.Layout.Msgs = append(model.Layout.Msgs, userMsg)

					aiMsg := ChatMsg{
						Content:   "",
						Role:      "AI",
						Timestamp: time.Now().Unix(),
					}

					model.Layout.Msgs = append(model.Layout.Msgs, aiMsg)

					model.Layout.ChatView.SetContent(model.Layout.RenderMsgs())
					model.Layout.ChatView.GotoBottom()
					model.Layout.ChatInput.Reset()
					model.Layout.IsLoading = true

					go func() {
						err := model.Provider.Run(context.TODO(), input, model.Layout.Session)
						if err != nil {
							log.Printf("Error: %v", err)
						}
					}()

					return model, model.Layout.Loader.Tick
				}
			}
		}

	case StreamMsg:
		if model.Layout.IsLoading {
			model.Layout.IsLoading = false
		}

		model.Layout.Msgs[len(model.Layout.Msgs)-1].Content += msg.Chunk
		model.Layout.Msgs[len(model.Layout.Msgs)-1].Timestamp = time.Now().Unix()

		model.Layout.ChatView.SetContent(model.Layout.RenderMsgs())
		model.Layout.ChatView.GotoBottom()

		return model, model.HandleStream
	}

	cmds := []tea.Cmd{}
	l, cmd := model.Layout.Update(msg)
	model.Layout = l.(LayoutView)
	model.Provider = model.Layout.Provider

	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model ChatModel) HandleStream() tea.Msg {
	content := <-model.Stream
	return StreamMsg{
		Chunk: content,
	}
}
