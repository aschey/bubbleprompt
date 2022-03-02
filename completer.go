package prompt

import (
	"strings"

	"github.com/aschey/bubbleprompt/commandinput"
	tea "github.com/charmbracelet/bubbletea"
)

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
	queueNext     bool
	ignore        bool
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
	return c.resetCompletions()
}

func (c completerModel) Update(msg tea.Msg, input commandinput.Model) (completerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case completionMsg:
		c.state = idle
		if c.ignore {
			// Request was in progress when resetCompletions was called, don't update suggestions
			c.ignore = false
		} else {
			c.suggestions = Suggestions(msg)
		}

		if c.queueNext {
			// Start another update if it was requested while running
			c.queueNext = false
			return c, c.updateCompletions(input)
		}
	}
	return c, nil
}

func (c *completerModel) updateCompletionsInternal(text string, cursorPos int) tea.Cmd {
	// If completer is already running or the text input hasn't changed, don't run the completer again
	textTrimmed := strings.TrimSpace(text)
	textBeforeCursor := text
	if cursorPos < len(textTrimmed) {
		textBeforeCursor = text[:cursorPos]
	}

	// No need to queue another update if the text hasn't changed
	if textBeforeCursor == c.prevText {
		return nil
	}

	// Text has changed, but the completer is already running
	// Run again once the current iteration has finished
	if c.state == running {
		c.queueNext = true
		return nil
	}

	c.state = running
	c.prevText = textBeforeCursor

	return func() tea.Msg {
		filtered := c.completerFunc(Document{
			Input:             text,
			InputBeforeCursor: textBeforeCursor,
			CursorPosition:    cursorPos,
		})
		return completionMsg(filtered)
	}
}

func (c *completerModel) updateCompletions(input commandinput.Model) tea.Cmd {
	// If completer is already running or the text input hasn't changed, don't run the completer again
	text := input.Value()
	cursorPos := input.Cursor()
	return c.updateCompletionsInternal(text, cursorPos)
}

func (c *completerModel) resetCompletions() tea.Cmd {
	if c.state == running {
		// If completion is currently running, ignore the next value and trigger another update
		// This helps speed up getting the next valid result for slow completers
		c.state = idle
		c.ignore = true
	}

	return c.updateCompletionsInternal("", 0)
}
