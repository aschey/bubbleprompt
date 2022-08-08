package executor

import (
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AsyncStringModel struct {
	outputFunc func() (string, error)
	output     *string
	spinner    spinner.Model
}

type outputMsg string

func NewAsyncStringModel(outputFunc func() (string, error)) AsyncStringModel {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
	return AsyncStringModel{outputFunc: outputFunc, spinner: spin}
}

func (m AsyncStringModel) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, func() tea.Msg {
		output, err := m.outputFunc()
		if err != nil {
			return ErrorMsg(err)
		}
		return outputMsg(output)
	})
}

func (m AsyncStringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	m.spinner, cmd = m.spinner.Update(msg)
	if msg, ok := msg.(outputMsg); ok {
		stringMsg := string(msg)
		m.output = &stringMsg
		return m, tea.Quit
	}
	return m, cmd
}

func (m AsyncStringModel) View() string {
	if m.output == nil {
		return m.spinner.View() + "Loading..."
	}
	return *m.output
}
