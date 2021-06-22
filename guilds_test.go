package dgrs

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func testGuild() *discordgo.Guild {
	return &discordgo.Guild{
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
}

func TestSetGuild(t *testing.T) {
	state, _ := obtainInstance()

	erroneousGuild := testGuild()
	erroneousGuild.Members[0].User = nil
	err := state.SetGuild(erroneousGuild)
	assert.ErrorIs(t, err, ErrMemberUserNil)

	guild := testGuild()
	err = state.SetGuild(guild)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyGuild, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyMember, "id", "memberid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Members[0]), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyRole, "id", "roleid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Roles[0]), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyChannel, "chanid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Channels[0]), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyEmoji, "id", "emojiid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Emojis[0]), res.Val())
}

func TestGuild(t *testing.T) {
	guild := testGuild()
	guild.Members = make([]*discordgo.Member, 0)
	guild.Roles = make([]*discordgo.Role, 0)
	guild.Channels = make([]*discordgo.Channel, 0)
	guild.Emojis = make([]*discordgo.Emoji, 0)

	state, session := obtainInstance()

	session.On("Guild", "id").Return(guild, nil)

	gr, err := state.Guild("id")
	assert.Nil(t, err)
	assert.Nil(t, gr)

	state.options.FetchAndStore = true
	gr, err = state.Guild("id")
	assert.Nil(t, err)
	assert.EqualValues(t, guild, gr)
}

func TestGuilds(t *testing.T) {

}
