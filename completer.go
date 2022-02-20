package prompt

import tea "github.com/charmbracelet/bubbletea"

type completerState int

const (
	idle completerState = iota
	running
)

type Completer func(input string) Suggestions

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
	return c.updateCompletions("")
}

func (c completerModel) Update(msg tea.Msg) (completerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case completionMsg:
		c.state = idle
		c.suggestions = Suggestions(msg)
	}
	return c, nil
}

func (c *completerModel) updateCompletions(text string) tea.Cmd {
	// If completer is already running or the text input hasn't changed, don't run the completer again
	if c.state == running || text == c.prevText {
		return nil
	}

	c.state = running
	c.prevText = text

	return func() tea.Msg {
		filtered := c.completerFunc(text)
		return completionMsg(filtered)
	}
}
