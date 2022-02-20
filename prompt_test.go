package prompt

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	prompt Model
}

type testCompleterModel struct {
	suggestions []Suggestion
}

func (m model) Init() tea.Cmd {
	return m.prompt.Init()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	p, cmd := m.prompt.Update(msg)
	m.prompt = p
	return m, cmd
}

func (m model) View() string {
	return m.prompt.View()
}

func (m completerModel) completer(input string) Suggestions {
	//time.Sleep(100 * time.Millisecond)
	return FilterHasPrefix(input, m.suggestions)
}

func executor(input string, selected *Suggestion, suggestions Suggestions) tea.Model {
	return NewAsyncStringModel(func() string {
		//time.Sleep(10 * time.Millisecond)
		return "test"
	})
}

type reader struct {
	iter int
}

func (r *reader) Read(p []byte) (n int, err error) {
	if r.iter == 0 {
		n = copy(p, "test")
		r.iter++
		return n, nil
	}
	//time.Sleep(1000 * time.Millisecond)
	n = copy(p, []byte{3})
	return n, nil
}

type writer struct {
	last string
}

const ansi = "[\u001B\u009B][[\\]()#;?]*(?:(?:(?:[a-zA-Z\\d]*(?:;[a-zA-Z\\d]*)*)?\u0007)|(?:(?:\\d{1,4}(?:;\\d{0,4})*)?[\\dA-PRZcf-ntqry=><~]))"

var re = regexp.MustCompile(ansi)

func (w *writer) Write(p []byte) (n int, err error) {
	last := strings.TrimSpace(re.ReplaceAllString(string(p), ""))
	if len(last) > 0 {
		w.last = last
	}
	return len(p), nil
}

func Test(t *testing.T) {

	suggestions := []Suggestion{
		{Name: "first option", Description: "test desc", Placeholder: "[hh]"},
		{Name: "second option", Description: "test desc2"},
		{Name: "third option", Description: "test desc2"},
		{Name: "fourth option", Description: "test desc2"},
		{Name: "fifth option", Description: "test desc2"},
	}

	completerModel := completerModel{suggestions: suggestions}

	m := model{prompt: New(
		completerModel.completer,
		executor,
	)}
	m.prompt.ready = true
	m.prompt.viewport = viewport.New(80, 30)

	// data := [2]byte{}
	// data[0] = 3
	// data = append(data, '3')
	// in := bytes.NewReader(data[:])

	in := reader{}
	out := writer{}
	if err := tea.NewProgram(m, tea.WithInput(&in), tea.WithOutput(&out)).Start(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
	lines := strings.Split(out.last, "\n")
	if !strings.Contains(lines[1], "first") {
		panic("fail")
	}
	// var buf []byte
	// out.Read(buf)
	// println("out", string(buf))
}
