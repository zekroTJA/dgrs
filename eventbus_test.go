package dgrs

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestPublish(t *testing.T) {
	state, _ := obtainInstance()
	const chKey = "testchan"

	type payload struct {
		Data string
	}

	sub := state.client.Subscribe(
		context.Background(),
		state.joinChanKeys(chKey))
	defer sub.Close()

	pl := payload{"data"}
	err := state.Publish(chKey, pl)
	assert.Nil(t, err)

	msg := <-sub.Channel()
	var rec payload
	err = json.Unmarshal([]byte(msg.Payload), &rec)
	assert.Nil(t, err)
	assert.Equal(t, pl, rec)
}

func TestSubscribe(t *testing.T) {
	state, _ := obtainInstance()
	const chKey = "testchan"

	type payload struct {
		Data string
	}

	pl := payload{"data"}
	var rec payload
	defer assert.Equal(t, &pl, &rec)

	cl := state.Subscribe(
		chKey,
		func(scan func(v interface{}) error) {
			assert.Nil(t, scan(&rec))
		})
	defer cl()

	data, err := json.Marshal(pl)
	assert.Nil(t, err)
	err = state.client.Publish(
		context.Background(),
		state.joinChanKeys(chKey),
		data).Err()
	assert.Nil(t, err)

	time.Sleep(100 * time.Millisecond)
}

func TestSubscribeDMs(t *testing.T) {
	var err error
	state, _, handle := obtainHookesInstance()

	chDm := testChannel("dm-id")
	chDm.Type = discordgo.ChannelTypeDM
	err = state.SetChannel(chDm)
	assert.Nil(t, err)
	chNoDm := testChannel("no-dm-id")
	err = state.SetChannel(chNoDm)
	assert.Nil(t, err)

	dm := testMessage("msg-id")
	dm.ChannelID = chDm.ID

	var rec *DirectMessageEvent
	cl := state.SubscribeDMs(func(e *DirectMessageEvent) {
		rec = e
	})
	defer cl()

	rec = nil
	handle(ds, &discordgo.MessageCreate{dm})
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, rec)

	state.options.BroadcastDMs = true

	rec = nil
	handle(ds, &discordgo.MessageCreate{dm})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, dm, rec.Message)
	assert.Equal(t, chDm, rec.Channel)
	assert.False(t, rec.IsUpdate)

	rec = nil
	handle(ds, &discordgo.MessageUpdate{dm, nil})
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, dm, rec.Message)
	assert.Equal(t, chDm, rec.Channel)
	assert.True(t, rec.IsUpdate)

	rec = nil
	dm.ChannelID = chNoDm.ID
	handle(ds, &discordgo.MessageCreate{dm})
	time.Sleep(100 * time.Millisecond)
	assert.Nil(t, rec)
}
