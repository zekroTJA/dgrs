package dgrc

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
)

func (s *State) set(key string, v interface{}, lifetime time.Duration) (err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	res := s.client.SetEX(s.getContext(), key, data, lifetime)
	return res.Err()
}

func (s *State) get(key string, v interface{}) (ok bool, err error) {
	res := s.client.Get(s.getContext(), key)
	data, err := res.Bytes()
	if err == redis.Nil {
		err = nil
		return
	}
	ok = true
	err = json.Unmarshal(data, v)
	return
}

func (s *State) del(keys ...string) (err error) {
	res := s.client.Del(s.getContext(), keys...)
	return res.Err()
}
