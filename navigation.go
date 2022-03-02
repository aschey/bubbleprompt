package prompt

func (m *Model) unselectSuggestion() {
	m.listPosition = -1
}

func (m Model) isSuggestionSelected() bool {
	return m.listPosition > -1
}

func (m *Model) nextSuggestion() {
	if len(m.completer.suggestions) == 0 {
		return
	}

	if m.listPosition < len(m.completer.suggestions)-1 {
		m.listPosition++
	} else {
		m.unselectSuggestion()
	}
}

func (m *Model) previousSuggestion() {
	if len(m.completer.suggestions) == 0 {
		return
	}

	m.listPosition--
	if m.isSuggestionSelected() {
	} else {
		m.unselectSuggestion()
	}
}

func (m Model) getSelectedSuggestion() *Suggestion {
	if m.isSuggestionSelected() {
		return &m.completer.suggestions[m.listPosition]
	}
	return nil
}
