# dgrs  [![](https://godoc.org/github.com/zekrotja/dgrs?status.svg)](https://pkg.go.dev/github.com/zekrotja/dgrs) [![Go Report Card](https://goreportcard.com/badge/github.com/zekroTJA/dgrs)](https://goreportcard.com/report/github.com/zekroTJA/dgrs) [![CI](https://github.com/zekroTJA/dgrs/actions/workflows/ci.yml/badge.svg)](https://github.com/zekroTJA/dgrs/actions/workflows/ci.yml) [![Coverage Status](https://coveralls.io/repos/github/zekroTJA/dgrs/badge.svg?branch=master)](https://coveralls.io/github/zekroTJA/dgrs?branch=master)

DiscordGo Redis State - or for short: **dgrs** - is a custom state manager for [DiscordGo](https://github.com/bwmarrin/discordgo) which uses a [Redis](https://redis.io) Instance to store and sync state.

This implementation has four core advantages:

1. The default state management of DiscordGo uses multi-layer maps where all cached objects are stored in the application heap. If you are dealing with a lot of data, this can really increase the load on the applications garbage collector and can eventually reduce the performance of your bot. By storing all of those objects in Redis *(which is also way more optimized for storing large amounts of data and making them quickly accessible)*, your applications GC is not responsible for keeping track of all of these objects.

1. By connecting to the same Redis instance, you can share state across multiple sharded replicas of your bot fairly easily. You can even use the auto shard-ID system to automatically distribute shard IDs across your network. There is also an event bus which is using the Redis PubSub system to sharde events across instances like received direct messages. 

1. As long as your Redis instance is up, the state is persistently cached and you don't need to build up your cache state from scratch at every restart of your bot, which can save a lot of time and unnessecary API calls.

1. You can set cache expirations for each type of state object after which the cached value is invalidated. This is not possible with the default state implementation of DiscordGo.

## Usage

```go
// Create a new DiscordGo session.
session, _ := discordgo.New("Bot " + token)

// Create the State instance passing the
// DiscordGo session and Redis client
// configuration.
state, err := dgrs.New(dgrs.Options{
	DiscordSession: session,
	RedisOptions: redis.Options{
		Addr: "localhost:6379",
	},
	FetchAndStore:  true,
})

guilds, err := state.Guilds()
if err != nil {
    log.Fatal(err)
}

for _, g := range guilds {
    fmt.Println(g.Name)
}
```

# ToDo

- [ ] Optimize state updating
- [ ] Optimize code documentation
- [x] Add more unit tests (> 80% coverage)
- [x] Bind events to state updates
- [x] Add custom marshal/unmarshal function option

---

© 2021 Ringo Hoffmann (zekro Development).  
Covered by the MIT License.
