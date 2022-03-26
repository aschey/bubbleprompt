package commandinput

import "github.com/charmbracelet/lipgloss"

type viewBuilder struct {
	view    string
	rawView string
	model   Model
}

func newViewBuilder(model Model) *viewBuilder {
	return &viewBuilder{model: model}
}

func (v viewBuilder) getView() string {
	if v.model.Cursor() == v.viewLen() {
		return v.view + v.cursorView(" ", lipgloss.NewStyle())
	}
	return v.view
}

func (v viewBuilder) viewLen() int {
	return len(v.rawView)
}

func (v *viewBuilder) render(newText string, style lipgloss.Style) {
	cursorPos := v.model.Cursor()

	viewLen := v.viewLen()
	if cursorPos >= viewLen && cursorPos < viewLen+len(newText) {
		v.view += v.renderAllWithCursor(newText, cursorPos-viewLen, style)
	} else {
		v.view += style.Render(newText)
	}
	v.rawView += newText
}

func (v viewBuilder) last() *byte {
	viewLen := v.viewLen()
	if viewLen == 0 {
		return nil
	}
	last := v.rawView[v.viewLen()-1]
	return &last
}

func (v viewBuilder) renderWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	view := v.cursorView(string(text[cursorPos]), s)
	view += s.Render(text[cursorPos+1:])
	return view
}

func (v viewBuilder) renderAllWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	view := ""
	view += s.Render(text[:cursorPos])
	view += v.renderWithCursor(text, cursorPos, s)
	return view
}

func (v viewBuilder) cursorView(view string, s lipgloss.Style) string {
	if v.model.textinput.Blink() {
		return s.Render(view)
	}
	return v.model.CursorStyle.Inline(true).Reverse(true).Render(view)
}
