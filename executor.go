package prompt

import (
	"github.com/aschey/bubbleprompt/executor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ExecutorFinishedMsg tea.Model

type executorModel struct {
	inner          tea.Model
	errorTextStyle lipgloss.Style
	err            error
}

func newExecutorModel(inner tea.Model, errorTextStyle lipgloss.Style, err error) *executorModel {
	return &executorModel{
		inner:          inner,
		errorTextStyle: errorTextStyle,
		err:            err,
	}
}

func (m executorModel) Init() tea.Cmd {
	return m.inner.Init()
}

func (m executorModel) Update(msg tea.Msg) (executorModel, tea.Cmd) {
	inner, cmd := m.inner.Update(msg)
	m.inner = inner
	if msg, ok := msg.(executor.ErrorMsg); ok {
		m.err = error(msg)
		return m, tea.Quit
	} else {
		return m, cmd
	}

}

func (m executorModel) View() string {
	if m.err != nil {
		return m.errorTextStyle.Render(m.err.Error()) + "\n"
	} else {
		return m.inner.View()
	}

}
