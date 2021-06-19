package dgrc

import "errors"

var (
	ErrNilObject   = errors.New("object is nil")
	ErrInvalidType = errors.New("invalid object type")
)
