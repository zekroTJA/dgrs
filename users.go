package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetUser(user *discordgo.User) (err error) {
	err = s.set(joinKeys(keyUser, user.ID), user, s.getLifetime(user))
	return
}

func (s *State) User(id string) (v *discordgo.User, err error) {
	v = &discordgo.User{}
	ok, err := s.get(joinKeys(keyUser, id), v)
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

func (s *State) Users() (v []*discordgo.User, err error) {
	v = make([]*discordgo.User, 0)
	err = s.list(joinKeys(keyUser, "*"), &v)
	return
}

func (s *State) RemoveUser(id string) (err error) {
	return s.del(joinKeys(keyUser, id))
}

func (s *State) SelfUser() (v *discordgo.User, err error) {
	return s.User("@me")
}

func (s *State) SetSelfUser(user *discordgo.User) (err error) {
	err = s.set(joinKeys(keyUser, "@me"), user, s.getLifetime(user))
	return
}
