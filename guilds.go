package dgrc

import "github.com/bwmarrin/discordgo"

// SetGuild sets the given guild object to the cache.
//
// If the guild object contains any members, roles, channels
// or emojis, these objects are also stroed subsequently in
// the cache.
func (s *State) SetGuild(guild *discordgo.Guild) (err error) {
	err = s.set(joinKeys(keyGuild, guild.ID), guild, s.getLifetime(guild))
	if err != nil {
		return
	}

	for _, m := range guild.Members {
		if err = s.SetMember(guild.ID, m); err != nil {
			return
		}
	}
	for _, r := range guild.Roles {
		if err = s.SetRole(guild.ID, r); err != nil {
			return
		}
	}
	for _, c := range guild.Channels {
		if err = s.SetChannel(c); err != nil {
			return
		}
	}
	for _, e := range guild.Emojis {
		if err = s.SetEmoji(guild.ID, e); err != nil {
			return
		}
	}

	return
}

// Guild tries to retrieve a guild by the given guild ID.
//
// If no guild was found and FetchAndStore is enabled, the object
// will be tried to be retrieved from the API. When this was successful,
// it is stored in the cache and the object is returned.
//
// Otherwise, if the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
func (s *State) Guild(id string) (v *discordgo.Guild, err error) {
	v = &discordgo.Guild{}
	ok, err := s.get(joinKeys(keyGuild, id), v)
	if !ok {
		if s.options.FetchAndStore {
			if v, err = s.session.Guild(id); v != nil && err == nil {
				err = s.SetGuild(v)
			}
		} else {
			v = nil
		}
	}

	return
}

// Guilds returns a list of guilds which are stored
// in the cache at the given moment.
func (s *State) Guilds() (v []*discordgo.Guild, err error) {
	v = make([]*discordgo.Guild, 0)
	err = s.list(joinKeys(keyGuild, "*"), &v)
	return
}

// RemoveGuild removes a guild object from the cache by the given ID.
func (s *State) RemoveGuild(id string) (err error) {
	return s.del(joinKeys(keyGuild, id))
}
