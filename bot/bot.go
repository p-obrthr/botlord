package bot

import (
	"github.com/bwmarrin/discordgo"
	"botlord/db"
	"botlord/models"
	"os"
	"strings"
	"log"
	"fmt"
)

type Bot struct {
	db *db.BotlordDb
	token string
	session *discordgo.Session
}

func NewBot() *Bot {
	token, exists := os.LookupEnv("DISCORD_BOT_TOKEN") 
	if !exists {
		log.Fatal("err: no discord bot token")
	}
	var err error
	db, err := db.InitDb()
	if err != nil {
		log.Fatalf("err database init: %v\n", err)
	}
	return &Bot {
		db: db,
		token: token,
	}
}

func(b *Bot) Start() {
	sess, err := discordgo.New("Bot " + b.token)	
	if err != nil {
		log.Fatalf("err init discordgo: %v\n", err)
	}

	sess.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		if m.Author.ID == s.State.User.ID {
			return 
		}


		switch {
		case m.Content == "!hi" :
			s.ChannelMessageSend(m.ChannelID, "hi")

		case strings.HasPrefix(m.Content, "!addQuote") :
			messageArray := strings.SplitN(m.Content, " ", 2)
			quoteText := messageArray[1]
			quote := models.NewQuote(quoteText)
			id, err := b.db.Insert(*quote)
			if err != nil {
			} else {
				s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("quote %s successfully inserted: id %d", quote.Text, id))
				fmt.Printf("quote successfully inserted: id %d", id)
			}
		 

		case m.Content == "!quote" :
			quoteText, err := b.db.GetRandomQuoteText()
			if err != nil {
				fmt.Printf("err: get random quote %v", err)
			}
			s.ChannelMessageSend(m.ChannelID, *quoteText) 
		}	
	})

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged

	err = sess.Open()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	fmt.Println("the bot is online...")
}
