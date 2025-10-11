// TODO
// - 2 modes: Chat, Debug
// - add some kind of debug view for reading reasoning steps
// - load chat history
// [x] top bar with session name / timestamp
// [x] add error msgs
// [x] bottom bar with active mode and active engine
// [x] commands menu: models, agents, session, exit

package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type StreamMsg struct {
	Chunk string
}

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type delegate struct{}

func (d delegate) Height() int                             { return 1 }
func (d delegate) Spacing() int                            { return 0 }
func (d delegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d delegate) Render(w io.Writer, m list.Model, index int, item list.Item) {
	var (
		normalStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#FFFFFF")).
				Padding(0, 0, 0, 0).
				Width(120)

		selectedStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("#11111b")).
				Foreground(lipgloss.Color("#b4befe")).
				Padding(0).
				Width(120)
	)
	i, ok := item.(menuItem)
	if !ok {
		return
	}
	style := normalStyle

	if index == m.Index() {
		style = selectedStyle
	}

	fmt.Fprint(w, style.Render(string(i.Title())))
}

type ChatView struct {
	user       *user.User
	provider   ChatProvider
	msgs       []ChatMsg
	streamChan chan string
	viewport   viewport.Model
	textarea   textarea.Model
	renderer   *glamour.TermRenderer
	windowSize tea.WindowSizeMsg
	menuMode   bool
	menuList   list.Model
	menuItems  []list.Item
	menuHeight int
}

func NewChatViewLight(agent ChatProvider) ChatView {
	ta := textarea.New()
	ta.Prompt = ""
	ta.ShowLineNumbers = false
	ta.FocusedStyle.CursorLine = lipgloss.NewStyle()
	ta.FocusedStyle.Base = lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		MarginBottom(1)
	ta.SetHeight(1)
	ta.Focus()

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(120),
	)

	vp := viewport.New(120, 35)
	vp.KeyMap = viewport.KeyMap{
		PageDown: key.NewBinding(key.WithKeys("pgdown", "space")),
		PageUp:   key.NewBinding(key.WithKeys("pgup")),
		Down:     key.NewBinding(key.WithKeys("down")),
		Up:       key.NewBinding(key.WithKeys("up")),
	}

	user, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	menuItems := []list.Item{
		menuItem{title: "/models", desc: "List available models"},
		menuItem{title: "/agents", desc: "List available agents"},
		menuItem{title: "/session", desc: "List session history"},
		menuItem{title: "/exit", desc: "Exit"},
	}

	list := list.New(menuItems, delegate{}, 0, 0)
	list.SetShowTitle(false)
	list.SetShowHelp(false)
	list.SetShowPagination(false)
	list.SetShowFilter(false)
	list.SetShowStatusBar(false)
	list.SetFilteringEnabled(false)

	list.KeyMap.CursorDown = key.NewBinding(key.WithKeys("down"))
	list.KeyMap.CursorUp = key.NewBinding(key.WithKeys("up"))

	return ChatView{
		user:       user,
		provider:   agent,
		streamChan: make(chan string, 100),
		viewport:   vp,
		textarea:   ta,
		renderer:   renderer,
		menuMode:   false,
		menuList:   list,
		menuItems:  menuItems,
	}
}

func (chat ChatView) Init() tea.Cmd {
	chat.provider.Stream(context.Background(), func(ctx context.Context, chunk []byte) error {
		chat.streamChan <- string(chunk)
		return nil
	})

	cmds := []tea.Cmd{}
	cmds = append(cmds, textarea.Blink)
	cmds = append(cmds, chat.handleStream)

	return tea.Batch(cmds...)
}

func (chat ChatView) handleStream() tea.Msg {
	content := <-chat.streamChan

	return StreamMsg{
		Chunk: content,
	}
}

