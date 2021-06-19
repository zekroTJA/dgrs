package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/go-redis/redis/v8"
	"github.com/zekrotja/dgrc"
)

func main() {
	dc, _ := discordgo.New("Bot NDE5ODM3NDcyMDQ2ODQxODY2.WpvqHQ.bpgniyO7WfBwXwEZ8R_4THJm7Zo")

	s := dgrc.New(dc, dgrc.Options{
		RedisOptions: redis.Options{
			Addr: "localhost:6379",
		},
		FetchAndStore: true,
	})

	fmt.Println(dc.Open())
	defer dc.Close()

	fmt.Println(s.Guild("526196711962705925"))
	fmt.Println(s.Guilds())

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
