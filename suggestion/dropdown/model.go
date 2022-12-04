package dropdown

import (
	"math"
	"strings"

	"github.com/aschey/bubbleprompt/formatter"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/suggestion"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
)

type completerState int

const (
	idle completerState = iota
	running
)

type Model[T any] struct {
	textInput          input.Input[T]
	state              completerState
	suggestions        []suggestion.Suggestion[T]
	lastKeyMsg         tea.KeyMsg
	maxSuggestions     int
	scrollPosition     int
	prevScroll         int
	selectedKey        *string
	prevRunes          []rune
	queueNext          bool
	ignoreCount        int
	selectionIndicator string
	scrollbar          string
	scrollbarThumb     string
	err                error
}

func NewDropdownSuggestionModel[T any](textInput input.Input[T]) *Model[T] {
	defaultMaxSuggestions := 6
	return &Model[T]{
		textInput:          textInput,
		state:              idle,
		maxSuggestions:     defaultMaxSuggestions,
		selectionIndicator: "",
		scrollbar:          " ",
		scrollbarThumb:     " ",
		prevRunes:          []rune(" "), // Need to set the previous text to something in order to force the initial render
	}
}

func (c *Model[T]) Init() tea.Cmd {
	// Since the user hasn't typed anything on init, call the completer with empty text
	return c.ResetSuggestions()
}

func (c *Model[T]) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case suggestion.SuggestionMsg[T]:
		if c.ignoreCount > 0 {
			// Request was in progress when resetSuggestions was called, don't update suggestions
			c.ignoreCount--
		} else {
			c.state = idle
			if msg.Suggestions == nil {
				c.suggestions = []suggestion.Suggestion[T]{}
			} else {
				c.suggestions = msg.Suggestions
			}

			c.err = msg.Err
			// Selection is out of range of the current view or the key is no longer present
			if c.scrollPosition > len(c.suggestions)-1 || c.SelectedSuggestion() == nil {
				c.UnselectSuggestion()
			}

			if c.queueNext {
				// Start another update if it was requested while running
				c.queueNext = false
				return c.UpdateSuggestions()
			}
		}
	case suggestion.PeriodicCompleterMsg:
		if !c.canUpdateSuggestions() {
			return suggestion.PeriodicCompleter(msg.NextTrigger)
		}
		return tea.Batch(c.forceUpdateSuggestions(), suggestion.PeriodicCompleter(msg.NextTrigger))
	case suggestion.OneShotCompleterMsg:
		if !c.canUpdateSuggestions() {
			return nil
		}
		return c.forceUpdateSuggestions()
	case tea.KeyMsg:
		c.lastKeyMsg = msg
		if msg.Type == tea.KeyTab {
			// Tab suggestion may have changed text so reset previous value
			c.prevRunes = []rune("")
		}
	}
	return nil
}

func (c Model[T]) canUpdateSuggestions() bool {
	runes := c.textInput.Runes()
	if len(c.textInput.SuggestionRunes(runes[:c.textInput.CursorIndex()])) == 0 {
		return true
	}

	switch c.lastKeyMsg.Type {
	case tea.KeyUp, tea.KeyDown, tea.KeyTab:
		return false
	default:
		return true
	}
}

func (c Model[T]) ScrollbarBounds() (int, int) {
	windowHeight := c.windowHeight()
	contentHeight := len(c.suggestions)
	// The zero-based index of the first element that will be shown when the content is scrolled to the bottom
	lastSegmentStart := contentHeight - windowHeight
	scrollbarHeight := int(math.Max(float64(windowHeight-lastSegmentStart), 1))
	scrollbarPos := float64(c.scrollPosition) * (float64(windowHeight-scrollbarHeight) / float64(lastSegmentStart))

	// If scrolling up, use ceiling operation to ensure the scrollbar is only at the top when the first row is shown
	// otherwise use floor operation
	var scrollbarTop int
	if c.prevScroll > c.scrollPosition {
		scrollbarTop = int(math.Ceil(scrollbarPos))
	} else {
		scrollbarTop = int(math.Floor(scrollbarPos))
	}

	return scrollbarTop, scrollbarTop + scrollbarHeight
}

func (c *Model[T]) UpdateSuggestions() tea.Cmd {
	return c.updateSuggestionsCmd(false)
}

func (c *Model[T]) forceUpdateSuggestions() tea.Cmd {
	return c.updateSuggestionsCmd(true)
}

func (c *Model[T]) updateSuggestionsCmd(forceUpdate bool) tea.Cmd {
	runes := c.textInput.Runes()
	cursorPos := c.textInput.CursorIndex()

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

	return suggestion.Complete
}

func (c *Model[T]) ResetSuggestions() tea.Cmd {
	if c.state == running {
		// If suggestion is currently running, ignore the next value and trigger another update
		// This helps speed up getting the next valid result for slow completers
		c.ignoreCount++
	}

	c.state = running
	c.prevRunes = []rune("")

	return suggestion.Complete
}

func (c *Model[T]) UnselectSuggestion() {
	c.selectedKey = nil
	c.scrollPosition = 0
	c.prevScroll = 0
	c.textInput.OnSuggestionUnselected()
}

