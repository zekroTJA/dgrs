package dgrs

import (
	"context"
	"fmt"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func testMember(id ...string) *discordgo.Member {
	gid := "id"
	if len(id) > 0 {
		gid = id[0]
	}

	return &discordgo.Member{
		GuildID: "guildid",
		User: &discordgo.User{
			ID:            gid,
			Username:      "username",
			Discriminator: "discriminator",
		},
	}
}

func TestSetMember(t *testing.T) {
	state, _ := obtainInstance()

	member := testMember()
	member.User = nil
	err := state.SetMember(member.GuildID, member)
	assert.ErrorIs(t, err, ErrUserNil)

	member = testMember()
	err = state.SetMember(member.GuildID, member)
	assert.Nil(t, err)

	res := state.client.Get(context.Background(), state.joinKeys(KeyMember, member.GuildID, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(member), res.Val())

	res = state.client.Get(context.Background(), state.joinKeys(KeyUser, "id"))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(member.User), res.Val())
}

func TestMember(t *testing.T) {
	member := testMember()

	state, session := obtainInstance()

	session.On("GuildMember", "guildid", "id").Return(member, nil)

	mr, err := state.Member("guildid", "id")
	assert.Nil(t, err)
	assert.Nil(t, mr)

	state.options.FetchAndStore = true
	mr, err = state.Member("guildid", "id")
	assert.Nil(t, err)
	assert.EqualValues(t, member, mr)
}

func TestMembers(t *testing.T) {
	members := make([]*discordgo.Member, 10)

	testMembers := func(exp []*discordgo.Member, rec []*discordgo.Member) {
		i := 0
		for _, eg := range exp {
			found := false
			for _, rg := range rec {
				if eg.User.ID == rg.User.ID {
					assert.Equal(t, eg, rg)
					i++
					found = true
					break
				}
			}
			assert.True(t, found, "Expected member not found in recovered members", eg.User.ID)
		}
		assert.Equal(t, 10, i, "Not all members were recovered")
	}

	for i := range members {
		m := testMember(fmt.Sprintf("id%d", i))
		members[i] = m
	}

	{
		state, _ := obtainInstance()

		for _, m := range members {
			assert.Nil(t, state.SetMember(m.GuildID, m))
		}

		recMembers, err := state.Members("guildid")
		assert.Nil(t, err)

		testMembers(members, recMembers)
	}

	{
		state, session := obtainInstance()

		session.On("GuildMembers", "guildid", "", 100).Return(members, nil)
		session.On("GuildMembers", "guildid", mock.Anything, 100).Return([]*discordgo.Member{}, nil)

		recMembers, err := state.Members("guildid")
		assert.Nil(t, err)
		assert.Equal(t, 0, len(recMembers))

		state.options.FetchAndStore = true
		recMembers, err = state.Members("guildid")
		assert.Nil(t, err)
		testMembers(members, recMembers)
	}
}

func TestRemoveMember(t *testing.T) {
	state, _ := obtainInstance()

	m1 := testMember(fmt.Sprintf("id%d", 1))
	assert.Nil(t, state.SetMember(m1.GuildID, m1))
	m2 := testMember(fmt.Sprintf("id%d", 2))
	assert.Nil(t, state.SetMember(m2.GuildID, m2))

	assert.Nil(t, state.RemoveMember(m1.GuildID, m1.User.ID))

	res := state.client.Get(context.Background(), state.joinKeys(KeyMember, m1.GuildID, m1.User.ID))
	assert.ErrorIs(t, res.Err(), redis.Nil)
	res = state.client.Get(context.Background(), state.joinKeys(KeyMember, m2.GuildID, m2.User.ID))
	assert.Nil(t, res.Err())
	assert.Equal(t, mustMarshal(m2), res.Val())
}
