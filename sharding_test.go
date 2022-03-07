package dgrs

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestReserveShard(t *testing.T) {
	state, _ := obtainInstance()
	state.options.ShardTimeout = 10 * time.Second

	id, err := state.ReserveShard(0)
	assert.Nil(t, err)
	assert.Equal(t, 0, id)

	time.Sleep(10 * time.Millisecond)

	state.stopHeartbeat = nil
	id, err = state.ReserveShard(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, id)

	time.Sleep(10 * time.Millisecond)

	state.stopHeartbeat = nil
	_, err = state.ReserveShard(0, 1)
	assert.ErrorIs(t, err, ErrShardIDAlreadyReserved)

	err = state.ReleaseShard(0, 0)
	assert.Nil(t, err)

	time.Sleep(10 * time.Millisecond)

	id, err = state.ReserveShard(0)
	assert.Nil(t, err)
	assert.Equal(t, 0, id)
}

func TestShards(t *testing.T) {
	state, _ := obtainInstance()
	state.options.ShardTimeout = 10 * time.Second

	t1 := state.startHeartbeat(0, 1)

	time.Sleep(100 * time.Millisecond)
	shards, err := state.Shards(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(shards))

	t2 := state.startHeartbeat(0, 2)
	defer t2()

	time.Sleep(100 * time.Millisecond)
	shards, err = state.Shards(0)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(shards))

	t1()
	time.Sleep(state.options.ShardTimeout + 5*time.Second)

	shards, err = state.Shards(0)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(shards))
}
