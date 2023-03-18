package dropdown

import (
	"math"

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
	formatters         suggestion.Formatters
	err                error
}

func NewDropdownSuggestionModel[T any](textInput input.Input[T], options ...Option[T]) *Model[T] {
	defaultMaxSuggestions := 6
	m := &Model[T]{
		textInput:          textInput,
		state:              idle,
		maxSuggestions:     defaultMaxSuggestions,
		selectionIndicator: "",
		scrollbar:          " ",
		scrollbarThumb:     " ",
		formatters:         suggestion.DefaultFormatters(),
		prevRunes: []rune(
			" ",
		), // Need to set the previous text to something in order to force the initial render
	}
	for _, option := range options {
		option(m)
	}

	return m
}

func (m *Model[T]) Init() tea.Cmd {
	// Since the user hasn't typed anything on init, call the completer with empty text
	return m.ResetSuggestions()
}

func (m *Model[T]) Update(msg tea.Msg) tea.Cmd {
	switch msg := msg.(type) {
	case suggestion.SuggestionMsg[T]:
		if m.ignoreCount > 0 {
			// Request was in progress when resetSuggestions was called, don't update suggestions
			m.ignoreCount--
		} else {
			m.state = idle
			if msg.Suggestions == nil {
				m.suggestions = []suggestion.Suggestion[T]{}
			} else {
				m.suggestions = msg.Suggestions
			}

			m.err = msg.Err
			// Selection is out of range of the current view or the key is no longer present
			if m.scrollPosition > len(m.suggestions)-1 || m.SelectedSuggestion() == nil {
				m.UnselectSuggestion()
			}

			if m.queueNext {
				// Start another update if it was requested while running
				m.queueNext = false
				return m.UpdateSuggestions()
			}
		}
	case suggestion.PeriodicCompleterMsg:
		if !m.canUpdateSuggestions() {
			return suggestion.PeriodicCompleter(msg.NextTrigger)
		}
		return tea.Batch(m.forceUpdateSuggestions(), suggestion.PeriodicCompleter(msg.NextTrigger))
	case suggestion.OneShotCompleterMsg:
		if !m.canUpdateSuggestions() {
			return nil
		}
		return m.forceUpdateSuggestions()
	case tea.KeyMsg:
		m.lastKeyMsg = msg
		switch msg.Type {
		case tea.KeyTab:
			// Tab suggestion may have changed text so reset previous value
			m.prevRunes = []rune("")
			m.NextSuggestion()
			m.updateIfUnselected()
		case tea.KeyUp:
			m.PreviousSuggestion()
			m.updateIfUnselected()
		case tea.KeyDown:
			m.NextSuggestion()
			m.updateIfUnselected()
		}
	}
	return nil
}

