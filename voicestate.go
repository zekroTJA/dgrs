package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetVoiceState(guildID string, vs *discordgo.VoiceState) (err error) {
	err = s.set(joinKeys(keyVoiceState, guildID, vs.UserID), vs, s.getLifetime(vs))
	return
}

func (s *State) VoiceState(guildID, userID string) (v *discordgo.VoiceState, err error) {
	v = &discordgo.VoiceState{}
	ok, err := s.get(joinKeys(keyVoiceState, guildID, userID), v)
	if !ok {
		v = nil
	}
	return
}

func (s *State) VoiceStates(guildID string) (v []*discordgo.VoiceState, err error) {
	v = make([]*discordgo.VoiceState, 0)
	err = s.list(joinKeys(keyVoiceState, guildID, "*"), &v)
	return
}

func (s *State) RemoveVoiceState(guildID, userID string) (err error) {
	return s.del(joinKeys(keyVoiceState, guildID, userID))
}
