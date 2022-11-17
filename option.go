package prompt

import (
	"github.com/aschey/bubbleprompt/editor"
	"github.com/aschey/bubbleprompt/renderer"
)

type Option[T any] func(model *Model[T]) error

func WithFormatters[T any](formatters editor.Formatters) Option[T] {
	return func(model *Model[T]) error {
		model.formatters = formatters
		return nil
	}
}

func WithMaxSuggestions[T any](maxSuggestions int) Option[T] {
	return func(model *Model[T]) error {
		model.SetMaxSuggestions(maxSuggestions)
		return nil
	}
}

func WithUnmanagedRenderer[T any]() Option[T] {
	return func(model *Model[T]) error {
		model.renderer = renderer.NewUnmanagedRenderer()
		return nil
	}
}

func WithViewportRenderer[T any](offset renderer.ViewportOffset) Option[T] {
	return func(model *Model[T]) error {
		model.renderer = renderer.NewViewportRenderer(offset)
		return nil
	}
}

func WithRenderer[T any](r renderer.Renderer) Option[T] {
	return func(model *Model[T]) error {
		model.renderer = r
		return nil
	}
}
