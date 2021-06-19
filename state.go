package dgrc

import (
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
)

// Options defines State preferences.
type Options struct {
	// You can pass a pre-initialized redis instance
	// if you already have one.
	RedisClient *redis.Client

	// Redis client options to connect to a redis
	// instance.
	RedisOptions redis.Options

	// Fetch requested values directly from the Discord API
	// and store them in the cache.
	FetchAndStore bool

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
	options *Options
	client  *redis.Client
	session *discordgo.Session
}

func New(session *discordgo.Session, opts Options) (s *State) {
	s = &State{}

	s.session = session
	if opts.RedisClient != nil {
		s.client = opts.RedisClient
	} else {
		s.client = redis.NewClient(&opts.RedisOptions)
	}

	s.options = &opts

	return
}

func (s *State) SetGuild(guild *discordgo.Guild) (err error) {
	err = s.set(joinKeys(keyGuild, guild.ID), guild, s.getLifetime(guild))
	return
}

func (s *State) Guild(id string) (v *discordgo.Guild, err error) {
	v = &discordgo.Guild{}
	ok, err := s.get(joinKeys(keyGuild, id), v)
	if !ok {
		if s.options.FetchAndStore {
			v, err = s.session.Guild(id)
		} else {
			v = nil
		}
	}
	return
}
