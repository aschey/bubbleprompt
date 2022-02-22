package prompt

import tea "github.com/charmbracelet/bubbletea"

type completerState int

const (
	idle completerState = iota
	running
)

type Document struct {
	InputBeforeCursor string
	Input             string
	CursorPosition    int
}

type Completer func(document Document) Suggestions

type completionMsg Suggestions

type completerModel struct {
	completerFunc Completer
	state         completerState
	suggestions   Suggestions
	prevText      string
}

func newCompleterModel(completerFunc Completer) completerModel {
	return completerModel{
		completerFunc: completerFunc,
		state:         idle,
		prevText:      " ", // Need to set the previous text to something in order to force the initial render
	}
}

func (c completerModel) Init() tea.Cmd {
	// Since the user hasn't typed anything on init, call the completer with empty text
	return c.updateCompletions(Model{})
}

func (c completerModel) Update(msg tea.Msg) (completerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case completionMsg:
		c.state = idle
		c.suggestions = Suggestions(msg)
	}
	return c, nil
}

func (c *completerModel) updateCompletions(m Model) tea.Cmd {
	// If completer is already running or the text input hasn't changed, don't run the completer again
	textBeforeCursor := m.textInput.Value()[:m.textInput.Cursor()]
	if c.state == running || textBeforeCursor == c.prevText {
		return nil
	}

	c.state = running
	c.prevText = textBeforeCursor

	return func() tea.Msg {
		filtered := c.completerFunc(Document{
			Input:             m.textInput.Value(),
			InputBeforeCursor: textBeforeCursor,
			CursorPosition:    m.textInput.Cursor(),
		})
		return completionMsg(filtered)
	}
}
