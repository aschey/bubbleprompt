package prompt

import (
	"github.com/aschey/bubbleprompt/input"
	"github.com/aschey/bubbleprompt/renderer"
)

type Option[T any] func(model *Model[T])

func WithFormatters[T any](formatters input.Formatters) Option[T] {
	return func(model *Model[T]) {
		model.formatters = formatters
	}
}

func WithMaxSuggestions[T any](maxSuggestions int) Option[T] {
	return func(model *Model[T]) {
		model.SetMaxSuggestions(maxSuggestions)
	}
}

func WithUnmanagedRenderer[T any]() Option[T] {
	return func(model *Model[T]) {
		model.renderer = renderer.NewUnmanagedRenderer()
	}
}

func WithViewportRenderer[T any](offset renderer.ViewportOffset) Option[T] {
	return func(model *Model[T]) {
		model.renderer = renderer.NewViewportRenderer(offset)
	}
}

func WithRenderer[T any](r renderer.Renderer) Option[T] {
	return func(model *Model[T]) {
		model.renderer = r
	}
}
