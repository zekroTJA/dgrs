package main

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/zekrotja/dgrc"
)

func main() {
	s := dgrc.New(nil, dgrc.Options{
		RedisOptions: redis.Options{
			Addr: "localhost:6379",
		},
	})

	fmt.Println(s)
}
