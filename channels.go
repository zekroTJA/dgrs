package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetChannel(channel *discordgo.Channel) (err error) {
	err = s.set(joinKeys(keyChannel, channel.ID), channel, s.getLifetime(channel))
	return
}

func (s *State) Channel(id string) (v *discordgo.Channel, err error) {
	v = &discordgo.Channel{}
	ok, err := s.get(joinKeys(keyChannel, id), v)
	if !ok {
		if s.options.FetchAndStore {
			if v, err = s.session.Channel(id); v != nil && err == nil {
				err = s.SetChannel(v)
			}
		} else {
			v = nil
		}
	}
	return
}

func (s *State) Channels(guildID string) (v []*discordgo.Channel, err error) {
	v = make([]*discordgo.Channel, 0)
	err = s.list(joinKeys(keyChannel, "*"), &v)

	vg := make([]*discordgo.Channel, 0)
	if guildID != "" {
		for _, c := range v {
			if c.GuildID == guildID {
				vg = append(vg, c)
			}
		}
		v = vg
	}

	return
}

func (s *State) RemoveChannel(id string) (err error) {
	return s.del(joinKeys(keyChannel, id))
}
