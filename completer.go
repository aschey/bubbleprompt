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
	selectedKey   *string
	prevText      string
	queueNext     bool
	ignoreCount   int
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
		if c.ignoreCount > 0 {
			// Request was in progress when resetCompletions was called, don't update suggestions
			c.ignoreCount--
		} else {
			c.state = idle
			c.suggestions = Suggestions(msg)
			if c.getSelectedSuggestion() == nil {
				c.unselectSuggestion()
			}
			if c.queueNext {
				// Start another update if it was requested while running
				c.queueNext = false
				return c, c.updateCompletions(input)
			}
		}
	}
	return c, nil
}

func (c *completerModel) updateCompletions(input commandinput.Model) tea.Cmd {
	text := strings.Split(input.Value()[:input.Cursor()], " ")[0]
	cursorPos := input.Cursor()

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

func (c *completerModel) resetCompletions() tea.Cmd {
	if c.state == running {
		// If completion is currently running, ignore the next value and trigger another update
		// This helps speed up getting the next valid result for slow completers
		c.ignoreCount++
	}

	c.state = running
	c.prevText = ""

	return func() tea.Msg {
		filtered := c.completerFunc(Document{
			Input:             "",
			InputBeforeCursor: "",
			CursorPosition:    0,
		})
		return completionMsg(filtered)
	}
}

func (c *completerModel) unselectSuggestion() {
	c.selectedKey = nil
}

func (c *completerModel) isSuggestionSelected() bool {
	return c.selectedKey != nil
}

func (c *completerModel) nextSuggestion() {
	if len(c.suggestions) == 0 {
		return
	}
	index := c.getSelectedIndex()
	if index < len(c.suggestions)-1 {
		c.selectedKey = c.suggestions[index+1].key()
	} else {
		c.unselectSuggestion()
	}
}

func (c *completerModel) previousSuggestion() {
	if len(c.suggestions) == 0 {
		return
	}

	index := c.getSelectedIndex()
	if index > 0 {
		c.selectedKey = c.suggestions[index-1].key()
	} else {
		c.unselectSuggestion()
	}
}

func (c *completerModel) getSelectedIndex() int {
	if c.isSuggestionSelected() {
		for i, suggestion := range c.suggestions {
			if *suggestion.key() == *c.selectedKey {
				return i
			}
		}
	}
	return -1
}

func (c *completerModel) getSelectedSuggestion() *Suggestion {
	if c.isSuggestionSelected() {
		for _, suggestion := range c.suggestions {
			if *suggestion.key() == *c.selectedKey {
				return &suggestion
			}
		}
	}
	return nil
}
