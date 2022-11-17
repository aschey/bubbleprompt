package editor

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ViewBuilder struct {
	view        string
	viewLen     int
	extraOffset int
	rawView     string
	cursor      int
	cursorStyle lipgloss.Style
	delimiter   string
	showCursor  bool
}

func NewViewBuilder(cursor int, cursorStyle lipgloss.Style, delimiter string, showCursor bool) *ViewBuilder {
	return &ViewBuilder{cursor: cursor, cursorStyle: cursorStyle, delimiter: delimiter, showCursor: showCursor}
}

func (v ViewBuilder) View() string {
	if v.cursor == v.viewLen {
		return v.view + v.cursorView(" ", lipgloss.NewStyle())
	}
	return v.view
}

func (v *ViewBuilder) Render(newRunes []rune, column int, style lipgloss.Style) {
	offset := column - 1
	if offset < 0 {
		offset = 0
	}

	cursorPos := v.cursor
	if offset+v.extraOffset > v.viewLen {
		newRunes = append([]rune(strings.Repeat(string(v.delimiter), offset+v.extraOffset-v.viewLen)), newRunes...)
	}
	if cursorPos >= v.viewLen && cursorPos < v.viewLen+len(newRunes) {
		v.view += v.renderAllWithCursor(newRunes, cursorPos-v.viewLen, style)
	} else {
		v.view += style.Render(string(newRunes))
	}
	v.rawView += string(newRunes)
	v.viewLen += len(newRunes)
}

func (v *ViewBuilder) ViewLen() int {
	return v.viewLen
}

func (v *ViewBuilder) RenderPlaceholder(newRunes []rune, offset int, style lipgloss.Style) {
	v.Render(newRunes, offset, style)
	// Add offset to account for the extra characters we added to the view that aren't part of what the user typed
	v.extraOffset += len(newRunes)
}

func (v ViewBuilder) Last() *byte {
	if v.viewLen == 0 {
		return nil
	}
	last := v.rawView[len(v.rawView)-1]
	return &last
}

func (v ViewBuilder) renderWithCursor(runes []rune, cursorPos int, s lipgloss.Style) string {
	view := v.cursorView(string(runes[cursorPos]), s)
	view += s.Render(string(runes[cursorPos+1:]))
	return view
}

func (v ViewBuilder) renderAllWithCursor(runes []rune, cursorPos int, s lipgloss.Style) string {
	view := ""
	view += s.Render(string(runes[:cursorPos]))
	view += v.renderWithCursor(runes, cursorPos, s)
	return view
}

func (v ViewBuilder) cursorView(view string, s lipgloss.Style) string {
	if !v.showCursor {
		return s.Render(view)
	}
	return v.cursorStyle.Inline(true).Reverse(true).Render(view)
}
