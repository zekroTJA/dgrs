package dgrs

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

var testGuild = &discordgo.Guild{
	ID:   "id",
	Name: "name",
	Members: []*discordgo.Member{
		{
			GuildID: "memberid",
			Nick:    "nick",
			User: &discordgo.User{
				ID: "memberid",
			},
		},
	},
	Roles: []*discordgo.Role{
		{
			ID:   "roleid",
			Name: "rolename",
		},
	},
	Channels: []*discordgo.Channel{
		{
			ID:   "chanid",
			Name: "channame",
		},
	},
	Emojis: []*discordgo.Emoji{
		{
			ID:   "emojiid",
			Name: "emojiname",
		},
	},
}

func TestSetGuild(t *testing.T) {
	state, _ := obtainInstance()

	guild := new(discordgo.Guild)
	copyObject(testGuild, guild)
	guild.Members[0].User = nil
	err := state.SetGuild(guild)
	assert.ErrorIs(t, err, ErrMemberUserNil)

	err = state.SetGuild(testGuild)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyGuild, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(testGuild), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyMember, "id", "memberid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(testGuild.Members[0]), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyRole, "id", "roleid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(testGuild.Roles[0]), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyChannel, "chanid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(testGuild.Channels[0]), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyEmoji, "id", "emojiid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(testGuild.Emojis[0]), res.Val())
}

func TestGuild(t *testing.T) {
	guild := *testGuild
	guild.Members = make([]*discordgo.Member, 0)
	guild.Roles = make([]*discordgo.Role, 0)
	guild.Channels = make([]*discordgo.Channel, 0)
	guild.Emojis = make([]*discordgo.Emoji, 0)

	state, session := obtainInstance()

	session.On("Guild", "id").Return(&guild, nil)

	gr, err := state.Guild("id")
	assert.Nil(t, err)
	assert.Nil(t, gr)

	state.options.FetchAndStore = true
	gr, err = state.Guild("id")
	assert.Nil(t, err)
	assert.EqualValues(t, &guild, gr)
}

func TestGuilds(t *testing.T) {

}
