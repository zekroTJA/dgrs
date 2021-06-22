package dgrs

import "github.com/bwmarrin/discordgo"

// SetVoiceState sets the given voice state object to the cache.
func (s *State) SetVoiceState(guildID string, vs *discordgo.VoiceState) (err error) {
	err = s.set(s.joinKeys(KeyVoiceState, guildID, vs.UserID), vs, s.getLifetime(vs))
	return
}

// VoiceState tries to retrieve a voice state by the given guild and user ID.
//
// Because voice states are tracked by the 'VoiceStateUpdate' event handler,
// an uncached voice state object will not be retrieved from the API on get.
func (s *State) VoiceState(guildID, userID string) (v *discordgo.VoiceState, err error) {
	v = &discordgo.VoiceState{}
	ok, err := s.get(s.joinKeys(KeyVoiceState, guildID, userID), v)
	if !ok {
		v = nil
	}
	return
}

// VoiceStates returns a list of voice states which are stored
// in the cache at the given moment on the given guild.
func (s *State) VoiceStates(guildID string) (v []*discordgo.VoiceState, err error) {
	v = make([]*discordgo.VoiceState, 0)
	err = s.list(s.joinKeys(KeyVoiceState, guildID, "*"), &v)
	return
}

// RemoveVoiceState removes a voice state object from the
// cache by the given guild and user ID.
func (s *State) RemoveVoiceState(guildID, userID string) (err error) {
	return s.del(s.joinKeys(KeyVoiceState, guildID, userID))
}
