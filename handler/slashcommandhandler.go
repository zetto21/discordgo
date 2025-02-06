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

    // 모든 샤드에서 명령어 등록
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

    fmt.Printf("총 %d 개의 명령이 성공적으로 등록되었습니다.\n", registeredCount)

    // 기존 명령 가져오기
    existingCommands, err := dg.ApplicationCommands(dg.State.User.ID, "")
    if err != nil {
        fmt.Println("애플리케이션 명령을 가져오는 중 오류가 발생했습니다:", err)
        return
    }

    var deletedCount int
    desiredCommandMap := make(map[string]bool)

    for _, cmd := range commands {
        desiredCommandMap[cmd.Name] = true
    }

    // 필요 없는 명령 삭제
    for _, cmd := range existingCommands {
        if !desiredCommandMap[cmd.Name] {
            err := dg.ApplicationCommandDelete(dg.State.User.ID, "", cmd.ID)
            if err != nil {
                fmt.Printf("명령 '%s' 삭제 중 오류: %v\n", cmd.Name, err)
            } else {
                deletedCount++
            }
        }
    }

    fmt.Printf("총 %d 개의 불필요한 명령이 삭제되었습니다.\n", deletedCount)
}