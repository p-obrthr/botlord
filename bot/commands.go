package bot

import (
	"botlord/models"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

type Command struct {
	Trigger     string
	Description string
	Use         string
	Execute     func(s *discordgo.Session, m *discordgo.MessageCreate, args string) error
}

func (b *Bot) InitCommands() {
	b.commands = &[]Command{
		{
			Trigger:     "!hi",
			Description: "Gruesst dich zurueck.",
			Use:         "",
			Execute:     b.CmdGreet,
		},
		{
			Trigger:     "!addQuote",
			Description: "Fuegt ein Zitat hinzu.",
			Use:         "[Zitattext]",
			Execute:     b.CmdAddQuote,
		},
		{
			Trigger:     "!quote",
			Use:         "",
			Description: "Liefert ein zufaelliges Zitat zurueck.",
			Execute:     b.CmdRandomQuote,
		},
		{
			Trigger:     "!commands",
			Use:         "",
			Description: "Gibt eine Liste aller verfuegbaren Kommandos zurueck.",
			Execute:     b.CmdListCommands,
		},
	}
}

func (b *Bot) CmdGreet(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	user := fmt.Sprintf("<@%s>", m.Author.ID)
	b.Reply(s, m, fmt.Sprintf("Meddl %s", user))
	return nil
}

func (b *Bot) CmdAddQuote(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
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
}

func (b *Bot) CmdRandomQuote(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	quoteText, err := b.db.GetRandomQuoteText()
	if err != nil {
		b.Reply(s, m, fmt.Sprintf("Fehler: %v", err))
		log.Printf("err: get random quote %v", err)
		return err
	}
	b.Reply(s, m, *quoteText)
	return nil
}

func (b *Bot) CmdListCommands(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	commandList := "Verfuegbare Kommandos:\n"
	for _, cmd := range *b.commands {
		commandList += fmt.Sprintf("- %s %s >> %s\n", cmd.Trigger, cmd.Use, cmd.Description)
	}
	b.Reply(s, m, commandList)
	return nil
}
