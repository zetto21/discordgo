package utils

import (
	"github.com/bwmarrin/discordgo"
)

type Command struct {
    Name    string
    Aliases []string
    Execute func(*discordgo.Session, *discordgo.MessageCreate, []string)
}

// getGuildName retrieves the name of the guild by its ID
func GetGuildName(s *discordgo.Session, guildID string) string {
	guild, err := s.State.Guild(guildID)
	if err != nil || guild == nil {
		return "Unknown Guild"
	}
	return guild.Name
}
