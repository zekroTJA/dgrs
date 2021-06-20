package dgrs

import "github.com/bwmarrin/discordgo"

// SetChannel sets the given channel object to the cache.
func (s *State) SetChannel(channel *discordgo.Channel) (err error) {
	err = s.set(joinKeys(KeyChannel, channel.ID), channel, s.getLifetime(channel))
	return
}

// Channel tries to retrieve a channel by the given channel ID.
//
// If no channel was found and FetchAndStore is enabled, the object
// will be tried to be retrieved from the API. When this was successful,
// it is stored in the cache and the object is returned.
//
// Otherwise, if the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
func (s *State) Channel(id string) (v *discordgo.Channel, err error) {
	v = &discordgo.Channel{}
	ok, err := s.get(joinKeys(KeyChannel, id), v)
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

// Channels returns a list of channels of the given guild ID
// which are stored in the cache at the given moment.
func (s *State) Channels(guildID string, forceFetch ...bool) (v []*discordgo.Channel, err error) {
	v = make([]*discordgo.Channel, 0)
	if err = s.list(joinKeys(KeyChannel, "*"), &v); err != nil {
		return
	}

	vg := make([]*discordgo.Channel, 0)
	if guildID != "" {
		for _, c := range v {
			if c.GuildID == guildID {
				vg = append(vg, c)
			}
		}
		v = vg

		if (len(v) == 0 || optBool(forceFetch)) && s.options.FetchAndStore {
			if v, err = s.session.GuildChannels(guildID); err != nil {
				return
			}
			for _, c := range v {
				if err = s.SetChannel(c); err != nil {
					return
				}
			}
		}
	}

	return
}

// RemoveChannel removes a channel object from the cache by the given ID.
func (s *State) RemoveChannel(id string) (err error) {
	return s.del(joinKeys(KeyChannel, id))
}
