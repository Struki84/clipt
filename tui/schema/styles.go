package schema

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	SessionBar lipgloss.Style
	ChatInput  lipgloss.Style
	ChatMenu   lipgloss.Style

	StatusLine struct {
		BaseStyle    lipgloss.Style
		Tab          lipgloss.Style
		Mode         lipgloss.Style
		ProviderType lipgloss.Style
		ProviderName lipgloss.Style
	}

	Msg struct {
		User     lipgloss.Style
		AI       lipgloss.Style
		Sys      lipgloss.Style
		Err      lipgloss.Style
		Internal lipgloss.Style
	}
}
