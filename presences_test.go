package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testPresence(id ...string) *discordgo.Presence {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.Presence{
		Status: discordgo.StatusOnline,
		User: &discordgo.User{
			ID:            gid,
			Username:      "username",
			Discriminator: "discriminator",
		},
	}
}

func TestSetPresence(t *testing.T) {
	state, _ := obtainInstance()

	presence := testPresence()
	presence.User = nil
	err := state.SetPresence("guildid", presence)
	assert.ErrorIs(t, err, ErrUserNil)

	presence = testPresence()
	err = state.SetPresence("guildid", presence)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyPresence, "guildid", "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(presence), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyUser, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(presence.User), res.Val())
}

func TestPresence(t *testing.T) {
	presence := testPresence()

	state, _ := obtainInstance()

	err := state.set(state.joinKeys(KeyPresence, "guildid", presence.User.ID), presence, state.getLifetime(presence))
	assert.Nil(t, err)

	pr, err := state.Presence("guildid", presence.User.ID)
	assert.Nil(t, err)
	assert.EqualValues(t, presence, pr)
}

func TestPresences(t *testing.T) {
	presences := make([]*discordgo.Presence, 10)

	testPresences := func(exp []*discordgo.Presence, rec []*discordgo.Presence) {
		i := 0
		for _, eg := range exp {
			found := false
			for _, rg := range rec {
				if eg.User.ID == rg.User.ID {
					assert.Equal(t, eg, rg)
					i++
					found = true
					break
				}
			}
			assert.True(t, found, "Expected presence not found in recovered presences", eg.User.ID)
		}
		assert.Equal(t, 10, i, "Not all presences were recovered")
	}

	for i := range presences {
		p := testPresence(fmt.Sprintf("id%d", i))
		presences[i] = p
	}

	state, _ := obtainInstance()

	for _, p := range presences {
		assert.Nil(t, state.SetPresence("guildid", p))
	}

	recPresences, err := state.Presences("guildid")
	assert.Nil(t, err)

	testPresences(presences, recPresences)
}

func TestRemovePresence(t *testing.T) {
	state, _ := obtainInstance()

	m1 := testPresence(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetPresence("guildid", m1))
	m2 := testPresence(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetPresence("guildid", m2))

	assert.Nil(t, state.RemovePresence("guildid", m1.User.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyPresence, "guildid", m1.User.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)
	res = state.client.Get(context.Background(), state.joinKeys(KeyPresence, "guildid", m2.User.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(m2), res.Val())
}
