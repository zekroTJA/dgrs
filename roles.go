package dgrc

import "github.com/bwmarrin/discordgo"

// SetRole sets the given role object to the cache.
func (s *State) SetRole(guildID string, role *discordgo.Role) (err error) {
	err = s.set(joinKeys(keyRole, guildID, role.ID), role, s.getLifetime(role))
	return
}

// Role tries to retrieve a role by the given guild and role ID.
//
// If no role was found and FetchAndStore is enabled, the object
// will be tried to be retrieved from the API. When this was successful,
// it is stored in the cache and the object is returned.
//
// Otherwise, if the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
func (s *State) Role(guildID, roleID string) (v *discordgo.Role, err error) {
	v = &discordgo.Role{}
	ok, err := s.get(joinKeys(keyRole, guildID, roleID), v)
	if !ok {
		if s.options.FetchAndStore {
			var roles []*discordgo.Role
			if roles, err = s.session.GuildRoles(guildID); roles != nil && err == nil {
				for _, r := range roles {
					if r.ID == roleID {
						v = r
					}
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

// Roles returns a list of roles which are stored
// in the cache at the given moment on the given guild.
func (s *State) Roles(guildID string) (v []*discordgo.Role, err error) {
	v = make([]*discordgo.Role, 0)
	err = s.list(joinKeys(keyRole, guildID, "*"), &v)
	return
}

// RemoveRole removes a role object from the cache by the given ID.
func (s *State) RemoveRole(guildID, roleID string) (err error) {
	return s.del(joinKeys(keyRole, guildID, roleID))
}
