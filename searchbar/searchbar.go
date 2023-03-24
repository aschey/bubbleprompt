package searchbar

import (
	"strings"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/ansi"
)

type Model[T any] struct {
	promptModel      prompt.Model[T]
	contentModel     tea.Model
	settings         searchbarSettings[T]
	searchbarHeight  int
	placeholderStart int
	placeholderLine  int
	searchBar        string
	borderWidth      int
}

func New[T any](inputHandler prompt.InputHandler[T], textInput input.Input[T],
	contentModel tea.Model, options ...Option[T],
) Model[T] {
	defaultMaxWith := 50
	settings := searchbarSettings[T]{
		maxWidth:       defaultMaxWith,
		label:          "Search:",
		searchbarStyle: lipgloss.NewStyle(),
		promptOptions:  []prompt.Option[T]{},
	}

	for _, option := range options {
		option(&settings)
	}

	promptModel := prompt.New(inputHandler, textInput,
		append(settings.promptOptions,
			prompt.WithUnmanagedRenderer[T](renderer.WithUseHistory(false)))...)
	borderWidth := 0
	hasBorder := promptModel.SuggestionManager().Formatters().Suggestions.GetBorderLeft()
	if hasBorder {
		borderWidth = 2
		textInput.SetPrompt("  ")
	} else {
		textInput.SetPrompt("")
	}

	searchBar := settings.searchbarStyle.PaddingRight(settings.maxWidth).Render(settings.label)
	searchbarLines := strings.Split(searchBar, "\n")
	searchbarHeight := len(searchbarLines)
	placeholderStart := 0
	placeholderLine := 0
	for i, line := range searchbarLines {
		if textIndex := strings.Index(line, settings.label); textIndex > -1 {
			placeholderLine = i
			placeholderStart = ansi.PrintableRuneWidth(line[:textIndex])
		}
	}
	return Model[T]{
		promptModel:      promptModel,
		contentModel:     contentModel,
		searchBar:        searchBar,
		settings:         settings,
		searchbarHeight:  searchbarHeight,
		placeholderStart: placeholderStart,
		placeholderLine:  placeholderLine,
		borderWidth:      borderWidth,
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
	updateContent := true
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.contentModel, cmd = m.contentModel.Update(tea.WindowSizeMsg{
			Width:  msg.Width,
			Height: msg.Height - m.searchbarHeight,
		})
		cmds = append(cmds, cmd)
		updateContent = false
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			cmds = append(cmds, prompt.Blur())
		}
	}

	if updateContent {
		m.contentModel, cmd = m.contentModel.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model[T]) PromptModel() *prompt.Model[T] {
	return &m.promptModel
}

func (m Model[T]) OverlayX() int {
	return len(m.settings.label) + m.placeholderStart + 1 - m.borderWidth
}

func (m Model[T]) OverlayY() int {
	return m.placeholderLine
}

func (m Model[T]) BaseView() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.searchBar, m.contentModel.View())
}

func (m Model[T]) Input() string {
	promptRenderer := m.promptModel.Renderer().(*renderer.UnmanagedRenderer)
	return promptRenderer.Input()
}

func (m Model[T]) Body() string {
	promptRenderer := m.promptModel.Renderer().(*renderer.UnmanagedRenderer)
	return promptRenderer.Body()
}

func (m Model[T]) OverlayView() string {
	if len(m.promptModel.SuggestionManager().Suggestions()) == 0 {
		return m.promptModel.View()
	}

	return lipgloss.JoinVertical(lipgloss.Left, m.Input(), m.Body())
}

func (m Model[T]) View() string {
	return renderer.PlaceOverlay(m.OverlayX(), m.OverlayY(), m.OverlayView(), m.BaseView())
}
