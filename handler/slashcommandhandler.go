package handler

import (
    "fmt"
    "github.com/bwmarrin/discordgo"
    "github.com/zetto21/discordgo/commands/slashCommands"
)

// 명령 목록을 불러오는 함수
func loadCommands() []*discordgo.ApplicationCommand {
    return []*discordgo.ApplicationCommand{
        slashCommands.Pingcommand,
    }
}

// Slash 명령을 등록하는 함수
func RegisterSlashCommands(dg *discordgo.Session) {
    if dg == nil {
        fmt.Println("오류: Discord 세션이 없습니다.")
        return
    }

    // 명령어 등록
    commands := loadCommands()

    var registeredCount int

    // 명령 등록
    for _, cmd := range commands {
        _, err := dg.ApplicationCommandCreate(dg.State.User.ID, "", cmd)
        if err != nil {
            fmt.Printf("'%s' 명령을 생성할 수 없습니다. %v\n", cmd.Name, err)
        } else {
            registeredCount++
        }
    }

    fmt.Printf("%d개의 명령어가 성공적으로 등록되었습니다.\n", registeredCount)
}
