package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type modelState int

const (
	completing modelState = iota
	executing
)

type Executor func(input string, selected *input.Suggestion, suggestions []input.Suggestion) (tea.Model, error)

type Model struct {
	completer               completerModel
	executor                Executor
	textInput               input.Input
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

func New(completer Completer, executor Executor, textInput input.Input, opts ...Option) Model {
	model := Model{
		completer:  newCompleterModel(completer, textInput, 6),
		executor:   executor,
		textInput:  textInput,
		Formatters: input.DefaultFormatters(),
	}
	model.SetScrollbarColor(lipgloss.Color("14"))
	model.SetScrollbarThumbColor(lipgloss.Color("240"))
	for _, opt := range opts {
		if err := opt(&model); err != nil {
			panic(err)
		}
	}

	return model
}

func (m *Model) SetScrollbarColor(color lipgloss.TerminalColor) {
	m.scrollbar = lipgloss.NewStyle().Background(color).Render(" ")
}

func (m *Model) SetScrollbarThumbColor(color lipgloss.TerminalColor) {
	m.scrollbarThumb = lipgloss.NewStyle().Background(color).Render(" ")
}

func (m *Model) SetMaxSuggestions(maxSuggestions int) {
	m.completer.maxSuggestions = maxSuggestions
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(textinput.Blink, m.completer.Init())
}

func (m Model) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	return m.viewport.View()
}
