package dgrs

import (
	"math/rand"
	"strconv"
	"time"
)

const shardIdKey = "meta:shardid"
const shardIdLifetime = 1 * time.Minute
const shartHearbeatInterval = 45 * time.Second

type Shard struct {
	ID            string    `json:"id"`
	LastHeartbeat time.Time `json:"lastheartbeat"`
}

func getRandomID() string {
	return strconv.Itoa(rand.Intn(999_999))
}

func (s *State) sendHearbeat(id string) {
	s.set(s.joinKeys(shardIdKey, id),
		Shard{
			ID:            id,
			LastHeartbeat: time.Now(),
		},
		shardIdLifetime)
}

func (s *State) startHeartbeat() func() {
	id := getRandomID()
	ticker := time.NewTicker(45 * time.Second)
	go func() {
		s.sendHearbeat(id)
		for range ticker.C {
			s.sendHearbeat(id)
		}
	}()
	return ticker.Stop
}

func (s *State) Shards() (shards []*Shard, err error) {
	shards = make([]*Shard, 0)
	err = s.list(s.joinKeys(shardIdKey, "*"), &shards)
	return
}
