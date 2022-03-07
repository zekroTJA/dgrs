package dgrs

import (
	"context"
	"encoding/json"
	"testing"

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

	cl := state.Subscribe(
		chKey,
		func(scan func(v interface{}) error) {
			var rec payload
			assert.Nil(t, scan(&rec))
			assert.Equal(t, pl, rec)
		})
	defer cl()

	data, err := json.Marshal(pl)
	assert.Nil(t, err)
	err = state.client.Publish(
		context.Background(),
		chKey,
		data).Err()
	assert.Nil(t, err)
}
