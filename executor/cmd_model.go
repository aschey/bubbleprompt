package executor

import (
	"github.com/aschey/bubbleprompt/internal"
	tea "github.com/charmbracelet/bubbletea"
)

type CmdModel struct {
	cmd    tea.Cmd
	output string
}

func NewCmdModel(output string, cmd tea.Cmd) CmdModel {
	return CmdModel{cmd, output}
}

func (m CmdModel) Init() tea.Cmd {
	return m.cmd
}

func (m CmdModel) Update(tea.Msg) (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

func (m CmdModel) View() string {
	return internal.AddNewlineIfMissing(m.output)
}
