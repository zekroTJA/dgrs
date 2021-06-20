package dgrs

import "github.com/bwmarrin/discordgo"

// DiscordSession describes required API functionalities of
// a discordgo.Session instance.
type DiscordSession interface {
	AddHandler(interface{}) func()
	Channel(channelID string) (*discordgo.Channel, error)
	GuildEmojis(guildID string) ([]*discordgo.Emoji, error)
	Guild(guildID string) (*discordgo.Guild, error)
	GuildMember(guildID, memberID string) (*discordgo.Member, error)
	ChannelMessage(channelID, messageID string) (*discordgo.Message, error)
	GuildRoles(guildID string) ([]*discordgo.Role, error)
	User(userID string) (*discordgo.User, error)
}
