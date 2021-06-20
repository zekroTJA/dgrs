package dgrc

import (
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
)

// Options defines State preferences.
type Options struct {
	// You can pass a pre-initialized redis instance
	// if you already have one.
	RedisClient redis.Cmdable

	// Redis client options to connect to a redis
	// instance.
	RedisOptions redis.Options

	// Fetch requested values directly from the Discord API
	// and store them in the cache.
	FetchAndStore bool

	// If set, all cache entries created by dgrc will be
	// flushed on initialization.
	FlushOnStartup bool

	// You can specify a timeout period for redis commands.
	// If not set, no timeout will be used.
	CommandTimeout time.Duration

	// You can specify either a general lifetime for
	// values stored in the cache or a per-type
	// lifetime which will override the default
	// lifetime for that specific object type.
	//
	// If no lifetime is set at all, a default value
	// of DefaultGeneralLifetime is used.
	Lifetimes Lifetimes
}

type Lifetimes struct {
	General,
	Guild,
	Member,
	User,
	Role,
	Channel,
	Emoji,
	Message,
	VoiceState time.Duration
}

type State struct {
	client  redis.Cmdable
	options *Options
	session *discordgo.Session
}

func New(session *discordgo.Session, opts Options) (s *State, err error) {
	s = &State{}

	s.session = session

	if opts.RedisClient != nil {
		s.client = opts.RedisClient
	} else {
		s.client = redis.NewClient(&opts.RedisOptions)
	}

	s.options = &opts

	if opts.FlushOnStartup {
		err = s.Flush()
	}

	session.AddHandler(func(se *discordgo.Session, e interface{}) {
		if err := s.onEvent(se, e); err != nil {
			log.Println("State Error: ", err)
		}
	})

	return
}

func (s *State) Flush(subKeys ...string) (err error) {
	subKeys = append(subKeys, "*")
	return s.flush(joinKeys(subKeys...))
}

func (s *State) onEvent(se *discordgo.Session, _e interface{}) (err error) {
	switch e := (_e).(type) {

	case *discordgo.Ready:
		for _, g := range e.Guilds {
			if err = s.SetGuild(g); err != nil {
				return
			}
			if err = s.SetSelfUser(e.User); err != nil {
				return
			}
		}

	case *discordgo.GuildCreate:
		err = s.SetGuild(e.Guild)
	case *discordgo.GuildUpdate:
		err = s.SetGuild(e.Guild)
	case *discordgo.GuildDelete:
		err = s.RemoveGuild(e.Guild.ID)

	case *discordgo.GuildMemberAdd:
		err = s.SetMember(e.GuildID, e.Member)
	case *discordgo.GuildMemberUpdate:
		err = s.SetMember(e.GuildID, e.Member)
	case *discordgo.GuildMembersChunk:
		for _, m := range e.Members {
			if err = s.SetMember(e.GuildID, m); err != nil {
				return
			}
		}
	case *discordgo.GuildMemberRemove:
		if e.Member.User != nil {
			err = s.RemoveMember(e.GuildID, e.Member.User.ID)
		}

	case *discordgo.GuildRoleCreate:
		err = s.SetRole(e.GuildID, e.Role)
	case *discordgo.GuildRoleUpdate:
		err = s.SetRole(e.GuildID, e.Role)
	case *discordgo.GuildRoleDelete:
		err = s.RemoveRole(e.GuildID, e.RoleID)

	case *discordgo.GuildEmojisUpdate:
		for _, em := range e.Emojis {
			if err = s.SetEmoji(e.GuildID, em); err != nil {
				return
			}
		}

	case *discordgo.ChannelCreate:
		err = s.SetChannel(e.Channel)
	case *discordgo.ChannelUpdate:
		err = s.SetChannel(e.Channel)
	case *discordgo.ChannelDelete:
		err = s.RemoveChannel(e.Channel.ID)

	case *discordgo.MessageCreate:
		err = s.SetMessage(e.Message)
	case *discordgo.MessageUpdate:
		err = s.SetMessage(e.Message)
	case *discordgo.MessageDelete:
		err = s.RemoveMessage(e.ChannelID, e.Message.ID)
	case *discordgo.MessageDeleteBulk:
		for _, m := range e.Messages {
			if err = s.RemoveMessage(e.ChannelID, m); err != nil {
				return
			}
		}

	case *discordgo.VoiceStateUpdate:
		err = s.SetVoiceState(e.GuildID, e.VoiceState)

	case *discordgo.PresenceUpdate:
		if e.Status == discordgo.StatusOffline {
			return
		}
		var m *discordgo.Member
		m, err = s.Member(e.GuildID, e.User.ID, true)
		if err != nil {
			return
		}
		if m == nil {
			// Member not found; this is a user changing state
			m = &discordgo.Member{
				GuildID: e.GuildID,
				User:    e.User,
			}
		} else {
			if e.User.Username != "" {
				m.User.Username = e.User.Username
			}
		}

		err = s.SetMember(e.GuildID, m)
	}

	return
}
