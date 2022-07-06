package input

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
	blink       bool
}

func NewViewBuilder(cursor int, cursorStyle lipgloss.Style, delimiter string, blink bool) *ViewBuilder {
	return &ViewBuilder{cursor: cursor, cursorStyle: cursorStyle, delimiter: delimiter, blink: blink}
}

func (v ViewBuilder) GetView() string {
	if v.cursor == v.viewLen {
		return v.view + v.cursorView(" ", lipgloss.NewStyle())
	}
	return v.view
}

func (v *ViewBuilder) Render(newText string, offset int, style lipgloss.Style) {
	cursorPos := v.cursor
	if offset+v.extraOffset > v.viewLen {
		newText = strings.Repeat(v.delimiter, offset+v.extraOffset-v.viewLen) + newText
	}
	if cursorPos >= v.viewLen && cursorPos < v.viewLen+len(newText) {
		v.view += v.renderAllWithCursor(newText, cursorPos-v.viewLen, style)
	} else {
		v.view += style.Render(newText)
	}
	v.rawView += newText
	v.viewLen += len(newText)
}

func (v *ViewBuilder) ViewLen() int {
	return v.viewLen
}

func (v *ViewBuilder) RenderPlaceholder(newText string, offset int, style lipgloss.Style) {
	v.Render(newText, offset, style)
	// Add offset to account for the extra characters we added to the view that aren't part of what the user typed
	v.extraOffset += len(newText)
}

func (v ViewBuilder) Last() *byte {
	if v.viewLen == 0 {
		return nil
	}
	last := v.rawView[len(v.rawView)-1]
	return &last
}

func (v ViewBuilder) renderWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	view := v.cursorView(string(text[cursorPos]), s)
	view += s.Render(text[cursorPos+1:])
	return view
}

func (v ViewBuilder) renderAllWithCursor(text string, cursorPos int, s lipgloss.Style) string {
	view := ""
	view += s.Render(text[:cursorPos])
	view += v.renderWithCursor(text, cursorPos, s)
	return view
}

func (v ViewBuilder) cursorView(view string, s lipgloss.Style) string {
	if v.blink {
		return s.Render(view)
	}
	return v.cursorStyle.Inline(true).Reverse(true).Render(view)
}
