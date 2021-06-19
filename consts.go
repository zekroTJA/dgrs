package dgrc

import "time"

const (
	DefaultGeneralLifetime = 5 * time.Minute

	keyPrefix     = "dgrc"
	keyGuild      = "guild"
	keyMember     = "member"
	keyUser       = "user"
	keyRole       = "role"
	keyChannel    = "chan"
	keyEmoji      = "emoji"
	keyMessage    = "message"
	keyVoiceState = "vs"

	keySeperator = ':'
)
