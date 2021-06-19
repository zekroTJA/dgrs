package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetMessage(message *discordgo.Message) (err error) {
	err = s.set(joinKeys(keyMessage, message.ChannelID, message.ID), message, s.getLifetime(message))
	return
}

func (s *State) Message(channelID, messageID string) (v *discordgo.Message, err error) {
	v = &discordgo.Message{}
	ok, err := s.get(joinKeys(keyMessage, channelID, messageID), v)
	if !ok {
		if s.options.FetchAndStore {
			if v, err = s.session.ChannelMessage(channelID, messageID); v != nil && err == nil {
				err = s.SetMessage(v)
			}
		} else {
			v = nil
		}
	}
	return
}

func (s *State) Messages(channelID string) (v []*discordgo.Message, err error) {
	v = make([]*discordgo.Message, 0)
	err = s.list(joinKeys(keyMessage, channelID, "*"), &v)
	return
}

func (s *State) RemoveMessage(channelID, messageID string) (err error) {
	return s.del(joinKeys(keyMessage, channelID, messageID))
}