func (chat ChatView) View() string {
	baseStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		Foreground(lipgloss.Color("#ebdbb2"))

	tabStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#7f849c")).
		PaddingRight(1)

	modeStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#b4befe")).
		Foreground(lipgloss.Color("#181825")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderLeft(true).
		BorderRight(false).
		BorderTop(false).
		BorderBottom(false)

	providerTypeSyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#b4befe")).
		Foreground(lipgloss.Color("#181825")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderRight(true)

	providerNameStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#45475a")).
		Foreground(lipgloss.Color("fff")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#45475a")).
		BorderLeft(true)

	providerType := providerTypeSyle.Render("MODEL")
	providerName := providerNameStyle.Render("GPT-4o")
	tab := tabStyle.Render("tab")
	mode := modeStyle.Render("CHAT")

	leftPart := lipgloss.JoinHorizontal(lipgloss.Top, providerType, providerName)
	rightPart := lipgloss.JoinHorizontal(lipgloss.Top, tab, mode)

	fillerWidth := chat.windowSize.Width - lipgloss.Width(leftPart) - lipgloss.Width(rightPart)
	if fillerWidth < 0 {
		fillerWidth = 0
	}
	filler := baseStyle.Width(fillerWidth).Render("")

	statusView := lipgloss.JoinHorizontal(lipgloss.Top, leftPart, filler, rightPart)

	sessionBarStyle := lipgloss.NewStyle().
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#11111b")).
		PaddingLeft(1).
		PaddingRight(1).
		BorderLeft(true).
		BorderRight(true).
		Width(chat.windowSize.Width - 4)

	sessionBar := sessionBarStyle.Render("New Session \n04 Oct 2025 23:34")
	joinVertical := lipgloss.JoinVertical(
		lipgloss.Center,
		sessionBar,
		chat.viewport.View(),
		chat.textarea.View(),
		statusView,
	)

	menuStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#11111b")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#ffffff")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		PaddingLeft(1).
		PaddingRight(1).
		Width(chat.windowSize.Width - 6).
		Height(chat.menuHeight)

	if chat.menuMode {
		chat.viewport.Height = chat.windowSize.Height - chat.textarea.Height() - chat.menuHeight - 6

		chat.menuList.SetSize(chat.textarea.Width()-6, chat.menuHeight)
		menu := menuStyle.Render(chat.menuList.View())

		joinVertical := lipgloss.JoinVertical(
			lipgloss.Center,
			sessionBar,
			chat.viewport.View(),
			menu,
			chat.textarea.View(),
			statusView,
		)

		return joinVertical
	}

	return joinVertical

}

func (chat ChatView) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		chat.windowSize = msg
		chat.viewport.Width = msg.Width
		chat.textarea.SetWidth(msg.Width - 4)
		chat.viewport.Height = msg.Height - chat.textarea.Height() - 6

		chat.viewport.SetContent(chat.renderMessages())

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return chat, tea.Quit

		case tea.KeyEnter:
			if chat.textarea.Focused() && strings.TrimSpace(chat.textarea.Value()) != "" {
				input := chat.textarea.Value()

				usrMsg := ChatMsg{
					Content:   input,
					Role:      "User",
					Timestamp: time.Now().Unix(),
				}

				chat.msgs = append(chat.msgs, usrMsg)

				aiMsg := ChatMsg{
					Content:   "",
					Role:      "AI",
					Timestamp: time.Now().Unix(),
				}

				chat.msgs = append(chat.msgs, aiMsg)

				chat.viewport.SetContent(chat.renderMessages())
				chat.textarea.Reset()
				chat.viewport.GotoBottom()

				go func() {
					err := chat.provider.Run(context.Background(), input)
					if err != nil {
						log.Println("Run error:", err)
					}
				}()

				return chat, nil
			}
		case tea.KeyEsc:
			if chat.menuMode {
				chat.menuMode = false
			}
		}

	case StreamMsg:
		chat.msgs[len(chat.msgs)-1].Content += msg.Chunk
		chat.msgs[len(chat.msgs)-1].Timestamp = time.Now().Unix()

		chat.viewport.SetContent(chat.renderMessages())
		chat.viewport.GotoBottom()

		return chat, chat.handleStream
	}

	cmds := []tea.Cmd{}

	ta, cmd := chat.textarea.Update(msg)
	chat.textarea = ta
	cmds = append(cmds, cmd)

	vp, cmd := chat.viewport.Update(msg)
	chat.viewport = vp
	cmds = append(cmds, cmd)

	input := chat.textarea.Value()
	if strings.HasPrefix(input, "/") {
		chat.menuMode = true
		filterValue := strings.TrimPrefix(input, "/")
		var filtered []list.Item

		for _, item := range chat.menuItems {

			if strings.Contains(strings.ToLower(item.FilterValue()), strings.ToLower(filterValue)) {
				filtered = append(filtered, item)
			}

		}

		chat.menuList.SetItems(filtered)
		chat.menuList.SetHeight(len(filtered))
		chat.menuHeight = len(filtered)
		// if len(filtered) > 0 {
		// 	chat.menuList.Select(0)
		// }

		var cmd tea.Cmd
		chat.menuList, cmd = chat.menuList.Update(msg)
		cmds = append(cmds, cmd)

	} else {
		chat.menuMode = false
		chat.menuList.SetItems(chat.menuItems)
		chat.menuList.SetHeight(len(chat.menuItems))
		chat.menuHeight = (len(chat.menuItems))
	}

	return chat, tea.Batch(cmds...)
}

