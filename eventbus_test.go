package dgrs

import (
	"context"
	"encoding/json"
	"testing"
	"time"

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

	time.Sleep(1 * time.Second)
}
