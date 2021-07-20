package dgrs

import (
	"github.com/bwmarrin/discordgo"
)

// SetMember sets the given member object to the cache.
func (s *State) SetMember(guildID string, member *discordgo.Member) (err error) {
	if member.User == nil {
		err = ErrUserNil
		return
	}

	err = s.set(s.joinKeys(KeyMember, guildID, member.User.ID), member, s.getLifetime(member))
	if err != nil {
		return
	}

	if err = s.SetUser(member.User); err != nil {
		return
	}

	return
}

// Member tries to retrieve a member by the given guild and member ID.
//
// If no member was found and FetchAndStore is enabled, the object
// will be tried to be retrieved from the API. When this was successful,
// it is stored in the cache and the object is returned.
//
// Otherwise, if the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
func (s *State) Member(guildID, memberID string, forceNoFetch ...bool) (v *discordgo.Member, err error) {
	v = &discordgo.Member{}
	ok, err := s.get(s.joinKeys(KeyMember, guildID, memberID), v)
	if !ok {
		if s.options.FetchAndStore && !optBool(forceNoFetch) {
			if v, err = s.session.GuildMember(guildID, memberID); v != nil && err == nil {
				err = s.SetMember(guildID, v)
			}
		} else {
			v = nil
		}
	}

	return
}

// Members returns a list of members of the given guild ID
// which are stored in the cache at the given moment.
func (s *State) Members(guildID string, forceFetch ...bool) (v []*discordgo.Member, err error) {
	v = make([]*discordgo.Member, 0)
	if err = s.list(s.joinKeys(KeyMember, guildID, "*"), &v); err != nil {
		return
	}

	if (len(v) == 0 || optBool(forceFetch)) && s.options.FetchAndStore {
		var last string
		var ms []*discordgo.Member
		for ms == nil || len(ms) > 0 {
			if ms != nil {
				last = ms[len(ms)-1].User.ID
			}
			if ms, err = s.session.GuildMembers(guildID, last, 100); err != nil {
				return
			}
			v = append(v, ms...)
			for _, m := range ms {
				if err = s.SetMember(guildID, m); err != nil {
					return
				}
			}
		}
	}

	return
}

// RemoveMember removes a member object from the cache by
// the given guild and member ID.
func (s *State) RemoveMember(guildID, memberID string) (err error) {
	return s.del(s.joinKeys(KeyMember, guildID, memberID))
}
