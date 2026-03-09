package render

import "fmt"

type RenderError struct {
	Message string
}

func (e RenderError) Error() string {
	return e.Message
}

func IsRenderError(err error) bool {
	_, ok := err.(RenderError)
	return ok
}

func NewRenderError(msg string) error {

	return RenderError{Message: msg}
}

func NewRenderErrorf(format string, args ...any) error {
	return NewRenderError(fmt.Sprintf(format, args...))
}
