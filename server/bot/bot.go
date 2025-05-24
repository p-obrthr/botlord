package bot

import (
	"botlord/db"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Bot struct {
	db            *db.BotlordDb
	token         string
	textChannelId string
	commands      *[]Command
	session       *discordgo.Session
	Logs          []string
}

func NewBot() *Bot {
	token, exists := os.LookupEnv("DISCORD_BOT_TOKEN")
	if !exists {
		log.Fatal("err: no discord bot token")
	}

	textChannelId, exists := os.LookupEnv("TEXT_CHANNEL_ID")

	var err error
	db, err := db.InitDb()
	if err != nil {
		log.Fatalf("err database init: %v\n", err)
	}

	bot := &Bot{
		db:            db,
		token:         token,
		textChannelId: textChannelId,
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
		b.AddLog(fmt.Sprintf("err init discordgo: %v\n", err))
	}

	sess.AddHandler(b.handleMessage)

	if b.textChannelId != "" {
		sess.AddHandler(b.handleChannelUpdate)
	}

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = sess.Open()
	if err != nil {
		b.AddLog(fmt.Sprintf("error closing session: %v\n", err))
	}

	b.session = sess
	b.AddLog("the bot is online...")
}

func (b *Bot) Stop() {
	if b.session != nil {
		err := b.session.Close()
		if err != nil {
			b.AddLog(fmt.Sprintf("error closing session: %v\n", err))
		} else {
			b.AddLog("the bot is offline...")
		}
	}
}

func (b *Bot) AddLog(log string) {
	location, err := time.LoadLocation("Europe/Berlin")
	if err != nil {
		location = time.UTC
	}
	timestamp := time.Now().In(location).Format("15:04:05")
	formattedLog := fmt.Sprintf("[%s] %s", timestamp, log)
	fmt.Println(formattedLog)
	b.Logs = append(b.Logs, formattedLog)
}
