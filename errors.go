package dgrc

import "errors"

var (
	ErrNilObject     = errors.New("object is nil")
	ErrMemberUserNil = errors.New("user object of member is nil")
	ErrInvalidType   = errors.New("invalid object type")
)
