package dropdown

type Option[T any] func(model *Model[T])

func WithMaxSuggestions[T any](maxSuggestions int) Option[T] {
	return func(model *Model[T]) {
		model.SetMaxSuggestions(maxSuggestions)
	}
}

func WithSelectionIndicator[T any](indicator string) Option[T] {
	return func(model *Model[T]) {
		model.SetSelectionIndicator(indicator)
	}
}
