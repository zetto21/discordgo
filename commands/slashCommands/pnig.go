package slashCommands

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

// 봇 정보 명령어 정의
var Pingcommand = &discordgo.ApplicationCommand{
	Name:        "핑",
	Description: "봇 현재 상태를 확인합니다.",
}

func PingSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	latency := s.HeartbeatLatency().Milliseconds()
	response := fmt.Sprintf("퐁! %dms", latency)
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: response,
		},
	})
}