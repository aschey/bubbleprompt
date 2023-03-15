package searchbar

import (
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model[T any] struct {
	promptModel  prompt.Model[T]
	contentModel tea.Model
	searchText   string
	searchBar    string
}

func New[T any](inputHandler prompt.InputHandler[T], textInput input.Input[T],
	contentModel tea.Model, options ...prompt.Option[T],
) Model[T] {
	formatters := suggestion.DefaultFormatters().Minimal()
	formatters.Suggestions = formatters.Suggestions.Border(lipgloss.RoundedBorder(), false, true, true)
	promptModel := prompt.New(inputHandler, textInput,
		append(options,
			prompt.WithUnmanagedRenderer[T](renderer.WithUseHistory(false)),
			prompt.WithFormatters[T](formatters))...)

	searchbarWidth := 50
	searchText := "Search:"
	searchBar := lipgloss.NewStyle().PaddingRight(searchbarWidth).Border(lipgloss.RoundedBorder()).Render(searchText)

	return Model[T]{
		promptModel:  promptModel,
		contentModel: contentModel,
		searchBar:    searchBar,
		searchText:   searchText,
	}
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
	promptRenderer := m.promptModel.Renderer().(*renderer.UnmanagedRenderer)
	suggestionManager := m.promptModel.SuggestionManager()
	vis := suggestionManager.VisibleSuggestions()
	maxNameLen := 0
	for _, vis := range vis {
		if len(vis.GetSuggestionText()) > maxNameLen {
			maxNameLen = len(vis.GetSuggestionText())
		}
	}
	offset := m.promptModel.SuggestionOffset() - spacing
	if offset < 0 {
		offset = 0
	}
	topView := strings.Repeat(" ", offset) + "╮" + strings.Repeat(" ", maxNameLen+spacing) + "╭"
	promptView := lipgloss.JoinVertical(lipgloss.Left, "  "+promptRenderer.Input(), topView, promptRenderer.Body())
	return renderer.PlaceOverlay(len(m.searchText), 1, promptView, view)
}
