package dgrs

import (
	"context"
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

var (
	ds = &discordgo.Session{}
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

func TestHandlerReady(t *testing.T) {
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
	assert.Contains(t, gr, guilds[0])
	assert.Contains(t, gr, guilds[1])
}

func TestHandlerGuilds(t *testing.T) {
	state, _, handler := obtainHookesInstance()
	guild := testGuild("id")

	handler(ds, &discordgo.GuildCreate{
		Guild: guild,
	})

	r, err := state.Guild("id")
	assert.Nil(t, err)
	assert.Equal(t, guild, r)

	guild.MemberCount = 5
	handler(ds, &discordgo.GuildUpdate{
		Guild: guild,
	})

	r, err = state.Guild("id")
	assert.Nil(t, err)
	assert.Equal(t, guild, r)
	assert.NotSame(t, guild, r)

	handler(ds, &discordgo.GuildDelete{
		Guild: guild,
	})
	r, err = state.Guild("id")
	assert.Nil(t, err)
	assert.Nil(t, r)
}

func TestHandlerMembers(t *testing.T) {
	state, _, handler := obtainHookesInstance()
	const guildID = "guildid"
	member := testMember("id")
	member.GuildID = guildID

	guild := testGuild(guildID)
	assert.Nil(t, state.SetGuild(guild))
	mcb := guild.MemberCount

	handler(ds, &discordgo.GuildMemberAdd{
		Member: member,
	})
	r, err := state.Member(guildID, "id")
	assert.Nil(t, err)
	assert.Equal(t, member, r)
	rg, err := state.Guild(guildID)
	assert.Nil(t, err)
	assert.Equal(t, mcb+1, rg.MemberCount)

	member.Nick = "Poggers"
	handler(ds, &discordgo.GuildMemberUpdate{
		Member: member,
	})
	r, err = state.Member(guildID, "id")
	assert.Nil(t, err)
	assert.Equal(t, member, r)
	assert.NotSame(t, member, r)
	rg, err = state.Guild(guildID)
	assert.Nil(t, err)
	assert.Equal(t, mcb+1, rg.MemberCount)

	handler(ds, &discordgo.GuildMemberRemove{
		Member: member,
	})
	r, err = state.Member(guildID, "id")
	assert.Nil(t, err)
	assert.Nil(t, r)
	rg, err = state.Guild(guildID)
	assert.Nil(t, err)
	assert.Equal(t, mcb, rg.MemberCount)

	members := []*discordgo.Member{
		testMember("id1"),
		testMember("id2"),
	}
	presences := []*discordgo.Presence{
		testPresence("id1"),
		testPresence("id2"),
	}
	handler(ds, &discordgo.GuildMembersChunk{
		Members:   members,
		Presences: presences,
		GuildID:   guildID,
	})
	r, err = state.Member(guildID, "id1")
	assert.Nil(t, err)
	assert.Equal(t, members[0], r)
	r, err = state.Member(guildID, "id2")
	assert.Nil(t, err)
	assert.Equal(t, members[1], r)
	rp, err := state.Presence(guildID, "id1")
	assert.Nil(t, err)
	assert.Equal(t, presences[0], rp)
	rp, err = state.Presence(guildID, "id2")
	assert.Nil(t, err)
	assert.Equal(t, presences[1], rp)
}

func TestHandlerRoles(t *testing.T) {
	state, _, handler := obtainHookesInstance()
	role := testRole("id")
	const guildID = "guildid"

	handler(ds, &discordgo.GuildRoleCreate{
		GuildRole: &discordgo.GuildRole{
			Role:    role,
			GuildID: guildID,
		},
	})

	r, err := state.Role(guildID, "id")
	assert.Nil(t, err)
	assert.Equal(t, role, r)

	role.Name = "newname"
	handler(ds, &discordgo.GuildRoleUpdate{
		GuildRole: &discordgo.GuildRole{
			Role:    role,
			GuildID: guildID,
		},
	})

	r, err = state.Role(guildID, "id")
	assert.Nil(t, err)
	assert.Equal(t, role, r)
	assert.NotSame(t, role, r)

	handler(ds, &discordgo.GuildRoleDelete{
		RoleID:  role.ID,
		GuildID: guildID,
	})
	r, err = state.Role(guildID, "id")
	assert.Nil(t, err)
	assert.Nil(t, r)
}

func TestFlush(t *testing.T) {
	{
		s, _ := obtainInstance()
		populateState(t, s)
		assert.Nil(t, s.Flush())
		res := s.client.Keys(context.Background(), s.options.KeyPrefix+string(keySeperator)+"*")
		assert.Nil(t, res.Err())
		assert.Equal(t, 0, len(res.Val()))
	}

	{
		s, _ := obtainInstance()
		populateState(t, s)
		res := s.client.Keys(context.Background(), s.options.KeyPrefix+string(keySeperator)+"*")
		lenPre := len(res.Val())
		assert.Nil(t, s.Flush(KeyGuild))
		res = s.client.Keys(context.Background(), s.options.KeyPrefix+string(keySeperator)+"*")
		assert.Nil(t, res.Err())
		assert.Equal(t, lenPre-1, len(res.Val()))
		res = s.client.Keys(context.Background(), s.options.KeyPrefix+string(keySeperator)+KeyGuild+string(keySeperator)+"*")
		assert.Nil(t, res.Err())
		assert.Equal(t, 0, len(res.Val()))
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
		FetchAndStore:  false,
		DiscordSession: session,
		KeyPrefix:      fmt.Sprintf("dgrctest%d", rand.Int()),
	})
	return
}

func populateState(t *testing.T, s *State) {
	t.Helper()
	assert.Nil(t, s.SetGuild(testGuild()))
	assert.Nil(t, s.SetMember("guildid", testMember()))
	assert.Nil(t, s.SetUser(testUser()))
	assert.Nil(t, s.SetChannel(testChannel()))
	assert.Nil(t, s.SetMessage(testMessage()))
	assert.Nil(t, s.SetRole("guildid", testRole()))
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
