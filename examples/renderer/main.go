package main

import (
	"fmt"
	"os"

	prompt "github.com/aschey/bubbleprompt"
	"github.com/aschey/bubbleprompt/completer"
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/parserinput"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/renderer"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type appModel struct {
	suggestions []input.Suggestion[any]
	textInput   *simpleinput.Model[any]
}

func (m appModel) Complete(promptModel prompt.Model[any]) ([]input.Suggestion[any], error) {
	if len(m.textInput.AllTokens()) > 1 {
		return nil, nil
	}

	return completer.FilterHasPrefix(m.textInput.CurrentTokenBeforeCursor(), m.suggestions), nil
}

func (m appModel) Execute(input string, promptModel *prompt.Model[any]) (tea.Model, error) {
	tokens := m.textInput.TokenValues()
	if len(tokens) == 0 {
		return executor.NewStringModel("No selection"), nil
	}
	switch tokens[0] {
	case "viewport":
		return executor.NewCmdModel("set viewport renderer", prompt.SetRenderer(renderer.NewViewportRenderer(0), true)), nil
	case "unmanaged":
		return executor.NewCmdModel("set unmanaged renderer", prompt.SetRenderer(renderer.NewUnmanagedRenderer(), true)), nil
	}
	return executor.NewStringModel("You selected " + tokens[0]), nil
}

func (m appModel) Update(msg tea.Msg) (prompt.AppModel[any], tea.Cmd) {
	return m, nil
}

func main() {
	textInput := simpleinput.New(simpleinput.WithLexerOptions(parserinput.WithCursorMode[any](textinput.CursorStatic)))
	suggestions := []input.Suggestion[any]{
		{Text: "unmanaged", Description: "use the unmanaged renderer"},
		{Text: "viewport", Description: "use the viewport renderer"},
	}

	appModel := appModel{
		suggestions: suggestions,
		textInput:   textInput,
	}

	promptModel, err := prompt.New[any](
		appModel,
		textInput,
	)
	if err != nil {
		panic(err)
	}

	if _, err := tea.NewProgram(promptModel, tea.WithFilter(prompt.MsgFilter)).Run(); err != nil {
		fmt.Printf("Could not start program :(\n%v\n", err)
		os.Exit(1)
	}
}
