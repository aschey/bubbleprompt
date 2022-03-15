package prompt

type Document struct {
	Text           string
	CursorPosition int
}

func (d Document) TextBeforeCursor() string {
	if d.CursorPosition >= len(d.Text) {
		return d.Text
	}
	return d.Text[:d.CursorPosition]
}
