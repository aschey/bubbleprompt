package searchbar

import (
	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model[T any] struct {
	promptModel  prompt.Model[T]
	contentModel tea.Model
	searchText   string
	searchBar    string
}

func New[T any](promptModel prompt.Model[T], contentModel tea.Model) Model[T] {
	searchbarWidth := 50
	searchText := "Search:"
	searchBar := lipgloss.NewStyle().PaddingRight(searchbarWidth).Border(lipgloss.RoundedBorder()).Render(searchText)
	return Model[T]{promptModel: promptModel, contentModel: contentModel, searchBar: searchBar, searchText: searchText}
}

func (m Model[T]) Init() tea.Cmd {
	return tea.Batch(m.promptModel.Init(), m.contentModel.Init())
}

func (m Model[T]) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	promptModel, cmd := m.promptModel.Update(msg)
	m.promptModel = promptModel.(prompt.Model[T])
	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		borderSize := 1
		m.contentModel, cmd = m.contentModel.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - (borderSize*2 + 1),
		})
	default:
		m.contentModel, cmd = m.contentModel.Update(msg)
	}
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m Model[T]) View() string {
	view := lipgloss.JoinVertical(lipgloss.Left, m.searchBar, m.contentModel.View())
	spacing := 2
	return renderer.PlaceOverlay(len(m.searchText)+spacing, 1, m.promptModel.View(), view)
}
