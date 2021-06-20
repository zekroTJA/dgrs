package dgrc

import (
	"context"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func joinKeys(keys ...string) string {
	n := len(keyPrefix) + len(keys)
	for i := 0; i < len(keys); i++ {
		n += len(keys[i])
	}

	b := strings.Builder{}
	b.Grow(n)
	b.WriteString(keyPrefix)
	for _, s := range keys {
		b.WriteRune(keySeperator)
		b.WriteString(s)
	}

	return b.String()
}

func (s *State) getContext() (ctx context.Context, cf context.CancelFunc) {
	if s.options.CommandTimeout > 0 {
		ctx, cf = context.WithTimeout(context.Background(), s.options.CommandTimeout)
	} else {
		cf = func() {}
		ctx = context.Background()
	}
	return
}

func (s *State) getLifetime(v interface{}) (d time.Duration) {
	lt := s.options.Lifetimes

	switch v.(type) {
	case *discordgo.Guild:
		d = lt.Guild
	case *discordgo.Member:
		d = lt.Member
	case *discordgo.User:
		d = lt.User
	case *discordgo.Role:
		d = lt.Role
	case *discordgo.Channel:
		d = lt.Channel
	case *discordgo.Emoji:
		d = lt.Emoji
	case *discordgo.Message:
		d = lt.Message
	case *discordgo.VoiceState:
		d = lt.VoiceState
	}

	if d < 0 {
		d = lt.General
	}

	return
}
