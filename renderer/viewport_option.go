package renderer

type viewportSettings struct {
	widthOffset  int
	heightOffset int
	useHistory   bool
}

type ViewportOption func(settings *viewportSettings)

func WithWidthOffset(offset int) ViewportOption {
	return func(settings *viewportSettings) {
		settings.widthOffset = offset
	}
}

func WithHeightOffset(offset int) ViewportOption {
	return func(settings *viewportSettings) {
		settings.heightOffset = offset
	}
}

func WithUseHistory(useHistory bool) ViewportOption {
	return func(settings *viewportSettings) {
		settings.useHistory = useHistory
	}
}
