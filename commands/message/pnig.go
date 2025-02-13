package message

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/zetto21/discordgo/utils"
)

// 봇 정보 명령어 정의
var Pingcommand = utils.Command{
    Name:    "핑",
    Aliases: []string{"vld", "ping","ㅍ"},
    Execute: PingMessageCommand,
}

func PingMessageCommand(s *discordgo.Session, m *discordgo.MessageCreate, args []string) {
    latency := s.HeartbeatLatency().Milliseconds()
    response := fmt.Sprintf("퐁! %dms", latency)
    _, _ = s.ChannelMessageSendReply(m.ChannelID, response, &discordgo.MessageReference{
        MessageID: m.ID,
        ChannelID: m.ChannelID,
        GuildID:   m.GuildID,
    })
}
