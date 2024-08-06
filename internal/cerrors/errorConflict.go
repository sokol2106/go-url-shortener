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
var ErrGetShortURLDelete = errors.New("conflict delete original URL")
var ErrGetShortURLNotFind = errors.New("conflict find original URL")
