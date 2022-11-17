package simpleinput

import (
	"strings"

	"github.com/alecthomas/participle/v2/lexer"
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/editor/parser"
	"github.com/charmbracelet/lipgloss"
)

type simpleFormatter struct {
	lexer             *lexer.StatefulDefinition
	selectedTextStyle lipgloss.Style
}

func (f simpleFormatter) Lex(input string, selectedToken *editor.Token) ([]parser.FormatterToken, error) {
	tokens, err := f.lexer.Lex("", strings.NewReader(input))
	if err != nil {
		return nil, err
	}
	lexerTokens, err := lexer.ConsumeAll(tokens)
	if err != nil {
		return nil, err
	}
	formatterTokens := []parser.FormatterToken{}
	for _, token := range lexerTokens {
		formatterToken := parser.FormatterToken{Value: token.Value}
		if selectedToken != nil && selectedToken.Start == token.Pos.Column-1 {
			formatterToken.Style = f.selectedTextStyle
		}
		formatterTokens = append(formatterTokens, formatterToken)
	}

	return formatterTokens, nil
}
