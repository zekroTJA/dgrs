package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testGuild(id ...string) *discordgo.Guild {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.Guild{
		ID:   gid,
		Name: "name",
		Members: []*discordgo.Member{
			{
				GuildID: gid,
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
				ID:      "chanid",
				Name:    "channame",
				GuildID: gid,
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
	assert.ErrorIs(t, err, ErrUserNil)

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
	{
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
		assert.Equal(t, guild, gr)
	}

	{
		guild := testGuild()
		guildRet := testGuild()
		guildRet.Members = make([]*discordgo.Member, 0)
		guildRet.Roles = make([]*discordgo.Role, 0)
		guildRet.Channels = make([]*discordgo.Channel, 0)
		guildRet.Emojis = make([]*discordgo.Emoji, 0)

		state, session := obtainInstance()

		session.On("Guild", "id").Return(guildRet, nil)
		m := guild.Members[0]
		state.set(state.joinKeys(KeyMember, guild.ID, m.User.ID), m, 0)
		r := guild.Roles[0]
		state.set(state.joinKeys(KeyRole, guild.ID, r.ID), r, 0)
		c := guild.Channels[0]
		state.set(state.joinKeys(KeyChannel, c.ID), c, 0)
		e := guild.Emojis[0]
		state.set(state.joinKeys(KeyEmoji, guild.ID, e.ID), e, 0)

		gr, err := state.Guild("id", true)
		assert.Nil(t, err)
		assert.Nil(t, gr)

		state.options.FetchAndStore = true
		gr, err = state.Guild("id", true)
		assert.Nil(t, err)
		assert.Equal(t, guild, gr)
	}
}

func TestGuilds(t *testing.T) {
	guilds := make([]*discordgo.Guild, 10)
	state, _ := obtainInstance()

	for i := range guilds {
		g := testGuild(fmt.Sprintf("id%d", i))
		guilds[i] = g
		assert.Nil(t, state.SetGuild(g))
	}

	recGuilds, err := state.Guilds()
	assert.Nil(t, err)

	i := 0
	for _, eg := range guilds {
		found := false
		for _, rg := range recGuilds {
			if eg.ID == rg.ID {
				assert.Equal(t, eg, rg)
				i++
				found = true
				break
			}
		}
		assert.True(t, found, "Expected guild not found in recovered guilds", eg.ID)
	}
	assert.Equal(t, 10, i, "Not all guilds were recovered")
}

func TestRemoveGuild(t *testing.T) {
	{
		state, _ := obtainInstance()

		g1 := testGuild(fmt.Sprintf("id%d", 1))
		assert.Nil(t, state.SetGuild(g1))
		g2 := testGuild(fmt.Sprintf("id%d", 2))
		assert.Nil(t, state.SetGuild(g2))

		assert.Nil(t, state.RemoveGuild(g1.ID))

		res := state.client.Get(context.Background(), state.joinKeys(KeyGuild, g1.ID))
		assert.ErrorIs(t, res.Err(), redis.Nil)
		res = state.client.Get(context.Background(), state.joinKeys(KeyGuild, g2.ID))
		assert.Nil(t, res.Err())
		assert.Equal(t, mustMarshal(g2), res.Val())
	}

	{
		state, _ := obtainInstance()

		g := testGuild()
		assert.Nil(t, state.SetGuild(g))

		membs, err := state.Members(g.ID)
		assert.Nil(t, err)
		assert.ElementsMatch(t, g.Members, membs)

		roles, err := state.Roles(g.ID)
		assert.Nil(t, err)
		assert.ElementsMatch(t, g.Roles, roles)

		chans, err := state.Channels(g.ID)
		assert.Nil(t, err)
		assert.ElementsMatch(t, g.Channels, chans)

		emojis, err := state.Emojis(g.ID)
		assert.Nil(t, err)
		assert.ElementsMatch(t, g.Emojis, emojis)

		err = state.RemoveGuild(g.ID, true)
		assert.Nil(t, err)

		membs, err = state.Members(g.ID)
		assert.Nil(t, err)
		assert.Empty(t, membs)

		roles, err = state.Roles(g.ID)
		assert.Nil(t, err)
		assert.Empty(t, roles)

		chans, err = state.Channels(g.ID)
		assert.Nil(t, err)
		assert.Empty(t, chans)

		emojis, err = state.Emojis(g.ID)
		assert.Nil(t, err)
		assert.Empty(t, emojis)
	}
}
