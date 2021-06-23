package dgrs

import "errors"

var (
	ErrNilObject          = errors.New("object is nil")
	ErrUserNil            = errors.New("user object is nil")
	ErrInvalidType        = errors.New("invalid object type")
	ErrSessionNotProvided = errors.New("when FetchAndStore is enabled, a discordgo session instance must be provided")
)
