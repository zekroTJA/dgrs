package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testChannel(id ...string) *discordgo.Channel {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.Channel{
		ID:      gid,
		Name:    "channame",
		GuildID: "guildid",
	}
}

func TestSetChannel(t *testing.T) {
	state, _ := obtainInstance()

	channel := testChannel()
	err := state.SetChannel(channel)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyChannel, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(channel), res.Val())
}

func TestChannel(t *testing.T) {
	channel := testChannel()

	state, session := obtainInstance()

	session.On("Channel", "id").Return(channel, nil)

	gr, err := state.Channel("id")
	assert.Nil(t, err)
	assert.Nil(t, gr)

	state.options.FetchAndStore = true
	gr, err = state.Channel("id")
	assert.Nil(t, err)
	assert.EqualValues(t, channel, gr)
}

func TestChannels(t *testing.T) {
	channels := make([]*discordgo.Channel, 10)

	testChannels := func(exp []*discordgo.Channel, rec []*discordgo.Channel) {
		i := 0
		for _, eg := range exp {
			found := false
			for _, rg := range rec {
				if eg.ID == rg.ID {
					assert.Equal(t, eg, rg)
					i++
					found = true
					break
				}
			}
			assert.True(t, found, "Expected channel not found in recovered channels", eg.ID)
		}
		assert.Equal(t, 10, i, "Not all channels were recovered")
	}

	for i := range channels {
		c := testChannel(fmt.Sprintf("id%d", i))
		channels[i] = c
	}

	{
		state, _ := obtainInstance()

		for _, c := range channels {
			assert.Nil(t, state.SetChannel(c))
		}

		recChannels, err := state.Channels("guildid")
		assert.Nil(t, err)

		testChannels(channels, recChannels)
	}

	{
		state, session := obtainInstance()

		session.On("GuildChannels", "guildid").Return(channels, nil)

		recChannels, err := state.Channels("guildid")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(recChannels))

		state.options.FetchAndStore = true
		recChannels, err = state.Channels("guildid")
		assert.Nil(t, err)
		testChannels(channels, recChannels)
	}
}

func TestRemoveChannel(t *testing.T) {
	state, _ := obtainInstance()

	c1 := testChannel(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetChannel(c1))
	c2 := testChannel(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetChannel(c2))

	assert.Nil(t, state.RemoveChannel(c1.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyChannel, c1.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)
	res = state.client.Get(context.Background(), state.joinKeys(KeyChannel, c2.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(c2), res.Val())
}
