package guilds

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/zetto21/discordgo/utils"
	"github.com/zetto21/discordgo/commands/slashCommands"
)



// Register slash commands in a map
var SlashCommands = map[string]func(*discordgo.Session, *discordgo.InteractionCreate){
	"ping": slashCommands.PingSlashCommand,
}


// InteractionCreate handles slash commands
func InteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	if i.Member.User.Bot {
		return
	}
	if i.Type != discordgo.InteractionApplicationCommand {
		return
	}

	commandName := i.ApplicationCommandData().Name
	guildName := utils.GetGuildName(s, i.GuildID)
	user := i.Member.User

	// 옵션 정보 수집
	var optionsInfo string
	if i.ApplicationCommandData().Options != nil {
		options := make([]string, 0)

		// 서브커맨드 처리
		for _, opt := range i.ApplicationCommandData().Options {
			if opt.Type == discordgo.ApplicationCommandOptionSubCommand {
				subOptions := make([]string, 0)
				for _, subOpt := range opt.Options {
					var optValue string
					switch subOpt.Type {
					case discordgo.ApplicationCommandOptionUser:
						if user := subOpt.UserValue(s); user != nil {
							optValue = fmt.Sprintf("<@%s>", user.ID)
						}
					case discordgo.ApplicationCommandOptionInteger:
						optValue = fmt.Sprint(subOpt.IntValue())
					case discordgo.ApplicationCommandOptionString:
						optValue = subOpt.StringValue()
					default:
						optValue = fmt.Sprint(subOpt.Value)
					}
					subOptions = append(subOptions, fmt.Sprintf("%s: %v", subOpt.Name, optValue))
				}
				options = append(options, fmt.Sprintf("%s [%s]", opt.Name, strings.Join(subOptions, ", ")))
			} else {
				// 일반 옵션 처리 (기존 코드)
				var optValue string
				switch opt.Type {
				case discordgo.ApplicationCommandOptionUser:
					if user := opt.UserValue(s); user != nil {
						optValue = fmt.Sprintf("<@%s>", user.ID)
					}
				default:
					optValue = fmt.Sprint(opt.Value)
				}
				options = append(options, fmt.Sprintf("%s: %v", opt.Name, optValue))
			}
		}
		if len(options) > 0 {
			optionsInfo = strings.Join(options, ", ")
		}
	}

	if command, exists := SlashCommands[commandName]; exists {
		command(s, i)
		log.Printf("슬래시 명령 사용됨: | 명령어: %s | 옵션: %s | 유저: %s (%s) | 길드: %s (%s)",
			commandName,
			optionsInfo,
			user.Username,
			user.ID,
			guildName,
			i.GuildID,
		)
	} else {
		log.Printf("\033[31m알 수 없는 슬래시 명령: | 명령어: %s | 옵션: %s | 유저: %s (%s) | 길드: %s (%s)\033[0m",
			commandName,
			optionsInfo,
			user.Username,
			user.ID,
			guildName,
			i.GuildID,
		)
	}
}
