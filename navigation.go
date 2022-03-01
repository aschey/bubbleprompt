package prompt

func (m *Model) unselectSuggestion() {
	m.listPosition = -1
}

func (m Model) isSuggestionSelected() bool {
	return m.listPosition > -1
}

func (m *Model) updateArgs() {
	m.textInput.Args = []string{}
	suggestion := m.getSelectedSuggestion()
	for _, arg := range suggestion.PositionalArgs {
		m.textInput.Args = append(m.textInput.Args, arg.Placeholder)
	}
}

func (m *Model) nextSuggestion() {
	if len(m.completer.suggestions) == 0 {
		return
	}

	if m.listPosition < len(m.completer.suggestions)-1 {
		m.listPosition++
		m.updateArgs()
	} else {
		m.textInput.Args = []string{}
		m.unselectSuggestion()
	}
}

func (m *Model) previousSuggestion() {
	if len(m.completer.suggestions) == 0 {
		return
	}

	m.listPosition--
	if m.isSuggestionSelected() {
		m.updateArgs()
	} else {
		m.textInput.Args = []string{}
		m.unselectSuggestion()
	}
}

func (m Model) getSelectedSuggestion() *Suggestion {
	if m.isSuggestionSelected() {
		return &m.completer.suggestions[m.listPosition]
	}
	return nil
}
