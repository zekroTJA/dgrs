package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetMember(guildID string, member *discordgo.Member) (err error) {
	if member.User == nil {
		err = ErrMemberUserNil
		return
	}

	err = s.set(joinKeys(keyMember, guildID, member.User.ID), member, s.getLifetime(member))
	return
}

func (s *State) Member(guildID, memberID string, forceNoFetch ...bool) (v *discordgo.Member, err error) {
	v = &discordgo.Member{}
	ok, err := s.get(joinKeys(keyUser, guildID, memberID), v)
	if !ok {
		if s.options.FetchAndStore && !(len(forceNoFetch) > 0 && forceNoFetch[0]) {
			if v, err = s.session.GuildMember(guildID, memberID); v != nil && err == nil {
				err = s.SetMember(guildID, v)
			}
		} else {
			v = nil
		}
	}

	if v != nil && v.User != nil {
		err = s.SetUser(v.User)
	}

	return
}

func (s *State) Members(guildID string) (v []*discordgo.Member, err error) {
	v = make([]*discordgo.Member, 0)
	err = s.list(joinKeys(keyMember, guildID, "*"), &v)
	return
}

func (s *State) RemoveMember(guildID, memberID string) (err error) {
	return s.del(joinKeys(keyMember, guildID, memberID))
}
