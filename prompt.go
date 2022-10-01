package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: https://no-color.org/

type errMsg error

type modelState int

const (
	completing modelState = iota
	executing
)

type Executor[T any] func(input string, selectedSuggestion *input.Suggestion[T]) (tea.Model, error)

const DefaultScrollbarColor = "13"
const DefaultScrollbarThumbColor = "14"

type Model[I any] struct {
	completer               completerModel[I]
	executor                Executor[I]
	textInput               input.Input[I]
	renderer                renderer.Renderer
	Formatters              input.Formatters
	executorModel           *executorModel
	modelState              modelState
	scrollbar               string
	scrollbarThumb          string
	lastTypedCursorPosition int
	typedText               string
	ready                   bool
	err                     error
}

func New[I any](completer Completer[I], executor Executor[I], textInput input.Input[I], opts ...Option[I]) (Model[I], error) {
	formatters := input.DefaultFormatters()
	model := Model[I]{
		completer:  newCompleterModel(completer, textInput, formatters.ErrorText, 6),
		executor:   executor,
		textInput:  textInput,
		renderer:   &renderer.UnmanagedRenderer{},
		Formatters: formatters,
	}
	model.SetScrollbarColor(lipgloss.Color(DefaultScrollbarColor))
	model.SetScrollbarThumbColor(lipgloss.Color(DefaultScrollbarThumbColor))
	for _, opt := range opts {
		if err := opt(&model); err != nil {
			return Model[I]{}, err
		}
	}

	return model, nil
}

func (m *Model[I]) SetScrollbarColor(color lipgloss.TerminalColor) {
	m.scrollbar = lipgloss.NewStyle().Background(color).Render(" ")
}

func (m *Model[I]) SetScrollbarThumbColor(color lipgloss.TerminalColor) {
	m.scrollbarThumb = lipgloss.NewStyle().Background(color).Render(" ")
}

func (m *Model[I]) SetMaxSuggestions(maxSuggestions int) {
	m.completer.maxSuggestions = maxSuggestions
}

func (m *Model[I]) SetRenderer(renderer renderer.Renderer) {
	m.renderer = renderer
}

var shutdown bool = false

func OnQuit() tea.QuitBehavior {
	if shutdown {
		return tea.Shutdown
	} else {
		return tea.PreventShutdown
	}
}

func (m Model[I]) Init() tea.Cmd {
	return tea.Batch(m.textInput.Init(), m.completer.Init())
}

func (m Model[I]) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.renderer.View()
}
