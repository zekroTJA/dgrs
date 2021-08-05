package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/zekrotja/dgrs"
)

var cmds = map[string]func(s *discordgo.Session, e *discordgo.MessageCreate, state *dgrs.State){
	"info": func(s *discordgo.Session, e *discordgo.MessageCreate, state *dgrs.State) {
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
	},

	"channels": func(s *discordgo.Session, e *discordgo.MessageCreate, state *dgrs.State) {
		chans, err := state.Channels(e.GuildID)
		if err != nil {
			log.Fatal(err)
		}

		for _, c := range chans {
			fmt.Println(c.Name)
		}
	},

	"members": func(s *discordgo.Session, e *discordgo.MessageCreate, state *dgrs.State) {
		membs, err := state.Members(e.GuildID, true)
		if err != nil {
			log.Fatal(err)
		}

		for _, m := range membs {
			fmt.Println(m.User.String())
		}
	},

	"messages": func(s *discordgo.Session, e *discordgo.MessageCreate, state *dgrs.State) {
		membs, err := state.Messages(e.ChannelID, true)
		if err != nil {
			log.Fatal(err)
		}

		for _, m := range membs {
			fmt.Printf("%s - %s\n", m.Author.String(), m.Content)
		}
	},
}

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

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageCreate) {
		if e.Author.Bot {
			return
		}

		if cmd, ok := cmds[strings.ToLower(e.Content)]; ok {
			cmd(s, e, state)
		}
	})

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
		time.Sleep(1000 * time.Millisecond)
		msg, _ := state.Message(e.ChannelID, e.MessageID)
		fmt.Println("--- ADD ------------")
		for _, m := range msg.Reactions {
			fmt.Println(m.Emoji.Name, m.Count)
		}
	})

	session.AddHandler(func(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
		time.Sleep(1000 * time.Millisecond)
		msg, _ := state.Message(e.ChannelID, e.MessageID)
		fmt.Println("--- REM ------------")
		for _, m := range msg.Reactions {
			fmt.Println(m.Emoji.Name, m.Count)
		}
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
