package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
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

type Executor[T any] func(input string) (tea.Model, error)

const DefaultScrollbarColor = "13"
const DefaultScrollbarThumbColor = "14"

type Model[I any] struct {
	completer               completerModel[I]
	executor                Executor[I]
	textInput               input.Input[I]
	viewport                viewport.Model
	Formatters              input.Formatters
	previousCommands        string
	executorModel           *executorModel
	modelState              modelState
	scrollbar               string
	scrollbarThumb          string
	lastTypedCursorPosition int
	typedText               string
	ready                   bool
	err                     error
}

func New[I any](completer Completer[I], executor Executor[I], textInput input.Input[I], opts ...Option[I]) Model[I] {
	model := Model[I]{
		completer:  newCompleterModel(completer, textInput, 6),
		executor:   executor,
		textInput:  textInput,
		Formatters: input.DefaultFormatters(),
	}
	model.SetScrollbarColor(lipgloss.Color(DefaultScrollbarColor))
	model.SetScrollbarThumbColor(lipgloss.Color(DefaultScrollbarThumbColor))
	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
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

func (m Model[I]) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.completer.Init())
}

func (m Model[I]) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}