func (chat ChatView) renderMessages() string {

	var styledMessages []string

	// sysMsgStyle := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("#181825")).
	// 	BorderStyle(lipgloss.ThickBorder()).
	// 	BorderForeground(lipgloss.Color("#181825")).
	// 	BorderLeft(true).
	// 	BorderRight(true).
	// 	BorderTop(false).
	// 	BorderBottom(false).
	// 	Padding(1).
	// 	Margin(1).
	// 	Align(lipgloss.Left).
	// 	Width(chat.viewport.Width - 4)

	// errMsgStyle := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("#181825")).
	// 	BorderStyle(lipgloss.ThickBorder()).
	// 	BorderForeground(lipgloss.Color("#e64553")).
	// 	BorderLeft(true).
	// 	BorderRight(true).
	// 	BorderTop(false).
	// 	BorderBottom(false).
	// 	Padding(1).
	// 	Margin(1).
	// 	Align(lipgloss.Left).
	// 	Width(chat.viewport.Width - 4)

	// internalMsg := lipgloss.NewStyle().
	// 	Background(lipgloss.Color("#181825")).
	// 	BorderStyle(lipgloss.ThickBorder()).
	// 	BorderForeground(lipgloss.Color("#fab387")).
	// 	BorderLeft(true).
	// 	BorderRight(true).
	// 	BorderTop(false).
	// 	BorderBottom(false).
	// 	Padding(1).
	// 	Margin(1).
	// 	Align(lipgloss.Left).
	// 	Width(chat.viewport.Width - 4)

	userStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("#181825")).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color("#b4befe")).
		BorderLeft(true).
		BorderRight(true).
		BorderTop(false).
		BorderBottom(false).
		Padding(1).
		Margin(1).
		Align(lipgloss.Left).
		Width(chat.viewport.Width - 4)

	aiStyle := lipgloss.NewStyle().
		Align(lipgloss.Left).
		Width(chat.viewport.Width - 2)

	for _, msg := range chat.msgs {
		if msg.Role == "User" {
			date := time.Unix(msg.Timestamp, 0).Format("2 Jan 2006 15:04")
			username := chat.user.Username
			fullMsg := fmt.Sprintf("%s\n\n%s (%s) ", msg.Content, username, date)

			styledMessages = append(styledMessages, userStyle.Render(fullMsg))

		} else if msg.Role == "AI" {
			renderedTxt, _ := chat.renderer.Render(msg.Content)

			styledMessages = append(styledMessages, aiStyle.Render(renderedTxt))
		}
	}

	return lipgloss.JoinVertical(lipgloss.Center, styledMessages...)
}

func ShowChatViewLight(agent ChatProvider) {
	file, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	p := tea.NewProgram(
		NewChatViewLight(agent),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
