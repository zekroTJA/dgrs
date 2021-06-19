package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetMember(member *discordgo.Member) (err error) {
	if member.User == nil {
		err = ErrMemberUserNil
		return
	}

	err = s.set(joinKeys(keyMember, member.GuildID, member.User.ID), member, s.getLifetime(member))
	return
}

func (s *State) Member(guildID, memberID string) (v *discordgo.Member, err error) {
	v = &discordgo.Member{}
	ok, err := s.get(joinKeys(keyGuild, guildID, memberID), v)
	if !ok {
		if s.options.FetchAndStore {
			if v, err = s.session.GuildMember(guildID, memberID); v != nil && err == nil {
				err = s.SetMember(v)
			}
		} else {
			v = nil
		}
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
