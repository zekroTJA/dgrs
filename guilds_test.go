package dgrs

import (
	"context"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestSetGuild(t *testing.T) {
	guild := &discordgo.Guild{
		ID:   "id",
		Name: "name",
		Members: []*discordgo.Member{
			{
				GuildID: "id",
				Nick:    "nick",
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

	state := obtainInstance()

	err := state.SetGuild(guild)
	assert.ErrorIs(t, err, ErrMemberUserNil)

	guild.Members[0].User = &discordgo.User{
		ID:       "userid",
		Username: "username",
	}

	err = state.SetGuild(guild)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), joinKeys(KeyGuild, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild), res.Val())

	res = state.client.Get(context.Background(), joinKeys(KeyMember, "id", "userid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Members[0]), res.Val())

	res = state.client.Get(context.Background(), joinKeys(KeyRole, "id", "roleid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Roles[0]), res.Val())

	res = state.client.Get(context.Background(), joinKeys(KeyChannel, "chanid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Channels[0]), res.Val())

	res = state.client.Get(context.Background(), joinKeys(KeyEmoji, "id", "emojiid"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(guild.Emojis[0]), res.Val())
}
