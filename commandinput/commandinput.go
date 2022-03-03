package commandinput

import (
	"encoding/csv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Arg struct {
	Text      string
	Style     lipgloss.Style
	Formatter func(arg string) string
}

type Model struct {
	textinput        textinput.Model
	Placeholder      string
	Prompt           string
	Args             []Arg
	TextStyle        lipgloss.Style
	CursorStyle      lipgloss.Style
	PlaceholderStyle lipgloss.Style
	//DefaultArgStyle  lipgloss.Style
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
	textModel := m.textinput

	value := m.Value()

	pos := m.Cursor()
	words := strings.SplitN(value[:pos], " ", 2)
	wordsFull := strings.SplitN(value, " ", 2)
	v := m.TextStyle.Render(words[0])
	if len(words) > 1 {
		v += " " + words[1]
	}

	argLen := len(m.Args)
	numWords := 0
	argStart := 0
	startPadding := ""
	if argLen > 0 {
		r := csv.NewReader(strings.NewReader(value))
		r.Comma = ' '
		r.LazyQuotes = true
		record, _ := r.Read()
		for _, w := range record {
			if len(w) > 0 {
				numWords++
			}
		}
		argStart = numWords - 1
		if argStart < 0 {
			argStart = 0
		} else if argStart > argLen {
			argStart = argLen
		}

		if !strings.HasSuffix(value, " ") {
			startPadding = " "
		}
	}

	if pos < len(wordsFull[0]) {
		v += m.renderWithCursor(wordsFull[0], pos, m.TextStyle)
		if len(wordsFull) > 1 {
			v += " " + wordsFull[1]
		}

		if strings.HasPrefix(m.Placeholder, value) {
			v += m.PlaceholderStyle.Render(m.Placeholder[len(value):])
		}
	} else if pos < len(value) {
		v += m.renderWithCursor(value, pos, lipgloss.NewStyle())
		if strings.HasPrefix(m.Placeholder, value) {
			v += m.PlaceholderStyle.Render(m.Placeholder[len(value):])
		}
	} else if pos < len(m.Placeholder) && strings.HasPrefix(m.Placeholder, value) {
		v += m.renderWithCursor(m.Placeholder, pos, m.PlaceholderStyle)
	} else if argStart == argLen || (numWords > argLen && value[len(value)-1] == ' ') {
		v += m.cursorView(" ", m.TextStyle)
	}

	if argLen > 0 && argStart < argLen {
		if pos == len(value) && (!strings.HasPrefix(m.Placeholder, value) || pos == len(m.Placeholder)) {
			v += m.renderWithCursor(startPadding+m.Args[argStart].Text, 0, m.Args[argStart].Style)

		} else {
			v += m.Args[argStart].Style.Render(startPadding + m.Args[argStart].Text)
		}
		if argStart < argLen {
			for _, arg := range m.Args[argStart+1:] {
				v += " " + arg.Style.Render(arg.Text)
			}
		}
	}

	return textModel.PromptStyle.Render(m.Prompt) + v
}

func (m Model) renderWithCursor(text string, pos int, s lipgloss.Style) string {
	v := m.cursorView(string(text[pos]), s)
	v += s.Render(text[pos+1:])
	return v
}

func (m Model) cursorView(v string, s lipgloss.Style) string {
	if m.textinput.Blink() {
		return s.Render(v)
	}
	return m.CursorStyle.Inline(true).Reverse(true).Render(v)
}
