package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/aschey/bubbleprompt/suggestion"
	"github.com/aschey/bubbleprompt/suggestion/dropdown"
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
	Init() tea.Cmd
	Update(msg tea.Msg) (InputHandler[T], tea.Cmd)
	Execute(input string, prompt *Model[T]) (tea.Model, error)
	Complete(prompt Model[T]) ([]suggestion.Suggestion[T], error)
}

type Model[T any] struct {
	suggestionManager       suggestion.Manager[T]
	inputHandler            InputHandler[T]
	textInput               input.Input[T]
	renderer                renderer.Renderer
	executionManager        *executionManager
	modelState              modelState
	lastTypedCursorPosition int
	typedRunes              []rune
	ready                   bool
	size                    tea.WindowSizeMsg
	sequenceNumber          int
	focusOnStart            bool
	err                     error
}

func New[T any](
	inputHandler InputHandler[T],
	textInput input.Input[T],
	opts ...Option[T],
) Model[T] {
	model := Model[T]{
		suggestionManager: dropdown.New(textInput),
		inputHandler:      inputHandler,
		textInput:         textInput,
		focusOnStart:      true,
		renderer:          renderer.NewUnmanagedRenderer(),
	}

	for _, opt := range opts {
		opt(&model)
	}

	return model
}

func (m *Model[T]) SuggestionManager() suggestion.Manager[T] {
	return m.suggestionManager
}

func (m Model[T]) TextInput() input.Input[T] {
	return m.textInput
}

func (m Model[T]) Renderer() renderer.Renderer {
	return m.renderer
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

type focusMsg bool

func Focus() tea.Cmd {
	return func() tea.Msg { return focusMsg(true) }
}

func Blur() tea.Cmd {
	return func() tea.Msg { return focusMsg(false) }
}

var shutdown bool = false

type quitAttempted struct{}

func MsgFilter(_ tea.Model, msg tea.Msg) tea.Msg {
	if _, ok := msg.(tea.QuitMsg); ok && !shutdown {
		return quitAttempted{}
	}

	return msg
}

func (m Model[T]) getInitFocus() tea.Cmd {
	if m.focusOnStart {
		return Focus()
	} else {
		return Blur()
	}
}

func (m Model[T]) Init() tea.Cmd {
	return tea.Batch(m.suggestionManager.Init(), m.inputHandler.Init(), m.getInitFocus())
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
