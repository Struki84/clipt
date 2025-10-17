package tui

import (
	"context"
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type StreamMsg struct {
	Chunk string
}

type ChatModel struct {
	Layout LayoutView

	Providers []ChatProvider
	Provider  ChatProvider
	Session   Session

	Stream chan string
}

func NewChatModel(providers []ChatProvider) ChatModel {
	provider := providers[0]
	return ChatModel{
		Providers: providers,
		Provider:  provider,
		Layout:    NewLayoutView(provider),
		Stream:    make(chan string),
	}
}

func (model ChatModel) Init() tea.Cmd {
	model.Provider.Stream(context.TODO(), func(ctx context.Context, chunk []byte) error {
		model.Stream <- string(chunk)
		return nil
	})

	cmds := []tea.Cmd{}
	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, model.HandleStream)

	return tea.Batch(cmds...)
}

func (model ChatModel) View() string {
	return model.Layout.View()
}

func (model ChatModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return model, tea.Quit
		case tea.KeyEnter:
			if !model.Layout.MenuActive {
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

					go func() {
						err := model.Provider.Run(context.TODO(), input)
						if err != nil {
							log.Printf("Error: %v", err)
						}
					}()

					return model, nil
				}
			}
		}

	case StreamMsg:
		model.Layout.Msgs[len(model.Layout.Msgs)-1].Content += msg.Chunk
		model.Layout.Msgs[len(model.Layout.Msgs)-1].Timestamp = time.Now().Unix()

		model.Layout.ChatView.SetContent(model.Layout.RenderMsgs())
		model.Layout.ChatView.GotoBottom()

		return model, model.HandleStream
	}

	cmds := []tea.Cmd{}
	l, cmd := model.Layout.Update(msg)
	model.Layout = l.(LayoutView)
	cmds = append(cmds, cmd)

	return model, tea.Batch(cmds...)
}

func (model ChatModel) HandleStream() tea.Msg {
	content := <-model.Stream
	return StreamMsg{
		Chunk: content,
	}
}
