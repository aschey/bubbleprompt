package simpleinput

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/charmbracelet/lipgloss"
)

type simpleinputFormatter struct {
	lexer             parser.Lexer
	selectedTextStyle lipgloss.Style
}

func (f simpleinputFormatter) Lex(
	input string,
	selectedToken *input.Token,
) ([]parser.FormatterToken, error) {
	tokens, err := f.lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	formatterTokens := []parser.FormatterToken{}
	for _, token := range tokens {
		formatterToken := parser.FormatterToken{Value: token.Value}
		if selectedToken != nil && selectedToken.Start == token.Start {
			formatterToken.Style = f.selectedTextStyle
		}
		formatterTokens = append(formatterTokens, formatterToken)
	}

	return formatterTokens, nil
}
