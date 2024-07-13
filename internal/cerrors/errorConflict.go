package cerrors

import "errors"

/*
type errorConflict struct {
	Err error
}

func NewConflictError() error {
	return &errorConflict{Err: errors.New("conflict add original URL ")}
}

func (e *errorConflict) Error() string {
	return e.Err.Error()
}

func (e *errorConflict) Unwrap() error {
	return e.Err
}
*/

var ErrNewShortURL = errors.New("conflict add original URL")
