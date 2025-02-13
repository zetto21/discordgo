// Events/Guilds/messageCreate.go
package guilds

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zetto21/discordgo/utils"
	"github.com/zetto21/discordgo/commands/message"
)

const prefix = "!"

var TextCommands = []utils.Command{
    message.Pingcommand,
}

func MessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // 오류 발생시 로그에 기록합니다.
    defer func() {
        if err := recover(); err != nil {
            log.Printf("오류 발생: %v", err)
        }
    }()

    // 봇이 보낸 메시지 또는 프리픽스가 없는 메시지는 무시합니다.
    if m.Author.Bot || !strings.HasPrefix(m.Content, prefix) {
        return
    }

    // 메시지에서 프리픽스를 제외한 나머지 부분을 명령어 인자로 처리합니다.
    args := strings.Fields(m.Content[len(prefix):])
    if len(args) == 0 {
        return // 명령어가 없는 경우 무시
    }

    // 첫 번째 단어를 명령어로 인식하여 처리
    commandName := strings.ToLower(args[0])

    // 서버 이름을 얻어옵니다.
    guildName := utils.GetGuildName(s, m.GuildID)

    // 명령어가 등록된 명령어 목록에 존재하는지 확인
    for _, cmd := range TextCommands {
        if cmd.Name == commandName || contains(cmd.Aliases, commandName) {
            // 명령어가 있으면 해당 명령을 실행
            cmd.Execute(s, m, args[1:]) // 명령어 실행
            log.Printf("메시지 명령 사용됨: | 명령어: %s | 옵션: %v | 유저: %s (%s) | 길드: %s (%s)", commandName, args[1:], m.Author.Username, m.Author.ID, guildName, m.GuildID)
            return
        }
    }

    // 등록되지 않은 명령어인 경우 경고 로그 출력
    log.Printf("\033[31m알 수 없는 명령: | 명령어: %s | 옵션: %v | 유저: %s (%s) | 길드: %s (%s)\033[0m", commandName, args[1:], m.Author.Username, m.Author.ID, guildName, m.GuildID)
}

func contains(slice []string, item string) bool {
    for _, s := range slice {
        if s == item {
            return true
        }
    }
    return false
}
