package dgrs

import (
	"strconv"
	"time"
)

const shardIdKey = "meta:shardid"
const shardIdLifetime = 1 * time.Minute
const shartHearbeatInterval = 45 * time.Second

type Shard struct {
	ID            int       `json:"id"`
	LastHeartbeat time.Time `json:"lastheartbeat"`
}

func (s *State) sendHearbeat(id int) {
	s.set(s.joinKeys(shardIdKey, strconv.Itoa(id)),
		Shard{
			ID:            id,
			LastHeartbeat: time.Now(),
		},
		shardIdLifetime)
}

func (s *State) startHeartbeat(id int) func() {
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

func (s *State) ReserveShard(cid ...int) (id int, err error) {
	shards, err := s.Shards()
	if err != nil {
		return
	}
	if len(cid) != 0 {
		// Take the passed shard as ID if not
		// already reserved.
		id = cid[0]
		if containsShard(shards, id) {
			err = ErrShardIDAlreadyReserved
			return
		}
	} else {
		// Take the next free shard ID.
		for i := 0; i < len(shards)+1; i++ {
			if !containsShard(shards, i) {
				id = i
			}
		}
	}
	s.stopHeartbeat = s.startHeartbeat(id)
	return
}

func containsShard(shards []*Shard, id int) bool {
	for _, s := range shards {
		if s.ID == id {
			return true
		}
	}
	return false
}
