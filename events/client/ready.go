package client

import (
	"fmt"
	"log"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	initialGuilds   = make(map[string]bool) // 초기 길드 상태 저장
	initialGuildsMu sync.Mutex              // 길드 상태 Mutex
	shardReadyCount = 0                     // 준비 완료된 샤드 수
	totalShards     = 1                     // 총 샤드 수 (초기값)
	shardMu         sync.Mutex              // 샤드 상태 Mutex
	shardGuildCount = make(map[int]int)     // 각 샤드별 서버 수 저장
	sessions        []*discordgo.Session    // 모든 샤드 세션 저장
	sessionsMu      sync.Mutex              // 세션 Mutex
)

// Ready 이벤트 핸들러
func Ready(s *discordgo.Session, r *discordgo.Ready) {
	shardMu.Lock()
	defer shardMu.Unlock()

	// 총 샤드 수를 갱신
	if s.Identify.Shard != nil && len(s.Identify.Shard) == 2 {
		totalShards = s.Identify.Shard[1]
		currentShard := s.Identify.Shard[0]
		// 현재 샤드의 서버 수 저장
		shardGuildCount[currentShard] = len(s.State.Guilds)
	}

	// 준비된 샤드 수 증가
	shardReadyCount++

	// 세션 저장
	sessionsMu.Lock()
	sessions = append(sessions, s)
	sessionsMu.Unlock()

	// 초기 길드 상태 저장
	initialGuildsMu.Lock()
	for _, guild := range r.Guilds {
		initialGuilds[guild.ID] = true
	}
	initialGuildsMu.Unlock()

	// 모든 샤드의 총 서버 수 계산
	totalGuilds := 0
	for _, count := range shardGuildCount {
		totalGuilds += count
	}

	// 모든 샤드가 준비된 경우 출력
	if shardReadyCount == totalShards {
		fmt.Printf("%s 로 로그인되었습니다. (총 샤드: %d개, 총 서버: %d개)\n",
			s.State.User.Username, totalShards, totalGuilds)
	}
}

// OnGuildJoin 이벤트 핸들러
func OnGuildJoin(s *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	initialGuildsMu.Lock()
	defer initialGuildsMu.Unlock()

	// 초기 길드 목록에 존재하는 경우라면 처리하지 않음 (이미 있는 길드)
	if _, exists := initialGuilds[guildCreate.Guild.ID]; exists {
		delete(initialGuilds, guildCreate.Guild.ID)
		return
	}

	// 새로운 길드에 초대된 경우 해당 샤드의 서버 수 증가
	shardMu.Lock()
	if s.Identify.Shard != nil {
		currentShard := s.Identify.Shard[0]
		shardGuildCount[currentShard]++
	}
	shardMu.Unlock()

	// 새로운 길드에 초대된 경우 처리
	logGuildJoin(s, guildCreate)
}

// 새로운 길드에 초대되었을 때 로그를 남기는 함수
func logGuildJoin(s *discordgo.Session, guildCreate *discordgo.GuildCreate) {
	owner := guildCreate.Guild.OwnerID
	ownerUser, err := s.User(owner)
	if err != nil {
		log.Println("오너 정보를 가져오는 데 오류 발생:", err)
		return
	}

	logMessage := fmt.Sprintf("%s(%s)에 봇 초대됨. owner = %s(%s) (서버 수: %d개)\n",
		guildCreate.Guild.Name, guildCreate.Guild.ID, ownerUser.Username, ownerUser.ID, len(s.State.Guilds))

	log.Print(logMessage)
}

// 길드에서 봇이 제거되었을 때 로그를 남기는 함수
func OnGuildLeave(s *discordgo.Session, guildDelete *discordgo.GuildDelete) {
	// 서버가 오프라인이 된 경우는 무시
	if guildDelete.Unavailable {
		return
	}

	// 해당 샤드의 서버 수 감소
	shardMu.Lock()
	if s.Identify.Shard != nil {
		currentShard := s.Identify.Shard[0]
		if shardGuildCount[currentShard] > 0 {
			shardGuildCount[currentShard]--
			fmt.Printf("샤드 %d의 서버 수 감소: %d\n", currentShard, shardGuildCount[currentShard]) // 디버그 로그 추가
		}
	}
	shardMu.Unlock()

	// 로그 메시지 출력
	fmt.Printf("서버에서 추방됨: %s (ID: %s) (현재 총 서버 수: %d)\n",
		guildDelete.ID,
		guildDelete.ID,
		len(s.State.Guilds))
}

// GetTotalGuilds returns the total number of guilds across all shards
func GetTotalGuilds() int {
    shardMu.Lock()
    defer shardMu.Unlock()
    
    totalGuilds := 0
    for _, count := range shardGuildCount {
        totalGuilds += count
    }
    return totalGuilds
}