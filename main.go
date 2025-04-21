package main

import (
	"github.com/bwmarrin/discordgo"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"log"
)


func main() {
	token, exists := os.LookupEnv("DISCORD_BOT_TOKEN") 
	if !exists {
		log.Fatal("err: no discord bot token")
	}
	sess, err := discordgo.New("Bot " + token)	
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return 
		}

		if m.Content == "!hi" {
			s.ChannelMessageSend(m.ChannelID, "hi")
		}
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}
	defer sess.Close()

	fmt.Println("the bot is online...")

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}
