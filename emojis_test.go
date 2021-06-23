package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testEmoji(id ...string) *discordgo.Emoji {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.Emoji{
		ID:   gid,
		Name: "emojiname",
	}
}

func TestSetEmoji(t *testing.T) {
	state, _ := obtainInstance()

	emoji := testEmoji()
	err := state.SetEmoji("guildid", emoji)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyEmoji, "guildid", "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(emoji), res.Val())
}

func TestEmoji(t *testing.T) {
	emojis := []*discordgo.Emoji{
		testEmoji("id1"),
		testEmoji("id2"),
	}

	state, session := obtainInstance()

	session.On("GuildEmojis", "guildid").Return(emojis, nil)

	er, err := state.Emoji("guildid", "id1")
	assert.Nil(t, err)
	assert.Nil(t, er)

	state.options.FetchAndStore = true

	er, err = state.Emoji("guildid", "id1")
	assert.Nil(t, err)
	assert.EqualValues(t, emojis[0], er)

	er, err = state.Emoji("guildid", "id2")
	assert.Nil(t, err)
	assert.EqualValues(t, emojis[1], er)
}

func TestEmojis(t *testing.T) {
	emojis := make([]*discordgo.Emoji, 10)

	testEmojis := func(exp []*discordgo.Emoji, rec []*discordgo.Emoji) {
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
			assert.True(t, found, "Expected emoji not found in recovered emojis", eg.ID)
		}
		assert.Equal(t, 10, i, "Not all emojis were recovered")
	}

	for i := range emojis {
		e := testEmoji(fmt.Sprintf("id%d", i))
		emojis[i] = e
	}

	{
		state, _ := obtainInstance()

		for _, e := range emojis {
			assert.Nil(t, state.SetEmoji("guildid", e))
		}

		recEmojis, err := state.Emojis("guildid")
		assert.Nil(t, err)

		testEmojis(emojis, recEmojis)
	}

	{
		state, session := obtainInstance()

		session.On("GuildEmojis", "guildid").Return(emojis, nil)

		resEmojis, err := state.Emojis("guildid")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(resEmojis))

		state.options.FetchAndStore = true
		resEmojis, err = state.Emojis("guildid")
		assert.Nil(t, err)
		testEmojis(emojis, resEmojis)
	}
}

func TestRemoveEmoji(t *testing.T) {
	state, _ := obtainInstance()

	e1 := testEmoji(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetEmoji("guildid", e1))
	e2 := testEmoji(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetEmoji("guildid", e2))

	assert.Nil(t, state.SetEmoji("guildid1", e1))

	assert.Nil(t, state.RemoveEmoji("guildid", e1.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyEmoji, "guildid", e1.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)

	res = state.client.Get(context.Background(), state.joinKeys(KeyEmoji, "guildid", e2.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(e2), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyEmoji, "guildid1", e1.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(e1), res.Val())
}
