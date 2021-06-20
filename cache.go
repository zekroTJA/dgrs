package dgrc

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

func (s *State) set(key string, v interface{}, lifetime time.Duration) (err error) {
	data, err := json.Marshal(v)
	if err != nil {
		return
	}
	ctx, cf := s.getContext()
	defer cf()
	res := s.client.Set(ctx, key, data, lifetime)
	return res.Err()
}

func (s *State) get(key string, v interface{}) (ok bool, err error) {
	ctx, cf := s.getContext()
	defer cf()
	res := s.client.Get(ctx, key)
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
	ctx, cf := s.getContext()
	defer cf()
	res := s.client.Del(ctx, keys...)
	return res.Err()
}

func (s *State) list(key string, v interface{}) (err error) {
	ctx, cf := s.getContext()
	defer cf()
	keys := s.client.Keys(ctx, joinKeys(keyGuild, "*"))
	if err = keys.Err(); err != nil {
		return
	}

	var vals []interface{}

	if len(keys.Val()) > 0 {
		ctx, cf := s.getContext()
		defer cf()
		res := s.client.MGet(ctx, keys.Val()...)
		if err = res.Err(); err != nil {
			return
		}
		vals = res.Val()
	}

	b := strings.Builder{}
	b.WriteRune('[')

	n := len(vals)
	if n > 0 {
		b.WriteString(vals[0].(string))
		if n > 1 {
			for _, v := range vals[1:] {
				b.WriteRune(',')
				b.WriteString(v.(string))
			}
		}
	}

	b.WriteRune(']')

	err = json.Unmarshal([]byte(b.String()), v)
	return
}
