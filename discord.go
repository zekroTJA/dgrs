package dgrs

import "github.com/bwmarrin/discordgo"

// DiscordSession describes required API functionalities of
// a discordgo.Session instance.
type DiscordSession interface {
	AddHandler(interface{}) func()
	Channel(channelID string, options ...discordgo.RequestOption) (*discordgo.Channel, error)
	GuildChannels(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Channel, error)
	GuildEmojis(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Emoji, error)
	Guild(guildID string, options ...discordgo.RequestOption) (*discordgo.Guild, error)
	GuildMember(guildID, memberID string, options ...discordgo.RequestOption) (*discordgo.Member, error)
	GuildMembers(guildID string, after string, limit int, options ...discordgo.RequestOption) ([]*discordgo.Member, error)
	ChannelMessage(channelID, messageID string, options ...discordgo.RequestOption) (*discordgo.Message, error)
	ChannelMessages(channelID string, limit int, beforeID, afterID, aroundID string, options ...discordgo.RequestOption) ([]*discordgo.Message, error)
	GuildRoles(guildID string, options ...discordgo.RequestOption) ([]*discordgo.Role, error)
	User(userID string, options ...discordgo.RequestOption) (*discordgo.User, error)
}
