package schema

import "github.com/charmbracelet/lipgloss"

type Styles struct {
	ChatHeader lipgloss.Style
	ChatInput  lipgloss.Style
	InfoLine   lipgloss.Style
	ChatMenu   struct {
		View          lipgloss.Style
		TitleNormal   lipgloss.Style
		TitleSelected lipgloss.Style
		Description   lipgloss.Style
	}

	StatusLine struct {
		BaseStyle    lipgloss.Style
		Tab          lipgloss.Style
		Mode         lipgloss.Style
		ProviderType lipgloss.Style
		ProviderName lipgloss.Style
		Loader       lipgloss.Style
	}

	Msg struct {
		User     lipgloss.Style
		AI       lipgloss.Style
		Sys      lipgloss.Style
		Err      lipgloss.Style
		Internal lipgloss.Style
	}
}
