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

	dc.AddHandler(func(_ *discordgo.Session, e *discordgo.Ready) {
		fmt.Println("Guilds:")
		guilds, err := s.Guilds()
		if err != nil {
			fmt.Println("Err: ", err)
			return
		}
		for _, g := range guilds {
			fmt.Println(" -", g.Name)
		}
	})

	fmt.Println(dc.Open())
	defer dc.Close()

	// fmt.Println(s.Guild("526196711962705925"))
	// fmt.Println(s.Guilds())
	// fmt.Println(s.Member("526196711962705925", "221905671296253953"))

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
}
