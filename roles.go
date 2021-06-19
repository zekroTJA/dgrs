package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetRole(guildID string, role *discordgo.Role) (err error) {
	err = s.set(joinKeys(keyRole, guildID, role.ID), role, s.getLifetime(role))
	return
}

func (s *State) Role(guildID, roleID string) (v *discordgo.Role, err error) {
	v = &discordgo.Role{}
	ok, err := s.get(joinKeys(keyRole, guildID, roleID), v)
	if !ok {
		if s.options.FetchAndStore {
			var roles []*discordgo.Role
			if roles, err = s.session.GuildRoles(guildID); roles != nil && err == nil {
				for _, r := range roles {
					if err = s.SetRole(guildID, r); err != nil {
						return
					}
				}
			}
		} else {
			v = nil
		}
	}
	return
}

func (s *State) Roles(guildID string) (v []*discordgo.Role, err error) {
	v = make([]*discordgo.Role, 0)
	err = s.list(joinKeys(keyRole, guildID, "*"), &v)
	return
}

func (s *State) RemoveRole(guildID, roleID string) (err error) {
	return s.del(joinKeys(keyRole, guildID, roleID))
}
