package client

import (
    "fmt"
    "sync"

    "github.com/bwmarrin/discordgo"
)

var (
    initialGuilds   = make(map[string]bool) // 초기 길드 상태 저장
    initialGuildsMu sync.Mutex              // 길드 상태 Mutex
)

// Ready 이벤트 핸들러
func Ready(s *discordgo.Session, r *discordgo.Ready) {
    initialGuildsMu.Lock()
    defer initialGuildsMu.Unlock()

    // 초기 길드 상태 저장
    for _, guild := range r.Guilds {
        initialGuilds[guild.ID] = true
    }

    // 서버 수 출력
    serverCount := len(s.State.Guilds)
    fmt.Printf("봇이 성공적으로 실행되었습니다. (서버 수: %d개)\n", serverCount)
}
