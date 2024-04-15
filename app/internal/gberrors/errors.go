package gberrors

type UnsupportedImageFormat struct {
	Format string
}

func (e UnsupportedImageFormat) Error() string {
	return "unsupported image format: " + e.Format
}
