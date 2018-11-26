package types

type Error interface {
	Code() int
	Summary() string
}

func NewError(code int, summary string) Error {
	return &innerError{
		code:    code,
		summary: summary,
	}
}

type innerError struct {
	code    int
	summary string
}

func (s *innerError) Code() int {
	return s.code
}

func (s *innerError) Summary() string {
	return s.summary
}
