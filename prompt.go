package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: https://no-color.org/

type errMsg error

type modelState int

const (
	completing modelState = iota
	executing
)

type Executor[T any] func(input string, selectedSuggestion *input.Suggestion[T]) (tea.Model, error)

type Model[T any] struct {
	completer               completerModel[T]
	executor                Executor[T]
	textInput               input.Input[T]
	renderer                renderer.Renderer
	Formatters              input.Formatters
	executorModel           *executorModel
	modelState              modelState
	lastTypedCursorPosition int
	typedText               string
	ready                   bool
	err                     error
}

func New[T any, I input.Input[T]](completer Completer[T], executor Executor[T], textInput I, opts ...Option[T]) (Model[T], error) {
	formatters := input.DefaultFormatters()
	model := Model[T]{
		completer:  newCompleterModel(completer, textInput, formatters.ErrorText, 6),
		executor:   executor,
		textInput:  textInput,
		renderer:   &renderer.UnmanagedRenderer{},
		Formatters: formatters,
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			return Model[T]{}, err
		}
	}

	return model, nil
}

func (m *Model[T]) SetMaxSuggestions(maxSuggestions int) {
	m.completer.maxSuggestions = maxSuggestions
}

func (m *Model[T]) SetRenderer(renderer renderer.Renderer) {
	m.renderer = renderer
}

var shutdown bool = false

type quitAttempted struct{}

func MsgFilter(_ tea.Model, msg tea.Msg) tea.Msg {
	if _, ok := msg.(tea.QuitMsg); ok && !shutdown {
		return quitAttempted{}
	}

	return msg
}

func (m Model[T]) Init() tea.Cmd {
	return tea.Batch(m.textInput.Init(), m.completer.Init())
}

func (m Model[T]) View() string {
	if !m.ready {
		return "\n  Initializing..."
	}
	view := m.renderer.View()
	if shutdown {
		// For the final view, add one more newline so the terminal prompt doesn't cut off the last line
		view += "\n"
	}
	return view
}
