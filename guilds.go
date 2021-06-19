package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetGuild(guild *discordgo.Guild) (err error) {
	err = s.set(joinKeys(keyGuild, guild.ID), guild, s.getLifetime(guild))
	return
}

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

	if v != nil {
		for _, m := range v.Members {
			if err = s.SetMember(id, m); err != nil {
				return
			}
		}
		for _, r := range v.Roles {
			if err = s.SetRole(id, r); err != nil {
				return
			}
		}
		for _, c := range v.Channels {
			if err = s.SetChannel(c); err != nil {
				return
			}
		}
		for _, e := range v.Emojis {
			if err = s.SetEmoji(id, e); err != nil {
				return
			}
		}
	}

	return
}

func (s *State) Guilds() (v []*discordgo.Guild, err error) {
	v = make([]*discordgo.Guild, 0)
	err = s.list(joinKeys(keyGuild, "*"), &v)
	return
}

func (s *State) RemoveGuild(id string) (err error) {
	return s.del(joinKeys(keyGuild, id))
}
