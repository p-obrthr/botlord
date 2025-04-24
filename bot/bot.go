package bot

import (
	"github.com/bwmarrin/discordgo"
	"botlord/db"
	"botlord/models"
	"os"
	"log"
	"fmt"
	"strings"
)

type Bot struct {
	db *db.BotlordDb
	token string
	textChannelId string
	commands *[]Command
}

func NewBot() *Bot {
	token, exists := os.LookupEnv("DISCORD_BOT_TOKEN") 
	if !exists {
		log.Fatal("err: no discord bot token")
	}
	textChannelId, exists := os.LookupEnv("TEXT_CHANNEL_ID") 
	if !exists {
		log.Fatal("err: not text channel id")
	}
	var err error
	db, err := db.InitDb()
	if err != nil {
		log.Fatalf("err database init: %v\n", err)
	}
	bot := &Bot {
		db: db,
		token: token,
		textChannelId: textChannelId,
	}
	bot.InitCommands()
	return bot
}

type Command struct {
	Trigger string
	Description string
	Use string
	Execute func(s *discordgo.Session, m *discordgo.MessageCreate, args string) error
}

func(b *Bot) InitCommands() {
	b.commands = &[]Command{
		{
			Trigger: "!hi",
			Description: "Gruesst dich zurueck.",
			Use: "",
			Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
				user := fmt.Sprintf("<@%s>", m.Author.ID)
				b.Reply(s, m, fmt.Sprintf("Meddl %s", user))
				return nil
			},
		},
		{
			Trigger: "!addQuote",
			Description: "Fuegt ein Zitat hinzu.",
			Use: "[Zitattext]",
			Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
				if args == "" {
					b.Reply(s, m, "Kein Zitat mitgegeben.")
					return nil
				}
				quote := models.NewQuote(args)
				id, err := b.db.Insert(*quote)
				if err != nil {
					b.Reply(s, m, fmt.Sprintf("err adding quote: %v", err))
					return err
				}
				b.Reply(s, m, "Zitat erfolgreich hinzugefuegt.")
				log.Printf("Quote successfully inserted: id %d", id)
				return nil
			},
		},
		{
			Trigger: "!quote",
			Use: "",
			Description: "Liefert ein zufaelliges Zitat zurueck.",
			Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
				quoteText, err := b.db.GetRandomQuoteText()
				if err != nil {
					b.Reply(s, m, fmt.Sprintf("Fehler: %v", err))
					log.Printf("err: get random quote %v", err)
					return err
				}
				b.Reply(s, m, *quoteText)
				return nil
			},
		},
		{
			Trigger: "!commands",
			Use: "",
			Description: "Gibt eine Liste aller verfuegbaren Kommandos zurueck.",
			Execute: func(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
				commandList := "Verfuegbare Kommandos:\n"
				for _, cmd := range *b.commands {
					commandList += fmt.Sprintf("- %s %s >> %s\n", cmd.Trigger, cmd.Use, cmd.Description)
				}
				b.Reply(s, m, commandList)
				return nil
			},
		},
	}
}

func(b *Bot) Reply(s *discordgo.Session, m *discordgo.MessageCreate, text string) {
	s.ChannelMessageSend(m.ChannelID, text)
}

func(b *Bot) Start() {
	sess, err := discordgo.New("Bot " + b.token)	
	if err != nil {
		log.Fatalf("err init discordgo: %v\n", err)
	}

	sess.AddHandler(b.handleMessage)
	sess.AddHandler(b.handleChannelUpdate)

	sess.Identify.Intents = discordgo.IntentsAllWithoutPrivileged
	err = sess.Open()
	if err != nil {
		log.Fatalf("err: %v\n", err)
	}

	fmt.Println("the bot is online...")
}

func (b *Bot) handleMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "!") {
		return
	}
	content := m.Content
	guiltyCommand := false
	for _, cmd := range *b.commands {
		if strings.HasPrefix(content, cmd.Trigger) {
			guiltyCommand = true
			args := ""
			if parts := strings.SplitN(content, " ", 2); len(parts) == 2 {
				args = parts[1]
			}
			err := cmd.Execute(s, m, args)
			if err != nil {
				log.Printf("err: executing command %s: %v", cmd.Trigger, err)
			}
			break
		}
	}

	if !guiltyCommand {
		b.Reply(s, m, "Ungueltiges Kommando -> siehe !commands fuer gueltige Kommandos und weitere Informationen.")
	}
}

func (b *Bot) handleChannelUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.BeforeUpdate == nil && vs.ChannelID != "" {
		user, err := s.User(vs.UserID)
		if err != nil {
			log.Printf("Fehler beim Abrufen des Benutzers: %v", err)
			return
		}

		msg := fmt.Sprintf("<@%s> ist dem Sprachkanal beigetreten!", user.ID)
		s.ChannelMessageSend(b.textChannelId, msg)
	}
}
