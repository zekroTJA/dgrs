package dgrs

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zekrotja/dgrs/mocks"
)

func TestNew(t *testing.T) {
	{
		s, err := New(Options{})
		assert.Nil(t, err)
		assert.Equal(t, defaultKeyPrefix, s.options.KeyPrefix)
	}

	{
		lt := Lifetimes{
			General: 1,
			Guild:   2,
		}
		s, err := New(Options{
			KeyPrefix: "customPrefix",
			Lifetimes: lt,
		})
		assert.Nil(t, err)
		assert.Equal(t, "customPrefix", s.options.KeyPrefix)
		assert.Equal(t, lt, s.options.Lifetimes)
	}

	{
		_, err := New(Options{
			FetchAndStore: true,
		})
		assert.ErrorIs(t, err, ErrSessionNotProvided)

		session := &mocks.DiscordSession{}
		session.On("AddHandler", mock.Anything).Return(func() {})
		_, err = New(Options{
			FetchAndStore:  true,
			DiscordSession: session,
		})
		assert.Nil(t, err)
		session.AssertCalled(t, "AddHandler", mock.Anything)
	}
}

func TestHandlers(t *testing.T) {
	// Dummy session instance, not actually used
	ds := &discordgo.Session{}

	{
		state, _, handler := obtainHookesInstance()
		self := testUser("selfuser")
		guilds := []*discordgo.Guild{
			testGuild("g1"),
			testGuild("g2"),
		}

		handler(ds, &discordgo.Ready{
			User:   self,
			Guilds: guilds,
		})

		rs, err := state.SelfUser()
		assert.Nil(t, err)
		assert.Equal(t, self, rs)

		gr, err := state.Guilds()
		assert.Nil(t, err)
		assert.Equal(t, guilds, gr)
	}
}

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
			KeyPrefix:     fmt.Sprintf("dgrctest%d", rand.Int()),
			MarshalFunc:   json.Marshal,
			UnmarshalFunc: json.Unmarshal,
		},
	}
	return
}

func obtainHookesInstance() (
	state *State,
	session *mocks.DiscordSession,
	handler func(*discordgo.Session, interface{}),
) {
	session = &mocks.DiscordSession{}
	session.On("AddHandler", mock.Anything).
		Run(func(args mock.Arguments) {
			handler = args[0].(func(*discordgo.Session, interface{}))
		}).
		Return(func() {})

	state, _ = New(Options{
		FetchAndStore:  true,
		DiscordSession: session,
	})
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