func (m *Model[T]) updateIfUnselected() tea.Cmd {
	if m.IsSuggestionSelected() {
		// Set the input to the suggestion's selected text
		return nil
	} else {
		// Need to update suggestions since we changed the text and the cursor position
		return m.UpdateSuggestions()
	}
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
	scrollbarPos := float64(
		c.scrollPosition,
	) * (float64(windowHeight-scrollbarHeight) / float64(lastSegmentStart))

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

func (m *Model[T]) UpdateSuggestions() tea.Cmd {
	return m.updateSuggestionsCmd(false)
}

func (m *Model[T]) forceUpdateSuggestions() tea.Cmd {
	return m.updateSuggestionsCmd(true)
}

func (m *Model[T]) updateSuggestionsCmd(forceUpdate bool) tea.Cmd {
	runes := m.textInput.Runes()
	cursorPos := m.textInput.CursorIndex()

	runesBeforeCursor := runes
	if cursorPos < len(runes) {
		runesBeforeCursor = runes[:cursorPos]
	}

	// No need to queue another update if the text hasn't changed
	// Don't trim whitespace here because cursor location affects suggestions
	if !forceUpdate && string(runesBeforeCursor) == string(m.prevRunes) {
		return nil
	}

	// Text has changed, but the completer is already running
	// Run again once the current iteration has finished
	if m.state == running {
		m.queueNext = true
		return nil
	}

	m.state = running
	m.prevRunes = runesBeforeCursor

	return suggestion.Complete
}

func (m *Model[T]) ResetSuggestions() tea.Cmd {
	if m.state == running {
		// If suggestion is currently running, ignore the next value and trigger another update
		// This helps speed up getting the next valid result for slow completers
		m.ignoreCount++
	}

	m.state = running
	m.prevRunes = []rune("")

	return suggestion.Complete
}

func (m *Model[T]) UnselectSuggestion() {
	m.selectedKey = nil
	m.scrollPosition = 0
	m.prevScroll = 0
	m.textInput.OnSuggestionUnselected()
}

func (m *Model[T]) ClearSuggestions() {
	m.UnselectSuggestion()
	m.suggestions = []suggestion.Suggestion[T]{}
}

func (m *Model[T]) SelectSuggestion(suggestion suggestion.Suggestion[T]) {
	m.selectedKey = suggestion.Key()
	m.textInput.OnSuggestionChanged(suggestion)
}

func (m *Model[T]) IsSuggestionSelected() bool {
	return m.selectedKey != nil
}

func (m *Model[T]) NextSuggestion() {
	if len(m.suggestions) == 0 {
		return
	}
	index := m.SelectedIndex()
	if index < len(m.suggestions)-1 {
		m.prevScroll = m.scrollPosition
		m.SelectSuggestion(m.suggestions[index+1])
		if index+1 >= m.scrollPosition+m.maxSuggestions {
			m.scrollPosition++
		}

	} else {
		m.UnselectSuggestion()
	}
}

func (m *Model[T]) PreviousSuggestion() {
	if len(m.suggestions) == 0 {
		return
	}

	index := m.SelectedIndex()
	if index > 0 {
		m.prevScroll = m.scrollPosition
		m.SelectSuggestion(m.suggestions[index-1])
		if index-1 < m.scrollPosition {
			m.scrollPosition--
		}
	} else {
		m.UnselectSuggestion()
	}
}

func (m *Model[T]) SelectedIndex() int {
	if m.IsSuggestionSelected() {
		for i, suggestion := range m.suggestions {
			if *suggestion.Key() == *m.selectedKey {
				return i
			}
		}
	}
	return -1
}

func (m *Model[T]) SelectedSuggestion() *suggestion.Suggestion[T] {
	if m.IsSuggestionSelected() {
		for _, suggestion := range m.suggestions {
			if *suggestion.Key() == *m.selectedKey {
				return &suggestion
			}
		}
	}
	return nil
}

func (c Model[T]) MaxSuggestionWidth() (int, int) {
	suggestions := c.Suggestions()

	maxNameLen := 0
	maxDescLen := 0

	// Determine longest name and description to calculate padding
	for _, cur := range suggestions {
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

	return maxNameLen, maxDescLen
}

func (c Model[T]) Render(paddingSize int) string {
	if c.Error() != nil {
		return c.formatters.ErrorText.Render(c.Error().Error())
	}

	suggestions := c.Suggestions()
	if len(suggestions) == 0 {
		return ""
	}

	maxNameLen, maxDescLen := c.MaxSuggestionWidth()

	numSuggestions := len(c.Suggestions())

	visibleSuggestions := c.VisibleSuggestions()
	scrollbarStart, scrollbarEnd := c.ScrollbarBounds()

	prompts := []string{}
	listPosition := c.SelectedIndex() - c.ScrollPosition()
	scrollbar := c.formatters.Scrollbar.Render(c.Scrollbar())
	scrollbarThumb := c.formatters.ScrollbarThumb.Render(c.ScrollbarThumb())
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

		line := cur.Render(
			selected,
			maxNameLen,
			maxDescLen,
			c.formatters,
			scrollbarView,
			c.SelectionIndicator(),
		)
		prompts = append(prompts, line)
	}
	hasBorder := c.formatters.Suggestions.GetBorderLeft()

	allPrompts := lipgloss.JoinVertical(lipgloss.Left, prompts...)

	if hasBorder {
		borderPadding := 2
		return c.formatters.Suggestions.
			Copy().
			MarginLeft(paddingSize - borderPadding).
			PaddingLeft(1).
			Render(allPrompts)
	} else {
		return c.formatters.Suggestions.
			Copy().
			PaddingLeft(paddingSize).
			Render(allPrompts)
	}
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

func (m *Model[T]) Formatters() suggestion.Formatters {
	return m.formatters
}

func (m *Model[T]) SetFormatters(formatters suggestion.Formatters) {
	m.formatters = formatters
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

func (m *Model[T]) ShouldChangeListPosition(msg tea.Msg) bool {
	if key, ok := msg.(tea.KeyMsg); ok {
		switch key.Type {
		case tea.KeyUp, tea.KeyDown, tea.KeyTab:
			return true
		}
	}

	return false
}
