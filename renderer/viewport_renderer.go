package renderer

import (
	"github.com/aschey/bubbleprompt/internal"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type ViewportRenderer struct {
	viewport viewport.Model
	history  string
	settings viewportSettings
}

func NewViewportRenderer(options ...ViewportOption) *ViewportRenderer {
	settings := viewportSettings{
		widthOffset:  0,
		heightOffset: 0,
		useHistory:   true,
	}
	for _, option := range options {
		option(&settings)
	}
	return &ViewportRenderer{settings: settings}
}

func (v *ViewportRenderer) View() string {
	return v.viewport.View()
}

func (v *ViewportRenderer) Initialize(msg tea.WindowSizeMsg) {
	v.SetSize(msg)
	v.viewport.KeyMap.Up = key.NewBinding(key.WithKeys("ctrl+up"))
	v.viewport.KeyMap.Down = key.NewBinding(key.WithKeys("ctrl+down"))
}

func (v *ViewportRenderer) SetSize(msg tea.WindowSizeMsg) {
	v.viewport.Width = msg.Width - v.settings.widthOffset
	v.viewport.Height = msg.Height - v.settings.heightOffset
}

func (v *ViewportRenderer) Update(msg tea.Msg) (Renderer, tea.Cmd) {
	viewport, cmd := v.viewport.Update(msg)
	v.viewport = viewport

	return v, cmd
}

func (v *ViewportRenderer) SetContent(content string) {
	v.viewport.SetContent(v.history + content)
}

func (v *ViewportRenderer) GetHistory() string {
	return v.history
}

func (v *ViewportRenderer) SetHistory(history string) tea.Cmd {
	if v.settings.useHistory {
		v.history = internal.AddNewlineIfMissing(history)
	}
	return nil
}

func (v *ViewportRenderer) AddOutput(output string) {
	if v.settings.useHistory {
		v.history += internal.AddNewlineIfMissing(output)
	}
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
