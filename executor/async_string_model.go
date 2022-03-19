package executor

import tea "github.com/charmbracelet/bubbletea"

type AsyncStringModel struct {
	outputFunc func() string
	output     *string
}

type outputMsg string

func NewAsyncStringModel(outputFunc func() string) AsyncStringModel {
	return AsyncStringModel{outputFunc: outputFunc}
}

func (m AsyncStringModel) Init() tea.Cmd {
	return func() tea.Msg {
		return outputMsg(m.outputFunc())
	}
}

func (m AsyncStringModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case outputMsg:
		stringMsg := string(msg)
		m.output = &stringMsg
		return m, tea.Quit
	}
	return m, nil
}

func (m AsyncStringModel) View() string {
	if m.output == nil {
		return "Loading..."
	}
	return *m.output
}
