package dgrs

import (
	"github.com/bwmarrin/discordgo"
)

// SetPresence sets the given presence object to the cache.
func (s *State) SetPresence(guildID string, presence *discordgo.Presence) (err error) {
	if presence.User == nil {
		err = ErrUserNil
		return
	}

	err = s.set(s.joinKeys(KeyPresence, guildID, presence.User.ID), presence, s.getLifetime(presence))
	if err != nil {
		return
	}

	err = s.SetUser(presence.User)

	return
}

// Presence tries to retrieve a presence by the given guild and user ID.
//
// If the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
func (s *State) Presence(guildID, userID string) (v *discordgo.Presence, err error) {
	v = &discordgo.Presence{}
	ok, err := s.get(s.joinKeys(KeyPresence, guildID, userID), v)
	if !ok {
		v = nil
	}

	return
}

// Presences returns a list of presences of the given guild ID
// which are stored in the cache at the given moment.
func (s *State) Presences(guildID string) (v []*discordgo.Presence, err error) {
	v = make([]*discordgo.Presence, 0)
	err = s.list(s.joinKeys(KeyPresence, guildID, "*"), &v)
	return
}

// RemovePresence removes a presence object from the cache by
// the given guild and user ID.
func (s *State) RemovePresence(guildID, userID string) (err error) {
	return s.del(s.joinKeys(KeyPresence, guildID, userID))
}
