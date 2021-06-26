package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testVoiceState(id ...string) *discordgo.VoiceState {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.VoiceState{
		UserID:    gid,
		SessionID: "sid",
		ChannelID: "chanid",
		GuildID:   "guildid",
	}
}

func TestSetVoiceState(t *testing.T) {
	state, _ := obtainInstance()

	vs := testVoiceState()
	err := state.SetVoiceState("guildid", vs)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyVoiceState, "guildid", "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(vs), res.Val())
}

func TestVoiceState(t *testing.T) {
	vs := testVoiceState()

	state, _ := obtainInstance()

	pr, err := state.VoiceState("guildid", "nonexistent")
	assert.Nil(t, err)
	assert.Nil(t, pr)

	err = state.set(state.joinKeys(KeyVoiceState, "guildid", vs.UserID), vs, state.getLifetime(vs))
	assert.Nil(t, err)

	pr, err = state.VoiceState("guildid", vs.UserID)
	assert.Nil(t, err)
	assert.EqualValues(t, vs, pr)
}

func TestVoiceStates(t *testing.T) {
	vss := make([]*discordgo.VoiceState, 10)

	testVss := func(exp []*discordgo.VoiceState, rec []*discordgo.VoiceState) {
		i := 0
		for _, eg := range exp {
			found := false
			for _, rg := range rec {
				if eg.UserID == rg.UserID {
					assert.Equal(t, eg, rg)
					i++
					found = true
					break
				}
			}
			assert.True(t, found, "Expected voicestate not found in recovered voicestates", eg.UserID)
		}
		assert.Equal(t, 10, i, "Not all voicestates were recovered")
	}

	for i := range vss {
		vs := testVoiceState(fmt.Sprintf("id%d", i))
		vss[i] = vs
	}

	state, _ := obtainInstance()

	for _, vs := range vss {
		assert.Nil(t, state.SetVoiceState("guildid", vs))
	}

	recVss, err := state.VoiceStates("guildid")
	assert.Nil(t, err)

	testVss(vss, recVss)
}

func TestRemoveVoiceState(t *testing.T) {
	state, _ := obtainInstance()

	vs1 := testVoiceState(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetVoiceState("guildid", vs1))
	vs2 := testVoiceState(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetVoiceState("guildid", vs2))

	assert.Nil(t, state.RemoveVoiceState("guildid", vs1.UserID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyVoiceState, "guildid", vs1.UserID))
	assert.ErrorIs(t, res.Err(), redis.Nil)
	res = state.client.Get(context.Background(), state.joinKeys(KeyVoiceState, "guildid", vs2.UserID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(vs2), res.Val())
}
