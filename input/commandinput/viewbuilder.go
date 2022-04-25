package commandinput

import "github.com/charmbracelet/lipgloss"

type viewBuilder[T CmdMetadataAccessor] struct {
	view    string
	rawView string
	model   Model[T]
}

func newViewBuilder[T CmdMetadataAccessor](model Model[T]) *viewBuilder[T] {
	return &viewBuilder[T]{model: model}
}

func (v viewBuilder[T]) getView() string {
	if v.model.Cursor() == v.viewLen() {
		return v.view + v.cursorView(" ", lipgloss.NewStyle())
	}
	return v.view
}

func (v viewBuilder[T]) viewLen() int {
	return len(v.rawView)
}

func (v *viewBuilder[T]) render(newText string, style lipgloss.Style) {
	cursorPos := v.model.Cursor()

	viewLen := v.viewLen()
	if cursorPos >= viewLen && cursorPos < viewLen+len(newText) {
		v.view += v.renderAllWithCursor(newText, cursorPos-viewLen, style)
	} else {
		v.view += style.Render(newText)
	}
	v.rawView += newText
}

func (v viewBuilder[T]) last() *byte {
	viewLen := v.viewLen()
	if viewLen == 0 {
		return nil
	}
	last := v.rawView[v.viewLen()-1]
	return &last
}

func (v viewBuilder[T]) renderWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	view := v.cursorView(string(text[cursorPos]), s)
	view += s.Render(text[cursorPos+1:])
	return view
}

func (v viewBuilder[T]) renderAllWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	view := ""
	view += s.Render(text[:cursorPos])
	view += v.renderWithCursor(text, cursorPos, s)
	return view
}

func (v viewBuilder[T]) cursorView(view string, s lipgloss.Style) string {
	if v.model.textinput.Blink() {
		return s.Render(view)
	}
	return v.model.CursorStyle.Inline(true).Reverse(true).Render(view)
}
