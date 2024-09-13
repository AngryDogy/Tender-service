package myerrors

type NotFoundError struct {
	Text string
}

func (e *NotFoundError) Error() string {
	return e.Text
}

type NoRightsError struct {
	Text string
}

func (e *NoRightsError) Error() string {
	return e.Text
}
