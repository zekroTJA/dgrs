package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/zekrotja/dgrs"
)

func main() {
	godotenv.Load()

	session, _ := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))

	state, err := dgrs.New(dgrs.Options{
		DiscordSession: session,
		RedisOptions: redis.Options{
			Addr: "localhost:6379",
		},
		FetchAndStore:  true,
		FlushOnStartup: true,
	})
	if err != nil {
		log.Fatal(err)
	}

	session.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		log.Printf("Logged in as %s", e.User.String())
	})

	session.AddHandler(func(_ *discordgo.Session, e *discordgo.MessageCreate) {
		if e.Author.Bot || e.Content != "info" {
			return
		}

		guild, err := state.Guild(e.GuildID)
		if err != nil {
			log.Fatal(err)
		}

		channel, err := state.Channel(e.ChannelID)
		if err != nil {
			log.Fatal(err)
		}

		member, err := state.Member(e.GuildID, e.Author.ID)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Guild: %+v\n", guild)
		log.Printf("Channel: %+v\n", channel)
		log.Printf("Member: %+v\n", member)
	})

	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
