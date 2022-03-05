package commandinput

import (
	"encoding/csv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Arg struct {
	Text             string
	PlaceholderStyle lipgloss.Style
	ArgStyle         lipgloss.Style
	Formatter        func(arg string) string
}

type Model struct {
	textinput        textinput.Model
	Placeholder      string
	Prompt           string
	Args             []Arg
	PromptStyle      lipgloss.Style
	TextStyle        lipgloss.Style
	CursorStyle      lipgloss.Style
	PlaceholderStyle lipgloss.Style
}

func New() Model {
	textinput := textinput.New()
	return Model{
		textinput:        textinput,
		Placeholder:      "",
		Prompt:           "> ",
		PlaceholderStyle: textinput.PlaceholderStyle,
	}
}

func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		cursorPos := m.Cursor()
		if m.shouldSkipUpdate(msg, cursorPos) || m.shouldSkipUpdate(msg, cursorPos-1) {
			return m, nil
		}
	}
	m.textinput, cmd = m.textinput.Update(msg)
	return m, cmd
}

func (m *Model) Focus() tea.Cmd {
	return m.textinput.Focus()
}

func (m Model) Value() string {
	return m.textinput.Value()
}

func (m *Model) SetValue(s string) {
	m.textinput.SetValue(s)
}

func (m Model) Cursor() int {
	return m.textinput.Cursor()
}

func (m *Model) SetCursor(pos int) {
	m.textinput.SetCursor(pos)
}

func (m Model) Focused() bool {
	return m.textinput.Focused()
}

func (m *Model) Blur() {
	m.textinput.Blur()
}

func (m Model) View() string {
	view := m.getViewBeforeCursor()
	view += m.getPlaceholder()
	argPlaceholders := m.getArgPlaceholders()
	if argPlaceholders == "" && m.Cursor() == len(m.Value()) {
		view += m.cursorView(" ", m.TextStyle)
	}
	view += argPlaceholders

	return m.PromptStyle.Render(m.Prompt) + view
}

func (m Model) formatArgs(text string) string {
	words := strings.Split(text, " ")
	view := ""
	for i, arg := range words {
		view += " "
		if i < len(m.Args) {
			view += m.Args[i].ArgStyle.Render(arg)
		} else {
			view += arg
		}
	}

	return view
}

func (m Model) shouldSkipUpdate(msg tea.KeyMsg, pos int) bool {
	text := m.Value()
	// Don't allow consecutive spaces because this interferes with rendering arguments
	return pos < len(text) && msg.String() == " " && text[pos] == ' '
}

func (m Model) getViewBeforeCursor() string {
	words := strings.SplitN(m.Value()[:m.Cursor()], " ", 2)
	view := m.TextStyle.Render(words[0])
	if len(words) > 1 {
		view += m.formatArgs(words[1])
	}

	return view
}

func (m Model) getPlaceholder() string {
	view := ""
	cursorPos := m.Cursor()
	value := m.Value()
	allText := strings.SplitN(value, " ", 2)
	command := allText[0]

	if cursorPos < len(command) {
		args := ""
		if len(allText) > 1 {
			args = m.formatArgs(allText[1])
		}
		view += m.renderWithPlaceholder(command, args, m.TextStyle)
	} else if cursorPos < len(value) {
		cursorPos := m.Cursor()
		before := strings.Split(value[:cursorPos], " ")
		isInWord := value[cursorPos] != ' '
		isMiddleOfWord := isInWord && value[cursorPos-1] != ' '
		wordsBeforeCursor := []string{}
		for _, w := range before {
			if len(w) > 0 {
				wordsBeforeCursor = append(wordsBeforeCursor, w)
			}
		}
		skipArgs := len(wordsBeforeCursor) - 1
		if isMiddleOfWord {
			skipArgs--
		}
		after := strings.Split(value[cursorPos:], " ")
		wordsAfterCursor := []string{}
		for _, w := range after {
			if len(w) > 0 {
				wordsAfterCursor = append(wordsAfterCursor, w)
			}
		}

		for i, arg := range wordsAfterCursor {
			if i > 0 || !isInWord {
				if cursorPos == len(m.Value()[:cursorPos]+view) {
					view += m.cursorView(" ", lipgloss.NewStyle())
				} else {
					view += " "
				}
			}

			idx := i + skipArgs
			style := lipgloss.NewStyle()
			if idx >= 0 && idx < len(m.Args) {
				style = m.Args[idx].ArgStyle
			}

			lenBefore := len(m.Value()[:cursorPos] + view)
			if cursorPos >= lenBefore && cursorPos < lenBefore+len(arg) {
				view += m.renderWithCursor(arg, cursorPos-lenBefore, style)
			} else {
				view += style.Render(arg)
			}
		}

	} else if cursorPos < len(m.Placeholder) && strings.HasPrefix(m.Placeholder, value) {
		view += m.renderWithCursor(m.Placeholder, cursorPos, m.PlaceholderStyle)
	}

	return view
}

func (m Model) renderWithPlaceholder(text string, args string, style lipgloss.Style) string {
	value := m.Value()
	view := m.renderWithCursor(text, m.Cursor(), style) + args
	if strings.HasPrefix(m.Placeholder, value) {
		view += m.PlaceholderStyle.Render(m.Placeholder[len(value):])
	}

	return view
}

func (m Model) getArgPlaceholders() string {
	argLen := len(m.Args)
	placeholderStart := m.getPlaceholderStart()

	if placeholderStart < 0 {
		placeholderStart = 0
	} else if placeholderStart > argLen {
		placeholderStart = argLen
	}

	if placeholderStart >= argLen {
		return ""
	}

	startPadding := ""
	value := m.Value()
	if !strings.HasSuffix(value, " ") {
		startPadding = " "
	}

	argView := ""
	cursorPos := m.Cursor()
	if cursorPos == len(value) && (!strings.HasPrefix(m.Placeholder, value) || cursorPos == len(m.Placeholder)) {
		argView += m.renderWithCursor(startPadding+m.Args[placeholderStart].Text, 0, m.Args[placeholderStart].PlaceholderStyle)
	} else {
		argView += m.Args[placeholderStart].PlaceholderStyle.Render(startPadding + m.Args[placeholderStart].Text)
	}

	for _, arg := range m.Args[placeholderStart+1:] {
		argView += " " + arg.PlaceholderStyle.Render(arg.Text)
	}

	return argView
}

func (m Model) getPlaceholderStart() int {
	numWords := 0
	reader := csv.NewReader(strings.NewReader(m.Value()))
	reader.Comma = ' '
	reader.LazyQuotes = true
	record, _ := reader.Read()
	for _, w := range record {
		if len(w) > 0 {
			numWords++
		}
	}

	return numWords - 1
}

func (m Model) renderWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	v := m.cursorView(string(text[cursorPos]), s)
	v += s.Render(text[cursorPos+1:])
	return v
}

func (m Model) cursorView(v string, s lipgloss.Style) string {
	if m.textinput.Blink() {
		return s.Render(v)
	}
	return m.CursorStyle.Inline(true).Reverse(true).Render(v)
}
