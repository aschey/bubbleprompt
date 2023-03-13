package renderer

type rendererSettings struct {
	widthOffset  int
	heightOffset int
	useHistory   bool
}

type Option func(settings *rendererSettings)

func WithWidthOffset(offset int) Option {
	return func(settings *rendererSettings) {
		settings.widthOffset = offset
	}
}

func WithHeightOffset(offset int) Option {
	return func(settings *rendererSettings) {
		settings.heightOffset = offset
	}
}

func WithUseHistory(useHistory bool) Option {
	return func(settings *rendererSettings) {
		settings.useHistory = useHistory
	}
}
