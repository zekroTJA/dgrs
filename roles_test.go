package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
)

func testRole(id ...string) *discordgo.Role {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.Role{
		ID:   gid,
		Name: "rolename",
	}
}

func TestSetRole(t *testing.T) {
	state, _ := obtainInstance()

	role := testRole()
	err := state.SetRole("guildid", role)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyRole, "guildid", "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(role), res.Val())
}

func TestRole(t *testing.T) {
	roles := []*discordgo.Role{
		testRole("id1"),
		testRole("id2"),
	}

	state, session := obtainInstance()

	session.On("GuildRoles", "guildid").Return(roles, nil)

	er, err := state.Role("guildid", "id1")
	assert.Nil(t, err)
	assert.Nil(t, er)

	state.options.FetchAndStore = true

	er, err = state.Role("guildid", "id1")
	assert.Nil(t, err)
	assert.EqualValues(t, roles[0], er)

	er, err = state.Role("guildid", "id2")
	assert.Nil(t, err)
	assert.EqualValues(t, roles[1], er)
}

func TestRoles(t *testing.T) {
	roles := make([]*discordgo.Role, 10)

	testRoles := func(exp []*discordgo.Role, rec []*discordgo.Role) {
		i := 0
		for _, eg := range exp {
			found := false
			for _, rg := range rec {
				if eg.ID == rg.ID {
					assert.Equal(t, eg, rg)
					i++
					found = true
					break
				}
			}
			assert.True(t, found, "Expected role not found in recovered roles", eg.ID)
		}
		assert.Equal(t, 10, i, "Not all roles were recovered")
	}

	for i := range roles {
		r := testRole(fmt.Sprintf("id%d", i))
		roles[i] = r
	}

	{
		state, _ := obtainInstance()

		for _, r := range roles {
			assert.Nil(t, state.SetRole("guildid", r))
		}

		recRoles, err := state.Roles("guildid")
		assert.Nil(t, err)

		testRoles(roles, recRoles)
	}

	{
		state, session := obtainInstance()

		session.On("GuildRoles", "guildid").Return(roles, nil)

		resRoles, err := state.Roles("guildid")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(resRoles))

		state.options.FetchAndStore = true
		resRoles, err = state.Roles("guildid")
		assert.Nil(t, err)
		testRoles(roles, resRoles)
	}
}

func TestRemoveRole(t *testing.T) {
	state, _ := obtainInstance()

	e1 := testRole(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetRole("guildid", e1))
	e2 := testRole(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetRole("guildid", e2))

	assert.Nil(t, state.SetRole("guildid1", e1))

	assert.Nil(t, state.RemoveRole("guildid", e1.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyRole, "guildid", e1.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)

	res = state.client.Get(context.Background(), state.joinKeys(KeyRole, "guildid", e2.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(e2), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyRole, "guildid1", e1.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(e1), res.Val())
}
