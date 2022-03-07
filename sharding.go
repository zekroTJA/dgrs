package dgrs

import (
	"strconv"
	"time"
)

const shardIdKey = "meta:shardid"

type Shard struct {
	ID            int       `json:"id"`
	Pool          int       `json:"pool"`
	LastHeartbeat time.Time `json:"lastheartbeat"`
}

func (s *State) Shards(pool int) (shards []*Shard, err error) {
	poolStr := strconv.Itoa(pool)
	shards = make([]*Shard, 0)
	err = s.list(s.joinKeys(shardIdKey, poolStr, "*"), &shards)
	return
}

func (s *State) ReserveShard(pool int, cid ...int) (id int, err error) {
	if s.stopHeartbeat != nil {
		err = ErrShardAlreadyAllocated
		return
	}
	shards, err := s.Shards(pool)
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
	s.stopHeartbeat = s.startHeartbeat(pool, id)
	return
}

func (s *State) ReleaseShard(pool, id int) (err error) {
	if s.stopHeartbeat != nil {
		s.stopHeartbeat()
	}
	err = s.del(s.joinKeys(shardIdKey, strconv.Itoa(pool), strconv.Itoa(id)))
	return
}

func (s *State) sendHearbeat(pool, id int) {
	s.set(s.joinKeys(shardIdKey, strconv.Itoa(pool), strconv.Itoa(id)),
		Shard{
			ID:            id,
			LastHeartbeat: time.Now(),
		},
		s.options.ShardTimeout)
}

func (s *State) startHeartbeat(pool, id int) func() {
	d := s.options.ShardTimeout / 4 * 3
	if s.options.ShardTimeout-d > 15*time.Second {
		d = s.options.ShardTimeout - 15*time.Second
	}

	ticker := time.NewTicker(d)
	go func() {
		s.sendHearbeat(pool, id)
		for range ticker.C {
			s.sendHearbeat(pool, id)
		}
	}()

	return ticker.Stop
}

func containsShard(shards []*Shard, id int) bool {
	for _, s := range shards {
		if s.ID == id {
			return true
		}
	}
	return false
}
