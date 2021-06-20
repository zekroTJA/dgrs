package dgrs

import "errors"

var (
	ErrNilObject          = errors.New("object is nil")
	ErrMemberUserNil      = errors.New("user object of member is nil")
	ErrInvalidType        = errors.New("invalid object type")
	ErrSessionNotProvided = errors.New("when FetchAndStore is enabled, a discordgo session instance must be provided")
)
