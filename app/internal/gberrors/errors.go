package gberrors

type UnsupportedImageFormat struct {
	Format string
}

func (e UnsupportedImageFormat) Error() string {
	return "unsupported image format: " + e.Format
}

type MemberNotFound struct {
	Username string
}

func (e MemberNotFound) Error() string {
	return "member not found: " + e.Username
}