func (c *Model[T]) ClearSuggestions() {
	c.UnselectSuggestion()
	c.suggestions = []suggestion.Suggestion[T]{}
}

func (c *Model[T]) SelectSuggestion(suggestion suggestion.Suggestion[T]) {
	c.selectedKey = suggestion.Key()
	c.textInput.OnSuggestionChanged(suggestion)
}

func (c *Model[T]) IsSuggestionSelected() bool {
	return c.selectedKey != nil
}

func (c *Model[T]) NextSuggestion() {
	if len(c.suggestions) == 0 {
		return
	}
	index := c.SelectedIndex()
	if index < len(c.suggestions)-1 {
		c.prevScroll = c.scrollPosition
		c.SelectSuggestion(c.suggestions[index+1])
		if index+1 >= c.scrollPosition+c.maxSuggestions {
			c.scrollPosition++
		}

	} else {
		c.UnselectSuggestion()
	}
}

func (c *Model[T]) PreviousSuggestion() {
	if len(c.suggestions) == 0 {
		return
	}

	index := c.SelectedIndex()
	if index > 0 {
		c.prevScroll = c.scrollPosition
		c.SelectSuggestion(c.suggestions[index-1])
		if index-1 < c.scrollPosition {
			c.scrollPosition--
		}
	} else {
		c.UnselectSuggestion()
	}
}

func (c *Model[T]) SelectedIndex() int {
	if c.IsSuggestionSelected() {
		for i, suggestion := range c.suggestions {
			if *suggestion.Key() == *c.selectedKey {
				return i
			}
		}
	}
	return -1
}

func (c *Model[T]) SelectedSuggestion() *suggestion.Suggestion[T] {
	if c.IsSuggestionSelected() {
		for _, suggestion := range c.suggestions {
			if *suggestion.Key() == *c.selectedKey {
				return &suggestion
			}
		}
	}
	return nil
}

func (c Model[T]) Render(paddingSize int, formatters formatter.Formatters) string {
	if c.Error() != nil {
		return formatters.ErrorText.Render(c.Error().Error())
	}
	maxNameLen := 0
	maxDescLen := 0

	// Determine longest name and description to calculate padding
	for _, cur := range c.Suggestions() {
		suggestionText := cur.GetSuggestionText()
		textWidth := runewidth.StringWidth(suggestionText)
		if textWidth > maxNameLen {
			maxNameLen = textWidth
		}

		descWidth := runewidth.StringWidth(cur.Description)
		if descWidth > maxDescLen {
			maxDescLen = descWidth
		}
	}
	numSuggestions := len(c.Suggestions())

	visibleSuggestions := c.VisibleSuggestions()
	scrollbarStart, scrollbarEnd := c.ScrollbarBounds()

	// Add left offset
	leftPadding := lipgloss.
		NewStyle().
		PaddingLeft(paddingSize).
		Render("")

	prompts := []string{}
	listPosition := c.SelectedIndex() - c.ScrollPosition()
	scrollbar := formatters.Scrollbar.Render(c.Scrollbar())
	scrollbarThumb := formatters.ScrollbarThumb.Render(c.ScrollbarThumb())
	for i, cur := range visibleSuggestions {
		selected := i == listPosition
		scrollbarView := ""
		if numSuggestions > c.MaxSuggestions() {
			if scrollbarStart <= i && i < scrollbarEnd {
				scrollbarView = scrollbarThumb
			} else {
				scrollbarView = scrollbar
			}
		}

		line := cur.Render(selected, leftPadding, maxNameLen, maxDescLen, formatters, scrollbarView, c.SelectionIndicator())
		prompts = append(prompts, line)
	}

	return strings.Join(prompts, "\n")
}

func (m *Model[T]) EnableScrollbar() {
	m.scrollbar = " "
	m.scrollbarThumb = " "
}

func (m *Model[T]) DisableScrollbar() {
	m.scrollbar = ""
	m.scrollbarThumb = ""
}

func (m *Model[T]) MaxSuggestions() int {
	return m.maxSuggestions
}

func (m *Model[T]) SetMaxSuggestions(maxSuggestions int) {
	m.maxSuggestions = maxSuggestions
}

func (m *Model[T]) SelectionIndicator() string {
	return m.selectionIndicator
}

func (m *Model[T]) SetSelectionIndicator(selectionIndicator string) {
	m.selectionIndicator = selectionIndicator
}

func (m *Model[T]) Suggestions() []suggestion.Suggestion[T] {
	return m.suggestions
}

func (m *Model[T]) windowHeight() int {
	windowHeight := len(m.Suggestions())
	if windowHeight > m.MaxSuggestions() {
		windowHeight = m.MaxSuggestions()
	}
	return windowHeight
}

func (m *Model[T]) VisibleSuggestions() []suggestion.Suggestion[T] {
	windowHeight := m.windowHeight()
	visibleSuggestions := m.Suggestions()[m.scrollPosition : m.scrollPosition+windowHeight]
	return visibleSuggestions
}

func (m *Model[T]) Error() error {
	return m.err
}

func (m *Model[T]) ScrollPosition() int {
	return m.scrollPosition
}

func (m *Model[T]) Scrollbar() string {
	return m.scrollbar
}

func (m *Model[T]) ScrollbarThumb() string {
	return m.scrollbarThumb
}
