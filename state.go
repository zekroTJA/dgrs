package dgrs

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

	// Discord session used to fetch unpresent data
	// and to hook event handlers into.
	DiscordSession DiscordSession

	// Redis client options to connect to a redis
	// instance.
	RedisOptions redis.Options

	// Fetch requested values directly from the Discord API
	// and store them in the cache.
	FetchAndStore bool

	// If set, all cache entries created by dgrs will be
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

	// The prefix used for the redis storage keys.
	//
	// Defaults to 'gdrs'.
	KeyPrefix string
}

// Lifetimes wrap a grid of lifetime specifications
// for each cache object.
type Lifetimes struct {
	OverrrideZero bool

	General,
	Guild,
	Member,
	User,
	Role,
	Channel,
	Emoji,
	Message,
	VoiceState,
	Presence time.Duration
}

// State utilizes a redis connection to be able to store and retrieve
// discordgo state objects.
//
// Also, because state hooks event handlers into the passed discord
// session, it is also possible to maintain the current state
// automatically.
type State struct {
	client  redis.Cmdable
	session DiscordSession
	options *Options
}

// New returns a new State instance with the passed
// options.
func New(opts Options) (s *State, err error) {
	s = &State{}

	s.session = opts.DiscordSession

	if opts.FetchAndStore && s.session == nil {
		err = ErrSessionNotProvided
		return
	}

	if opts.RedisClient != nil {
		s.client = opts.RedisClient
	} else {
		s.client = redis.NewClient(&opts.RedisOptions)
	}

	if opts.KeyPrefix == "" {
		opts.KeyPrefix = defaultKeyPrefix
	}

	s.options = &opts

	if opts.FlushOnStartup {
		err = s.Flush()
	}

	if s.session != nil {
		s.session.AddHandler(func(se *discordgo.Session, e interface{}) {
			if err := s.onEvent(se, e); err != nil {
				log.Println("State Error: ", err)
			}
		})
	}

	return
}

// Flush deletes all keys in the cache stored by dgrs.
//
// You can also specify sub keys like KeyGuild to only remove
// all guild entries, for example.
func (s *State) Flush(subKeys ...string) (err error) {
	subKeys = append(subKeys, "*")
	return s.flush(s.joinKeys(subKeys...))
}

func (s *State) onEvent(_ *discordgo.Session, _e interface{}) (err error) {
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
		for _, p := range e.Presences {
			if err = s.SetPresence(e.GuildID, p); err != nil {
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
		s.SetPresence(e.GuildID, &e.Presence)

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
