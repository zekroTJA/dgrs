package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testUser(id ...string) *discordgo.User {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.User{
		ID:            gid,
		Username:      "username",
		Discriminator: "discriminator",
	}
}

func TestSetUser(t *testing.T) {
	state, _ := obtainInstance()

	user := testUser()
	err := state.SetUser(user)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyUser, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(user), res.Val())
}

func TestUser(t *testing.T) {
	user := testUser()

	state, session := obtainInstance()

	session.On("User", "id").Return(user, nil)

	ur, err := state.User("id")
	assert.Nil(t, err)
	assert.Nil(t, ur)

	state.options.FetchAndStore = true
	ur, err = state.User("id")
	assert.Nil(t, err)
	assert.EqualValues(t, user, ur)
}

func TestUsers(t *testing.T) {
	users := make([]*discordgo.User, 10)
	state, _ := obtainInstance()

	for i := range users {
		u := testUser(fmt.Sprintf("id%d", i))
		users[i] = u
		assert.Nil(t, state.SetUser(u))
	}

	recUser, err := state.Users()
	assert.Nil(t, err)

	i := 0
	for _, eg := range users {
		found := false
		for _, rg := range recUser {
			if eg.ID == rg.ID {
				assert.Equal(t, eg, rg)
				i++
				found = true
				break
			}
		}
		assert.True(t, found, "Expected user not found in recovered users", eg.ID)
	}
	assert.Equal(t, 10, i, "Not all users were recovered")
}

func TestRemoveUser(t *testing.T) {
	state, _ := obtainInstance()

	u1 := testUser(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetUser(u1))
	u2 := testUser(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetUser(u2))

	assert.Nil(t, state.RemoveUser(u1.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyUser, u1.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)
	res = state.client.Get(context.Background(), state.joinKeys(KeyUser, u2.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(u2), res.Val())
}

func TestSelfUser(t *testing.T) {
	state, _ := obtainInstance()
	user := testUser("@me")

	res := state.client.Set(context.Background(), state.joinKeys(KeyUser, selfUserKey), mustMarshal(user), 0)
	assert.Nil(t, res.Err())

	self, err := state.SelfUser()
	assert.Nil(t, err)
	assert.Equal(t, user, self)
}

func TestSetSelfUser(t *testing.T) {
	state, _ := obtainInstance()
	user := testUser("@me")

	err := state.SetSelfUser(user)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyUser, selfUserKey))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(user), res.Val())
}

func TestUserGuilds(t *testing.T) {
	state, _ := obtainInstance()
	gids := []string{"1", "3", "4", "6"}
	negids := []string{"2", "5", "7"}
	for _, gid := range gids {
		guild := &discordgo.Guild{
			ID: gid,
			Members: []*discordgo.Member{
				{User: &discordgo.User{ID: "uid"}},
				{User: &discordgo.User{ID: "uid2"}},
			},
		}
		assert.Nil(t, state.SetGuild(guild))
	}
	for _, gid := range negids {
		guild := &discordgo.Guild{
			ID: gid,
			Members: []*discordgo.Member{
				{User: &discordgo.User{ID: "uid1"}},
				{User: &discordgo.User{ID: "uid2"}},
			},
		}
		assert.Nil(t, state.SetGuild(guild))
	}

	obtGids, err := state.UserGuilds("uid")
	assert.Nil(t, err)
	assert.ElementsMatch(t, obtGids, gids)
}
