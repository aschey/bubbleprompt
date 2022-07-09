package prompt

import (
	"github.com/aschey/bubbleprompt/executor"
	"github.com/aschey/bubbleprompt/input"
	tea "github.com/charmbracelet/bubbletea"
)

type executorModel struct {
	inner     tea.Model
	errorText input.Text
	err       error
}

func newExecutorModel(inner tea.Model, errorText input.Text, err error) *executorModel {
	return &executorModel{
		inner:     inner,
		errorText: errorText,
		err:       err,
	}
}

func (m executorModel) Init() tea.Cmd {
	return m.inner.Init()
}

func (m executorModel) Update(msg tea.Msg) (executorModel, tea.Cmd) {
	inner, cmd := m.inner.Update(msg)
	m.inner = inner
	switch msg := msg.(type) {
	case executor.ErrorMsg:
		m.err = error(msg)
		return m, tea.Quit
	default:
		return m, cmd
	}

}

func (m executorModel) View() string {
	if m.err != nil {
		return m.errorText.Format(m.err.Error()) + "\n"
	} else {
		return m.inner.View()
	}

}
