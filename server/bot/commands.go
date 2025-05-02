package bot

import (
	"botlord/models"
	"fmt"
	"strconv"

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
			Trigger:     "!deleteQuote",
			Description: "Loescht ein bestimmtes Zitat.",
			Use:         "[Id]",
			Execute:     b.CmdDeleteQuote,
		},
		{
			Trigger:     "!quotes",
			Description: "Gibt eine Tabelle aller Zitat zurueck",
			Use:         "",
			Execute:     b.CmdListQuotes,
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
	b.AddLog(fmt.Sprintf("User %s greeted", user))
	return nil
}

func (b *Bot) CmdAddQuote(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	if args == "" {
		b.Reply(s, m, "Kein Zitat mitgegeben.")
		b.AddLog("Add quote failed: no quote text provided")
		return nil
	}
	quote := models.NewQuote(args)
	id, err := b.db.Insert(*quote)
	if err != nil {
		b.Reply(s, m, "Fehler beim Hinzufuegen des Zitats.")
		b.AddLog(fmt.Sprintf("err adding quote: %v", err))
		return err
	}
	b.Reply(s, m, "Zitat erfolgreich hinzugefuegt.")
	b.AddLog(fmt.Sprintf("Quote successfully inserted: id %d", id))
	return nil
}

func (b *Bot) CmdDeleteQuote(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	if args == "" {
		b.Reply(s, m, "Keine Id zum Loeschen angegeben.")
		b.AddLog("Delete quote failed: no quote id provided")
		return nil
	}
	id, err := strconv.Atoi(args)
	if err != nil {
		b.AddLog(fmt.Sprintf("err while args converting to string: %v", err))
		return nil
	}
	deleteErr := b.db.Delete(id)
	if deleteErr != nil {
		b.Reply(s, m, "Beim Loeschen des Zitats ist ein Fehler aufgetreten")
		b.AddLog("err while deleting quote")
	}
	b.Reply(s, m, "Zitat erfolgreich geloescht.")
	b.AddLog(fmt.Sprintf("quote id %d deleted successfully", id))

	return nil
}

func (b *Bot) CmdListQuotes(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
    quotes, err := b.db.GetAllQuotes()
    if err != nil {
        b.Reply(s, m, fmt.Sprintf("Fehler beim Abrufen der Zitate: %v", err))
		b.AddLog("err fetching all quotes")
        return err
    }
    formattedQuotes := models.PrintQuotes(quotes)
    b.Reply(s, m, formattedQuotes)
    b.AddLog("List quotes responed successfully")
    return nil
}

func (b *Bot) CmdRandomQuote(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	quoteText, err := b.db.GetRandomQuoteText()
	if err != nil {
		b.Reply(s, m, fmt.Sprintf("Fehler: %v", err))
		b.AddLog(fmt.Sprintf("err: get random quote %v", err))
		return err
	}
	b.Reply(s, m, *quoteText)
	b.AddLog("Random Quote responsed successfully")
	return nil
}

func (b *Bot) CmdListCommands(s *discordgo.Session, m *discordgo.MessageCreate, args string) error {
	commandList := "Verfuegbare Kommandos:\n"
	for _, cmd := range *b.commands {
		commandList += fmt.Sprintf("- %s %s >> %s\n", cmd.Trigger, cmd.Use, cmd.Description)
	}
	b.Reply(s, m, commandList)
	b.AddLog("List Commands responsed successfully")
	return nil
}
