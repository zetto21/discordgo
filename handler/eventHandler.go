// handler/eventHandler.go
package handler

import (
	"github.com/bwmarrin/discordgo"
	"github.com/zetto21/discordgo/events/client"
	"github.com/zetto21/discordgo/events/guilds"
)

func RegisterEventHandlers(dg *discordgo.Session) {
	dg.AddHandler(client.Ready)
	dg.AddHandler(guilds.InteractionCreate)
	dg.AddHandler(guilds.MessageCreate)
}
