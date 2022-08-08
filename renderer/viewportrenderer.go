package renderer

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewportRenderer struct {
	viewport viewport.Model
	history  string
}

func NewViewportRenderer() *ViewportRenderer {
	return &ViewportRenderer{}
}

func (v *ViewportRenderer) View() string {
	return v.viewport.View()
}

func (v *ViewportRenderer) Initialize(msg tea.WindowSizeMsg) {
	v.viewport = viewport.New(msg.Width, msg.Height-1)
	v.viewport.KeyMap.Up = key.NewBinding(key.WithKeys("ctrl+up"))
	v.viewport.KeyMap.Down = key.NewBinding(key.WithKeys("ctrl+down"))
}

func (v *ViewportRenderer) SetSize(msg tea.WindowSizeMsg) {
	v.viewport.Width = msg.Width
	v.viewport.Height = msg.Height - 1
}

func (v *ViewportRenderer) Update(msg tea.Msg) (Renderer, tea.Cmd) {
	viewport, cmd := v.viewport.Update(msg)
	v.viewport = viewport

	return v, cmd
}

func (v *ViewportRenderer) SetContent(content string) {
	v.viewport.SetContent(v.history + content)
}

func (v *ViewportRenderer) AddOutput(output string) {
	v.history += (output + "\n")
}

func (v *ViewportRenderer) FinishUpdate() tea.Cmd {
	return nil
}

func (v *ViewportRenderer) GotoBottom(msg tea.Msg) {
	keyMsg, isKeyMsg := msg.(tea.KeyMsg)
	if !isKeyMsg || (keyMsg.String() != "ctrl+up" && keyMsg.String() != "ctrl+down") {
		v.viewport.GotoBottom()
	}
}
