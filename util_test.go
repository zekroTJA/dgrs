package dgrs

import (
	"context"
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestJoinKeys(t *testing.T) {
	sep := string(keySeperator)
	s, _ := New(Options{})

	res := s.joinKeys("a", "b", "cd")
	exp := defaultKeyPrefix + sep + "a" + sep + "b" + sep + "cd"
	assert.Equal(t, exp, res)

	res = s.joinKeys("a")
	exp = defaultKeyPrefix + sep + "a"
	assert.Equal(t, exp, res)

	res = s.joinKeys()
	exp = defaultKeyPrefix
	assert.Equal(t, exp, res)

	s, _ = New(Options{
		KeyPrefix: "kekw",
	})

	res = s.joinKeys("a", "b", "cd")
	exp = "kekw" + sep + "a" + sep + "b" + sep + "cd"
	assert.Equal(t, exp, res)
}

func TestGetContext(t *testing.T) {
	{
		s, _ := New(Options{})

		ctx, fn := s.getContext()
		assert.NotNil(t, fn)
		assert.Same(t, context.Background(), ctx)
	}

	{
		s, _ := New(Options{
			CommandTimeout: 1 * time.Second,
		})

		ctx, fn := s.getContext()
		exCtx, _ := context.WithTimeout(context.Background(), 1*time.Second)
		assert.NotNil(t, fn)
		expDeadline, _ := exCtx.Deadline()
		actDeadline, _ := ctx.Deadline()
		assert.WithinDuration(t, expDeadline, actDeadline, 100*time.Microsecond)
	}
}

func TestGetLifetime(t *testing.T) {
	{
		s, _ := New(Options{})
		testLifetimes(t, s, Lifetimes{})
	}

	{
		s, _ := New(Options{
			Lifetimes: Lifetimes{
				General: 1,
			},
		})
		testLifetimes(t, s, Lifetimes{})
	}

	{
		s, _ := New(Options{
			Lifetimes: Lifetimes{
				OverrrideZero: true,
				General:       1,
			},
		})
		testLifetimes(t, s, Lifetimes{
			Guild:      1,
			Member:     1,
			User:       1,
			Role:       1,
			Channel:    1,
			Emoji:      1,
			Message:    1,
			VoiceState: 1,
			Presence:   1,
		})
	}

	{
		s, _ := New(Options{
			Lifetimes: Lifetimes{
				Guild:      1,
				Member:     2,
				User:       3,
				Role:       4,
				Channel:    5,
				Emoji:      6,
				Message:    7,
				VoiceState: 8,
				Presence:   9,
			},
		})
		testLifetimes(t, s, Lifetimes{
			Guild:      1,
			Member:     2,
			User:       3,
			Role:       4,
			Channel:    5,
			Emoji:      6,
			Message:    7,
			VoiceState: 8,
			Presence:   9,
		})
	}

	{
		s, _ := New(Options{
			Lifetimes: Lifetimes{
				Guild:  1,
				Member: 2,
			},
		})
		testLifetimes(t, s, Lifetimes{
			Guild:  1,
			Member: 2,
		})
	}
}

func TestOptBool(t *testing.T) {
	assert.False(t, optBool(nil))
	assert.False(t, optBool([]bool{}))
	assert.False(t, optBool([]bool{false}))
	assert.False(t, optBool([]bool{false, false}))
	assert.False(t, optBool([]bool{false, true}))

	assert.True(t, optBool([]bool{true}))
	assert.True(t, optBool([]bool{true, false}))
	assert.True(t, optBool([]bool{true, true}))
}

// ---- HELPERS ----

func testLifetimes(t *testing.T, s *State, expLts Lifetimes) {
	t.Helper()

	var d time.Duration

	d = s.getLifetime(&discordgo.Guild{})
	assert.Equal(t, expLts.Guild, d)
	d = s.getLifetime(&discordgo.Member{})
	assert.Equal(t, expLts.Member, d)
	d = s.getLifetime(&discordgo.User{})
	assert.Equal(t, expLts.User, d)
	d = s.getLifetime(&discordgo.Role{})
	assert.Equal(t, expLts.Role, d)
	d = s.getLifetime(&discordgo.Channel{})
	assert.Equal(t, expLts.Channel, d)
	d = s.getLifetime(&discordgo.Emoji{})
	assert.Equal(t, expLts.Emoji, d)
	d = s.getLifetime(&discordgo.Message{})
	assert.Equal(t, expLts.Message, d)
	d = s.getLifetime(&discordgo.VoiceState{})
	assert.Equal(t, expLts.VoiceState, d)
	d = s.getLifetime(&discordgo.Presence{})
	assert.Equal(t, expLts.Presence, d)
}
