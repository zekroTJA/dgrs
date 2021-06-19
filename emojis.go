package dgrc

import "github.com/bwmarrin/discordgo"

func (s *State) SetEmoji(guildID string, emoji *discordgo.Emoji) (err error) {
	err = s.set(joinKeys(keyEmoji, guildID, emoji.ID), emoji, s.getLifetime(emoji))
	return
}

func (s *State) Emoji(guildID, emojiID string) (v *discordgo.Emoji, err error) {
	v = &discordgo.Emoji{}
	ok, err := s.get(joinKeys(keyEmoji, guildID, emojiID), v)
	if !ok {
		if s.options.FetchAndStore {
			var emojis []*discordgo.Emoji
			if emojis, err = s.session.GuildEmojis(guildID); emojis != nil && err == nil {
				for _, e := range emojis {
					if e.ID == emojiID {
						v = e
					}
					if err = s.SetEmoji(guildID, e); err != nil {
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

func (s *State) Emojis(guildID string) (v []*discordgo.Emoji, err error) {
	v = make([]*discordgo.Emoji, 0)
	err = s.list(joinKeys(keyEmoji, guildID, "*"), &v)
	return
}

func (s *State) RemoveEmoji(guildID, emojiID string) (err error) {
	return s.del(joinKeys(keyEmoji, guildID, emojiID))
}
