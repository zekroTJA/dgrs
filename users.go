package dgrs

import "github.com/bwmarrin/discordgo"

// SetUser sets the given user object to the cache.
func (s *State) SetUser(user *discordgo.User) (err error) {
	err = s.set(s.joinKeys(KeyUser, user.ID), user, s.getLifetime(user))
	return
}

// User tries to retrieve a user by the given user ID.
//
// If no user was found and FetchAndStore is enabled, the object
// will be tried to be retrieved from the API. When this was successful,
// it is stored in the cache and the object is returned.
//
// Otherwise, if the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
func (s *State) User(id string) (v *discordgo.User, err error) {
	v = &discordgo.User{}
	ok, err := s.get(s.joinKeys(KeyUser, id), v)
	if !ok {
		if s.options.FetchAndStore {
			if v, err = s.session.User(id); v != nil && err == nil {
				err = s.SetUser(v)
			}
		} else {
			v = nil
		}
	}
	return
}

// Users returns a list of users which are stored
// in the cache at the given moment.
func (s *State) Users() (v []*discordgo.User, err error) {
	v = make([]*discordgo.User, 0)
	err = s.list(s.joinKeys(KeyUser, "*"), &v)
	return
}

// RemoveUser removes a user object from the cache by the given ID.
func (s *State) RemoveUser(id string) (err error) {
	return s.del(s.joinKeys(KeyUser, id))
}

// SelfUser returns the current user object of the authenticated
// account.
//
// This object is retrieved on receiving the 'Ready' event.
func (s *State) SelfUser() (v *discordgo.User, err error) {
	return s.User(selfUserKey)
}

// SetSelfUser allows to set a custom user object as self
// user to the cache.
func (s *State) SetSelfUser(user *discordgo.User) (err error) {
	err = s.set(s.joinKeys(KeyUser, selfUserKey), user, s.getLifetime(user))
	return
}

// UserGuilds returns a slice of Guild IDs the user is
// member of.
//
// If forceFetch is passed as true, the list of guilds
// is fetched from cache. There is no deep fetch, that
// means, the particular guilds and members are not
// fetched.
func (s *State) UserGuilds(id string, forceFetch ...bool) (res []string, err error) {
	ok, err := s.get(s.joinKeys(KeyUserGuilds, id), &res)

	if !ok && s.options.FetchAndStore || optBool(forceFetch) {
		res = make([]string, 0)
		var guilds []*discordgo.Guild
		guilds, err = s.Guilds()
		if err != nil {
			return
		}
		var membs []*discordgo.Member
		for _, guild := range guilds {
			membs, err = s.Members(guild.ID)
			if err != nil {
				return
			}
		membLoop:
			for _, memb := range membs {
				if memb.User.ID == id {
					res = append(res, guild.ID)
					break membLoop
				}
			}
		}
		err = s.setUserGuilds(id, res)
	}

	return
}

// AddUserGuilds adds a Guild ID to the list of guild IDs
// the given user is member of.
func (s *State) AddUserGuilds(userID string, guildIDs ...string) (err error) {
	var guilds []string
	if guilds, err = s.UserGuilds(userID); err != nil {
		return
	}

	for _, id := range guildIDs {
		if !stringSliceContains(guilds, id) {
			guilds = append(guilds, id)
		}
	}

	err = s.setUserGuilds(userID, guilds)

	return
}

// RemoveUserGuilds removes a Guild ID from the list of guild IDs
// the given user is member of.
func (s *State) RemoveUserGuilds(userID string, guildIDs ...string) (err error) {
	var guilds []string
	if guilds, err = s.UserGuilds(userID); err != nil {
		return
	}

	for _, id := range guildIDs {
		i := stringSliceIndex(guilds, id)
		if i != -1 {
			guilds = append(guilds[:i], guilds[i+1:]...)
		}
	}

	err = s.setUserGuilds(userID, guilds)

	return
}

func (s *State) setUserGuilds(userID string, guildIDs []string) error {
	return s.set(s.joinKeys(KeyUserGuilds, userID), guildIDs, s.getLifetime((*discordgo.Member)(nil)))
}
