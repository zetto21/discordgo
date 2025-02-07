package main

import (
	"os"
	"os/signal"
	"fmt"
	"syscall"

	"github.com/zetto21/discordgo/handler"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)


func main() {
    // .env 파일 로드
    err := godotenv.Load()
    if err != nil {
        fmt.Println(".env 파일 로드 실패. 환경 변수를 확인하세요.")
        return
    }

    // 환경 변수에서 토큰 가져오기
    token := os.Getenv("token")
    if token == "" {
        fmt.Println("에러: 디스코드 봇 토큰이 .env 파일에 설정되지 않았습니다.")
        return
    }

    // Discord 세션 생성
    dg, err := discordgo.New("Bot " + token)
    if err != nil {
        fmt.Println("Discord 세션 생성 중 오류 발생:", err)
        return
    }

    // Discord Intents 설정
    dg.Identify.Intents = discordgo.Intent(16291)

    // Ready 이벤트 핸들러 등록
    dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
        serverCount := len(s.State.Guilds)
        fmt.Printf("봇이 성공적으로 실행되었습니다. (서버 수: %d개)\n", serverCount)
    })

    // 이벤트 핸들러 등록
    handler.RegisterEventHandlers(dg)

    // Discord 연결 열기
    err = dg.Open()
    if err != nil {
        fmt.Println("연결을 여는 중 오류가 발생했습니다:", err)
        return
    }
    defer func() {
        dg.Close()
        fmt.Println("연결이 종료되었습니다.")
    }()

    // 슬래시 명령어 등록
    handler.RegisterSlashCommands(dg)

    // 종료 신호 대기
    stop := make(chan os.Signal, 1)
    signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
    <-stop

    fmt.Println("프로그램이 종료됩니다.")
}
