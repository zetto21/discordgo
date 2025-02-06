package message

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func PingMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
	latency := s.HeartbeatLatency().Milliseconds()
	response := fmt.Sprintf("퐁! %dms", latency)
	_, _ =  s.ChannelMessageSendReply(m.ChannelID, response, &discordgo.MessageReference{
		MessageID: m.ID,
		ChannelID: m.ChannelID,
		GuildID:   m.GuildID,
	})
}