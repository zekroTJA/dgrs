package dgrs

import (
	"encoding/json"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/zekrotja/dgrs/mocks"
)

// ---- HELPERS --------------------------------

func obtainInstance() *State {
	godotenv.Load()

	return &State{
		client: redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDRESS"),
		}),
		session: &mocks.DiscordSession{},
		options: &Options{},
	}
}

func mustMarshal(v interface{}) string {
	res, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(res)
}
