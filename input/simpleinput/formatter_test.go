package simpleinput_test

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/input/simpleinput"
	"github.com/aschey/bubbleprompt/parser"
	"github.com/charmbracelet/lipgloss"
)

type alternatingFormatter struct {
	lexer         parser.Lexer
	evenTextStyle lipgloss.Style
	oddTextStyle  lipgloss.Style
}

func (f alternatingFormatter) Lex(
	input string,
	_selectedToken *input.Token,
) ([]parser.FormatterToken, error) {
	tokens, err := f.lexer.Lex(input)
	if err != nil {
		return nil, err
	}

	formatterTokens := []parser.FormatterToken{}
	for i, token := range tokens {
		formatterToken := parser.FormatterToken{Value: token.Value}
		if i%2 == 0 {
			formatterToken.Style = f.evenTextStyle
		} else {
			formatterToken.Style = f.oddTextStyle
		}
		formatterTokens = append(formatterTokens, formatterToken)
	}

	return formatterTokens, nil
}

func ExampleWithFormatter() {
	simpleinput.New(simpleinput.WithFormatter[any](func(lexer parser.Lexer) parser.Formatter {
		return alternatingFormatter{
			lexer:         lexer,
			evenTextStyle: lipgloss.NewStyle().Foreground(lipgloss.Color("1")),
			oddTextStyle:  lipgloss.NewStyle().Foreground(lipgloss.Color("2")),
		}
	}))
}
