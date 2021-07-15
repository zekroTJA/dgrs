package dgrs

import "github.com/bwmarrin/discordgo"

// SetGuild sets the given guild object to the cache.
//
// If the guild object contains any members, roles, channels
// or emojis, these objects are also stroed subsequently in
// the cache.
func (s *State) SetGuild(guild *discordgo.Guild) (err error) {
	err = s.set(s.joinKeys(KeyGuild, guild.ID), guild, s.getLifetime(guild))
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
//
// Optionally, when hydrate is set to true, members, roles,
// channels and emojis will also be obtained from cache and
// added to the guild object.
func (s *State) Guild(id string, hydrate ...bool) (v *discordgo.Guild, err error) {
	v = &discordgo.Guild{}
	ok, err := s.get(s.joinKeys(KeyGuild, id), v)
	if !ok {
		if s.options.FetchAndStore {
			if v, err = s.session.Guild(id); v != nil && err == nil {
				err = s.SetGuild(v)
			}
		} else {
			v = nil
		}
	}

	if v != nil && optBool(hydrate) {
		if v.Members, err = s.Members(id); err != nil {
			return
		}
		if v.Roles, err = s.Roles(id); err != nil {
			return
		}
		if v.Channels, err = s.Channels(id); err != nil {
			return
		}
		if v.Emojis, err = s.Emojis(id); err != nil {
			return
		}
	}

	return
}

// Guilds returns a list of guilds which are stored
// in the cache at the given moment.
func (s *State) Guilds() (v []*discordgo.Guild, err error) {
	v = make([]*discordgo.Guild, 0)
	err = s.list(s.joinKeys(KeyGuild, "*"), &v)
	return
}

// RemoveGuild removes a guild object from the cache by the given ID.
//
// When dehydrate is passed as true, objects linked to this guild
// (members, roles, voice states, emojis, channels and messages)
// are purged from cache as well.
func (s *State) RemoveGuild(id string, dehydrate ...bool) (err error) {
	if err = s.del(s.joinKeys(KeyGuild, id)); err != nil {
		return
	}

	if optBool(dehydrate) {
		if err = s.delPattern(s.joinKeys(KeyMember, id, "*")); err != nil {
			return
		}
		if err = s.delPattern(s.joinKeys(KeyRole, id, "*")); err != nil {
			return
		}
		if err = s.delPattern(s.joinKeys(KeyVoiceState, id, "*")); err != nil {
			return
		}
		if err = s.delPattern(s.joinKeys(KeyEmoji, id, "*")); err != nil {
			return
		}
		var chans []*discordgo.Channel
		if chans, err = s.Channels(id); err != nil {
			return
		}
		for _, c := range chans {
			if err = s.RemoveChannel(c.ID, true); err != nil {
				return
			}
		}
	}

	return
}
