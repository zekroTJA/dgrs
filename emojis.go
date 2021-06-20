package dgrc

import "github.com/bwmarrin/discordgo"

// SetEmoji sets the given emoji object to the cache.
func (s *State) SetEmoji(guildID string, emoji *discordgo.Emoji) (err error) {
	err = s.set(joinKeys(keyEmoji, guildID, emoji.ID), emoji, s.getLifetime(emoji))
	return
}

// Emoji tries to retrieve a channel by the given guild and emoji ID.
//
// If no emoji was found and FetchAndStore is enabled, the object
// will be tried to be retrieved from the API. When this was successful,
// it is stored in the cache and the object is returned.
//
// Otherwise, if the object was not found in the cache and FetchAndStore
// is disabled, nil is returned.
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

// Emojis returns a list of emojis of the given guild ID
// which are stored in the cache at the given moment.
func (s *State) Emojis(guildID string) (v []*discordgo.Emoji, err error) {
	v = make([]*discordgo.Emoji, 0)
	err = s.list(joinKeys(keyEmoji, guildID, "*"), &v)
	return
}

// RemoveEmoji removes an emoji object from the cache by the given ID.
func (s *State) RemoveEmoji(guildID, emojiID string) (err error) {
	return s.del(joinKeys(keyEmoji, guildID, emojiID))
}
