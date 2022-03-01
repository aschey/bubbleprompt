package commandinput

import (
	"encoding/csv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	textinput        textinput.Model
	Placeholder      string
	Prompt           string
	Args             []string
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
	argStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	textModel := m.textinput
	styleText := m.TextStyle.Render

	value := m.Value()

	pos := m.Cursor()
	v := styleText(value[:pos])

	argLen := len(m.Args)
	numWords := 0
	argView := ""
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
		argStart := numWords - 1
		if argStart < 0 {
			argStart = 0
		} else if argStart > argLen {
			argStart = argLen
		}
		argView = strings.Join(m.Args[argStart:], " ")
		if !strings.HasSuffix(value, " ") {
			argView = " " + argView
		}
	}

	if pos < len(value) {
		v += m.renderWithCursor(value, pos, m.TextStyle)
		if strings.HasPrefix(m.Placeholder, value) {
			v += m.PlaceholderStyle.Render(m.Placeholder[len(value):])
		}
	} else if pos < len(m.Placeholder) && strings.HasPrefix(m.Placeholder, value) {
		v += m.renderWithCursor(m.Placeholder, pos, m.PlaceholderStyle)
	} else if argLen == 0 || (numWords > argLen && value[len(value)-1] == ' ') {
		v += m.cursorView(" ", m.TextStyle)
	}

	if len(argView) > 0 && pos == len(value) && (!strings.HasPrefix(m.Placeholder, value) || pos == len(m.Placeholder)) {
		v += m.renderWithCursor(argView, 0, argStyle)
	} else {
		v += argStyle.Render(argView)
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
