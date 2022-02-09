package dgrs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestShards(t *testing.T) {
	state, _ := obtainInstance()

	t1 := state.startHeartbeat()

	time.Sleep(100 * time.Millisecond)
	shards, err := state.Shards()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(shards))

	t2 := state.startHeartbeat()
	defer t2()

	time.Sleep(100 * time.Millisecond)
	shards, err = state.Shards()
	assert.Nil(t, err)
	assert.Equal(t, 2, len(shards))

	t1()
	time.Sleep(1*time.Minute + 10*time.Second)

	shards, err = state.Shards()
	assert.Nil(t, err)
	assert.Equal(t, 1, len(shards))
}
