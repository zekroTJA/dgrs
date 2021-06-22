package dgrs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/zekrotja/dgrs/mocks"
)

// ---- HELPERS --------------------------------

func init() {
	rand.Seed(time.Now().Unix())
}

func obtainInstance() (state *State, session *mocks.DiscordSession) {
	godotenv.Load()
	session = &mocks.DiscordSession{}
	state = &State{
		client: redis.NewClient(&redis.Options{
			Addr: os.Getenv("REDIS_ADDRESS"),
		}),
		session: session,
		options: &Options{
			KeyPrefix: fmt.Sprintf("dgrctest%d", rand.Int()),
		},
	}
	return
}

func mustMarshal(v interface{}) string {
	res, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(res)
}

func mustUnmarshal(data string, v interface{}) {
	err := json.Unmarshal([]byte(data), v)
	if err != nil {
		panic(err)
	}
}

func copyObject(src interface{}, dest interface{}) {
	data := mustMarshal(src)
	mustUnmarshal(data, dest)
}
