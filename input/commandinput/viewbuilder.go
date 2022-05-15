package commandinput

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type viewBuilder[T CmdMetadataAccessor] struct {
	view        string
	viewLen     int
	extraOffset int
	rawView     string
	model       Model[T]
}

func newViewBuilder[T CmdMetadataAccessor](model Model[T]) *viewBuilder[T] {
	return &viewBuilder[T]{model: model}
}

func (v viewBuilder[T]) getView() string {
	if v.model.Cursor() == v.viewLen {
		return v.view + v.cursorView(" ", lipgloss.NewStyle())
	}
	return v.view
}

func (v *viewBuilder[T]) render(newText string, offset int, style lipgloss.Style) {
	cursorPos := v.model.Cursor()
	if offset+v.extraOffset > v.viewLen {
		newText = strings.Repeat(v.model.defaultDelimiter, offset+v.extraOffset-v.viewLen) + newText
	}
	if cursorPos >= v.viewLen && cursorPos < v.viewLen+len(newText) {
		v.view += v.renderAllWithCursor(newText, cursorPos-v.viewLen, style)
	} else {
		v.view += style.Render(newText)
	}
	v.rawView += newText
	v.viewLen += len(newText)
}

func (v *viewBuilder[T]) renderPlaceholder(newText string, offset int, style lipgloss.Style) {
	v.render(newText, offset, style)
	// Add offset to account for the extra characters we added to the view that aren't part of what the user typed
	v.extraOffset += len(newText)
}

func (v viewBuilder[T]) last() *byte {
	if v.viewLen == 0 {
		return nil
	}
	last := v.rawView[len(v.rawView)-1]
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
