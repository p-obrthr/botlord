package bot

import (
	"botlord/db"
	"fmt"
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	db            *db.BotlordDb
	token         string
	textChannelId string
	commands      *[]Command
	session       *discordgo.Session
}

func NewBot() *Bot {
	token, exists := os.LookupEnv("DISCORD_BOT_TOKEN")
	if !exists {
		log.Fatal("err: no discord bot token")
	}

	//textChannelId, exists := os.LookupEnv("TEXT_CHANNEL_ID")
	//if !exists {
	//	log.Fatal("err: not text channel id")
	//}

	var err error
	db, err := db.InitDb()
	if err != nil {
		log.Fatalf("err database init: %v\n", err)
	}

	bot := &Bot{
		db:    db,
		token: token,
		//textChannelId: textChannelId,
	}
	bot.InitCommands()
	return bot
}

func (b *Bot) Reply(s *discordgo.Session, m *discordgo.MessageCreate, text string) {
	s.ChannelMessageSend(m.ChannelID, text)
}

func (b *Bot) Start() {
	sess, err := discordgo.New("Bot " + b.token)
	if err != nil {
		log.Fatalf("err init discordgo: %v\n", err)
	}

	sess.AddHandler(b.handleMessage)
	//sess.AddHandler(b.handleChannelUpdate)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = sess.Open()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	b.session = sess
	fmt.Println("the bot is online...")
}

func (b *Bot) Stop() {
	if b.session != nil {
		err := b.session.Close()
		if err != nil {
			log.Printf("error closing session: %v\n", err)
		} else {
			fmt.Println("the bot is offline...")
		}
	}
}
