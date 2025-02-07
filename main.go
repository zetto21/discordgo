package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

    "github.com/zetto21/discordgo/handler"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
)

// 슬래시 명령어 등록을 위한 뮤텍스 추가
var registerCommandsMutex sync.Mutex

// runShard 함수: 개별 샤드 실행
func runShard(shardID, totalShards int, token string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Discord Intents 설정
	intents := discordgo.Intent(16291)

	// Discord 세션 생성
	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		fmt.Printf("샤드 %d: Discord 세션 생성 중 오류 발생: %v\n", shardID, err)
		return
	}

	

	// 샤드 정보 설정
	dg.Identify.Intents = intents
	dg.Identify.Shard = &[2]int{shardID, totalShards}

	// Ready 이벤트 핸들러 등록
	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		serverCount := len(s.State.Guilds)
		fmt.Printf("샤드 %d: 봇이 성공적으로 실행되었습니다. (서버 수: %d개)\n", shardID, serverCount)
	})

	// 이벤트 핸들러 등록
	handler.RegisterEventHandlers(dg)

	// Discord 연결 열기
	err = dg.Open()
	if err != nil {
		fmt.Printf("샤드 %d: 연결을 여는 중 오류가 발생했습니다: %v\n", shardID, err)
		return
	}
	defer func() {
		dg.Close()
		fmt.Printf("샤드 %d: 연결이 종료되었습니다.\n", shardID)
	}()

	// 첫 번째 샤드에서만 슬래시 명령어 등록
	if shardID == 0 {
		registerCommandsMutex.Lock()
		handler.RegisterSlashCommands(dg)
		registerCommandsMutex.Unlock()
	}

	// 종료 신호 대기
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-stop

	fmt.Printf("샤드 %d: 프로그램이 종료됩니다.\n", shardID)
}

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

	// 샤드 수 설정 (3개 이상으로 설정)
	totalShards := 2 // 원하는 샤드 수로 변경 가능
    fmt.Printf("수동 샤드 수: %d\n", totalShards)

	// WaitGroup으로 샤드 실행 관리
	var wg sync.WaitGroup
	for shardID := 0; shardID < totalShards; shardID++ {
		wg.Add(1)
		go runShard(shardID, totalShards, token, &wg)
	}

	fmt.Println("모든 샤드가 실행 중입니다. 종료하려면 Ctrl+C를 누르세요.")
	wg.Wait() // 모든 샤드가 종료될 때까지 대기
	fmt.Println("모든 샤드가 종료되었습니다. 프로그램을 종료합니다.")
}
