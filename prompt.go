package prompt

import (
	"github.com/aschey/bubbleprompt/editor"
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

type InputHandler[T any] interface {
	Update(msg tea.Msg) (InputHandler[T], tea.Cmd)
	Execute(input string, prompt *Model[T]) (tea.Model, error)
	Complete(prompt Model[T]) ([]editor.Suggestion[T], error)
}

type Model[T any] struct {
	completionManager       completionManager[T]
	inputHandler            InputHandler[T]
	textInput               editor.Editor[T]
	renderer                renderer.Renderer
	formatters              editor.Formatters
	executionManager        *executionManager
	modelState              modelState
	lastTypedCursorPosition int
	typedRunes              []rune
	ready                   bool
	size                    tea.WindowSizeMsg
	err                     error
}

func New[T any](inputHandler InputHandler[T], textInput editor.Editor[T], opts ...Option[T]) (Model[T], error) {
	formatters := editor.DefaultFormatters()
	model := Model[T]{
		completionManager: newCompletionManager(textInput, formatters.ErrorText, 6),
		inputHandler:      inputHandler,
		textInput:         textInput,
		renderer:          &renderer.UnmanagedRenderer{},
		formatters:        formatters,
	}

	for _, opt := range opts {
		if err := opt(&model); err != nil {
			return Model[T]{}, err
		}
	}

	return model, nil
}

func (m *Model[T]) SetMaxSuggestions(maxSuggestions int) {
	m.completionManager.maxSuggestions = maxSuggestions
}

func (m Model[T]) Formatters() editor.Formatters {
	return m.formatters
}

func (m *Model[T]) SetFormatters(formatters editor.Formatters) {
	m.formatters = formatters
}

func (m Model[T]) SelectedSuggestion() *editor.Suggestion[T] {
	return m.completionManager.getSelectedSuggestion()
}

func (m Model[T]) TextInput() editor.Editor[T] {
	return m.textInput
}

type rendererMsg struct {
	renderer      renderer.Renderer
	retainHistory bool
}

func SetRenderer(r renderer.Renderer, retainHistory bool) tea.Cmd {
	return func() tea.Msg {
		return rendererMsg{
			renderer:      r,
			retainHistory: retainHistory,
		}
	}
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
	return tea.Batch(m.textInput.Init(), m.completionManager.Init(m))
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
