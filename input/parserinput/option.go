package parserinput

type Option func(model *LexerModel) error

func WithDelimiterTokens(tokens ...string) Option {
	return func(model *LexerModel) error {
		model.delimiterTokens = tokens
		return nil
	}
}

func WithDelimiters(delimiters ...string) Option {
	return func(model *LexerModel) error {
		model.delimiters = delimiters
		return nil
	}
}
