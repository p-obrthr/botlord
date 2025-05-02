package bot

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

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
				b.AddLog(fmt.Sprintf("err: executing command %s: %v", cmd.Trigger, err))
			}
			break
		}
	}

	if !guiltyCommand {
		b.Reply(s, m, "Ungueltiges Kommando -> siehe !commands fuer gueltige Kommandos und weitere Informationen.")
		b.AddLog(fmt.Sprintf("user requested unguilty cmd"))
	}
}

func (b *Bot) handleChannelUpdate(s *discordgo.Session, vs *discordgo.VoiceStateUpdate) {
	if vs.BeforeUpdate == nil && vs.ChannelID != "" {
		user, err := s.User(vs.UserID)
		if err != nil {
			b.AddLog(fmt.Sprintf("err while fetching user data: %v", err))
			return
		}

		fmt.Println(vs.ChannelID)

		msg := fmt.Sprintf("<@%s> ist dem Sprachkanal beigetreten!", user.ID)
		s.ChannelMessageSend(b.textChannelId, msg)
		b.AddLog(fmt.Sprintf("notified text chat that somebody entered the voice channel"))
	}
}
