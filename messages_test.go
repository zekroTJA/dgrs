package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testMessage(id ...string) *discordgo.Message {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	user := &discordgo.User{
		ID:            "userid",
		Username:      "username",
		Discriminator: "discriminator",
	}

	return &discordgo.Message{
		GuildID:   "guildid",
		ID:        gid,
		ChannelID: "chanid",
		Content:   "content",
		Author:    user,
		Member: &discordgo.Member{
			GuildID: "guildid",
			User:    user,
		},
	}
}

func TestSetMessage(t *testing.T) {
	{
		state, _ := obtainInstance()

		message := testMessage()
		err := state.SetMessage(message)
		assert.Nil(t, err)

		res := state.client.Get(context.Background(), state.joinKeys(KeyMessage, "chanid", "id"))
		assert.Nil(t, res.Err())
		assert.Equal(t, mustMarshal(message), res.Val())

		res = state.client.Get(context.Background(), state.joinKeys(KeyMember, "guildid", "userid"))
		assert.Nil(t, res.Err())
		assert.Equal(t, mustMarshal(message.Member), res.Val())

		res = state.client.Get(context.Background(), state.joinKeys(KeyUser, "userid"))
		assert.Nil(t, res.Err())
		assert.Equal(t, mustMarshal(message.Member.User), res.Val())
	}

	{
		state, _ := obtainInstance()

		message := testMessage()
		message.Member = nil
		err := state.SetMessage(message)
		assert.Nil(t, err)

		res := state.client.Get(context.Background(), state.joinKeys(KeyMessage, "chanid", "id"))
		assert.Nil(t, res.Err())
		assert.Equal(t, mustMarshal(message), res.Val())

		res = state.client.Get(context.Background(), state.joinKeys(KeyUser, "userid"))
		assert.Nil(t, res.Err())
		assert.Equal(t, mustMarshal(message.Author), res.Val())
	}
}

func TestMessage(t *testing.T) {
	message := testMessage()

	state, session := obtainInstance()

	session.On("ChannelMessage", "chanid", "id").Return(message, nil)

	mr, err := state.Message("chanid", "id")
	assert.Nil(t, err)
	assert.Nil(t, mr)

	state.options.FetchAndStore = true
	mr, err = state.Message("chanid", "id")
	assert.Nil(t, err)
	assert.EqualValues(t, message, mr)
}

func TestMessages(t *testing.T) {
	messages := make([]*discordgo.Message, 10)

	testMessages := func(exp []*discordgo.Message, rec []*discordgo.Message) {
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
			assert.True(t, found, "Expected message not found in recovered messages", eg.ID)
		}
		assert.Equal(t, 10, i, "Not all messages were recovered")
	}

	for i := range messages {
		m := testMessage(fmt.Sprintf("id%d", i))
		messages[i] = m
	}

	{
		state, _ := obtainInstance()

		for _, m := range messages {
			assert.Nil(t, state.SetMessage(m))
		}

		recMessages, err := state.Messages("chanid")
		assert.Nil(t, err)

		testMessages(messages, recMessages)
	}

	{
		state, session := obtainInstance()

		session.On("ChannelMessages", "chanid", 100, "", "", "").Return(messages, nil)
		session.On("ChannelMessages", "chanid", 100, mock.Anything, "", "").Return([]*discordgo.Message{}, nil)

		recMessages, err := state.Messages("chanid")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(recMessages))

		state.options.FetchAndStore = true
		recMessages, err = state.Messages("chanid")
		assert.Nil(t, err)
		testMessages(messages, recMessages)
	}
}

func TestRemoveMessage(t *testing.T) {
	state, _ := obtainInstance()

	m1 := testMessage(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetMessage(m1))
	m2 := testMessage(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetMessage(m2))

	assert.Nil(t, state.RemoveMessage(m1.ChannelID, m1.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyMessage, m1.ChannelID, m1.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)
	res = state.client.Get(context.Background(), state.joinKeys(KeyMessage, m2.ChannelID, m2.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(m2), res.Val())
}
