package prompt

import (
	"math"
	"strings"
	"time"

	"github.com/aschey/bubbleprompt/input"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type completerState int

const (
	idle completerState = iota
	running
)

type Completer[T any] func(prompt Model[T]) ([]input.Suggestion[T], error)

type completionMsg[T any] struct {
	suggestions []input.Suggestion[T]
	err         error
}

type PeriodicCompleterMsg struct {
	NextTrigger time.Duration
}

type OneShotCompleterMsg struct{}

func PeriodicCompleter(nextTrigger time.Duration) tea.Cmd {
	return tea.Tick(nextTrigger, func(_ time.Time) tea.Msg {
		return PeriodicCompleterMsg{NextTrigger: nextTrigger}
	})
}

func OneShotCompleter(nextTrigger time.Duration) tea.Cmd {
	return tea.Tick(nextTrigger, func(_ time.Time) tea.Msg {
		return OneShotCompleterMsg{}
	})
}

type completerModel[T any] struct {
	completerFunc  Completer[T]
	state          completerState
	textInput      input.Input[T]
	suggestions    []input.Suggestion[T]
	errorText      lipgloss.Style
	lastKeyMsg     tea.KeyMsg
	maxSuggestions int
	scroll         int
	prevScroll     int
	selectedKey    *string
	prevRunes      []rune
	queueNext      bool
	ignoreCount    int
	err            error
}

func newCompleterModel[T any](completerFunc Completer[T], textInput input.Input[T], errorText lipgloss.Style, maxSuggestions int) completerModel[T] {
	return completerModel[T]{
		textInput:      textInput,
		completerFunc:  completerFunc,
		state:          idle,
		maxSuggestions: maxSuggestions,
		errorText:      errorText,
		prevRunes:      []rune(" "), // Need to set the previous text to something in order to force the initial render
	}
}

func (c completerModel[T]) Init() tea.Cmd {
	// Since the user hasn't typed anything on init, call the completer with empty text
	return c.resetCompletions(Model[T]{})
}

func (c completerModel[T]) Update(msg tea.Msg, prompt Model[T]) (completerModel[T], tea.Cmd) {
	switch msg := msg.(type) {
	case completionMsg[T]:
		if c.ignoreCount > 0 {
			// Request was in progress when resetCompletions was called, don't update suggestions
			c.ignoreCount--
		} else {
			c.state = idle
			if msg.suggestions == nil {
				c.suggestions = []input.Suggestion[T]{}
			} else {
				c.suggestions = msg.suggestions
			}

			c.err = msg.err
			// Selection is out of range of the current view or the key is no longer present
			if c.scroll > len(c.suggestions)-1 || c.getSelectedSuggestion() == nil {
				c.unselectSuggestion()
			}

			if c.queueNext {
				// Start another update if it was requested while running
				c.queueNext = false
				return c, c.updateCompletions(prompt)
			}
		}
	case PeriodicCompleterMsg:
		if !c.canUpdateCompletions() {
			return c, PeriodicCompleter(msg.NextTrigger)
		}
		return c, tea.Batch(c.forceUpdateCompletions(prompt), PeriodicCompleter(msg.NextTrigger))
	case OneShotCompleterMsg:
		if !c.canUpdateCompletions() {
			return c, nil
		}
		return c, c.forceUpdateCompletions(prompt)
	case tea.KeyMsg:
		c.lastKeyMsg = msg
		if msg.Type == tea.KeyTab {
			// Tab completion may have changed text so reset previous value
			c.prevRunes = []rune("")
		}
	}
	return c, nil
}

func (c completerModel[T]) canUpdateCompletions() bool {
	runes := c.textInput.Runes()
	if len(c.textInput.CompletionRunes(runes[:c.textInput.CursorIndex()])) == 0 {
		return true
	}

	switch c.lastKeyMsg.Type {
	case tea.KeyUp, tea.KeyDown, tea.KeyTab:
		return false
	default:
		return true
	}
}

func (c completerModel[T]) scrollbarBounds(windowHeight int) (int, int) {
	contentHeight := len(c.suggestions)
	// The zero-based index of the first element that will be shown when the content is scrolled to the bottom
	lastSegmentStart := contentHeight - windowHeight
	scrollbarHeight := int(math.Max(float64(windowHeight-lastSegmentStart), 1))
	scrollbarPos := float64(c.scroll) * (float64(windowHeight-scrollbarHeight) / float64(lastSegmentStart))

	// If scrolling up, use ceiling operation to ensure the scrollbar is only at the top when the first row is shown
	// otherwise use floor operation
	var scrollbarTop int
	if c.prevScroll > c.scroll {
		scrollbarTop = int(math.Ceil(scrollbarPos))
	} else {
		scrollbarTop = int(math.Floor(scrollbarPos))
	}

	return scrollbarTop, scrollbarTop + scrollbarHeight
}

func (c completerModel[T]) Render(paddingSize int, formatters input.Formatters) string {
	if c.err != nil {
		return c.errorText.Render(c.err.Error())
	}
	maxNameLen := 0
	maxDescLen := 0

	// Determine longest name and description to calculate padding
	for _, cur := range c.suggestions {
		completionText := cur.GetCompletionText()
		textWidth := runewidth.StringWidth(completionText)
		if textWidth > maxNameLen {
			maxNameLen = textWidth
		}

		descWidth := runewidth.StringWidth(cur.Description)
		if descWidth > maxDescLen {
			maxDescLen = descWidth
		}
	}
	numSuggestions := len(c.suggestions)
	windowHeight := numSuggestions
	if windowHeight > c.maxSuggestions {
		windowHeight = c.maxSuggestions
	}
	visibleSuggestions := c.suggestions[c.scroll : c.scroll+windowHeight]
	scrollbarStart, scrollbarEnd := c.scrollbarBounds(windowHeight)

	// Add left offset
	leftPadding := lipgloss.
		NewStyle().
		PaddingLeft(paddingSize).
		Render("")

	prompts := []string{}
	listPosition := c.getSelectedIndex() - c.scroll
	scrollbar := formatters.Scrollbar.Render(" ")
	scrollbarThumb := formatters.ScrollbarThumb.Render(" ")
	for i, cur := range visibleSuggestions {
		selected := i == listPosition
		scrollbarView := ""
		if numSuggestions > c.maxSuggestions {
			if scrollbarStart <= i && i < scrollbarEnd {
				scrollbarView = scrollbarThumb
			} else {
				scrollbarView = scrollbar
			}
		}

		line := cur.Render(selected, leftPadding, maxNameLen, maxDescLen, formatters, scrollbarView)
		prompts = append(prompts, line)
	}

	return strings.Join(prompts, "\n")
}

func (c *completerModel[T]) updateCompletions(prompt Model[T]) tea.Cmd {
	return c.updateCompletionsCmd(prompt, false)
}

func (c *completerModel[T]) forceUpdateCompletions(prompt Model[T]) tea.Cmd {
	return c.updateCompletionsCmd(prompt, true)
}

func (c *completerModel[T]) updateCompletionsCmd(prompt Model[T], forceUpdate bool) tea.Cmd {
	input := prompt.textInput
	runes := input.Runes()
	cursorPos := input.CursorIndex()

	runesBeforeCursor := runes
	if cursorPos < len(runes) {
		runesBeforeCursor = runes[:cursorPos]
	}

	// No need to queue another update if the text hasn't changed
	// Don't trim whitespace here because cursor location affects suggestions
	if !forceUpdate && string(runesBeforeCursor) == string(c.prevRunes) {
		return nil
	}

	// Text has changed, but the completer is already running
	// Run again once the current iteration has finished
	if c.state == running {
		c.queueNext = true
		return nil
	}

	c.state = running
	c.prevRunes = runesBeforeCursor

	return func() tea.Msg {
		filtered, err := c.completerFunc(prompt)
		return completionMsg[T]{suggestions: filtered, err: err}
	}
}

func (c *completerModel[T]) resetCompletions(prompt Model[T]) tea.Cmd {
	if c.state == running {
		// If completion is currently running, ignore the next value and trigger another update
		// This helps speed up getting the next valid result for slow completers
		c.ignoreCount++
	}

	c.state = running
	c.prevRunes = []rune("")

	return func() tea.Msg {
		filtered, err := c.completerFunc(prompt)
		return completionMsg[T]{suggestions: filtered, err: err}
	}
}

func (c *completerModel[T]) unselectSuggestion() {
	c.selectedKey = nil
	c.scroll = 0
	c.prevScroll = 0
	c.textInput.OnSuggestionUnselected()
}

func (c *completerModel[T]) clearSuggestions() {
	c.unselectSuggestion()
	c.suggestions = []input.Suggestion[T]{}
}

func (c *completerModel[T]) selectSuggestion(suggestion input.Suggestion[T]) {
	c.selectedKey = suggestion.Key()
	c.textInput.OnSuggestionChanged(suggestion)
}

func (c *completerModel[T]) isSuggestionSelected() bool {
	return c.selectedKey != nil
}

func (c *completerModel[T]) nextSuggestion() {
	if len(c.suggestions) == 0 {
		return
	}
	index := c.getSelectedIndex()
	if index < len(c.suggestions)-1 {
		c.prevScroll = c.scroll
		c.selectSuggestion(c.suggestions[index+1])
		if index+1 >= c.scroll+c.maxSuggestions {
			c.scroll++
		}

	} else {
		c.unselectSuggestion()
	}
}

func (c *completerModel[T]) previousSuggestion() {
	if len(c.suggestions) == 0 {
		return
	}

	index := c.getSelectedIndex()
	if index > 0 {
		c.prevScroll = c.scroll
		c.selectSuggestion(c.suggestions[index-1])
		if index-1 < c.scroll {
			c.scroll--
		}
	} else {
		c.unselectSuggestion()
	}
}

func (c *completerModel[T]) getSelectedIndex() int {
	if c.isSuggestionSelected() {
		for i, suggestion := range c.suggestions {
			if *suggestion.Key() == *c.selectedKey {
				return i
			}
		}
	}
	return -1
}

func (c *completerModel[T]) getSelectedSuggestion() *input.Suggestion[T] {
	if c.isSuggestionSelected() {
		for _, suggestion := range c.suggestions {
			if *suggestion.Key() == *c.selectedKey {
				return &suggestion
			}
		}
	}
	return nil
}
