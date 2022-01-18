package main

import (
	"github.com/DisgoOrg/disgo/core/events"
	"github.com/DisgoOrg/disgo/discord"
	"github.com/DisgoOrg/disgolink/lavalink"
	"os"

	"github.com/DisgoOrg/disgo/core/bot"
	"github.com/DisgoOrg/disgolink/disgolink"
	"github.com/DisgoOrg/log"
)

var (
	token   = os.Getenv("bot_token")
	guildID = discord.Snowflake(os.Getenv("guild_id"))
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	disgo, err := bot.New(token)
	if err != nil {
		panic(err)
	}
	link := disgolink.New(disgo)

	if _, err = disgo.SetGuildCommands(guildID, []discord.ApplicationCommandCreate{
		discord.SlashCommandCreate{
			Name:        "play",
			Description: "plays music",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionString{
					Name:        "query",
					Description: "what to play",
					Required:    true,
				},
			},
			DefaultPermission: false,
		},
	}); err != nil {
		panic(err)
		return
	}

	disgo.AddEventListeners(&events.ListenerAdapter{
		OnApplicationCommandInteraction: func(event *events.ApplicationCommandInteractionEvent) {
			data := event.SlashCommandInteractionData()
			if data.CommandName != "play" {
				return
			}
			link.BestRestClient().LoadItemHandler(*data.Options.String("query"), lavalink.NewResultHandler(
				func(track lavalink.Track) {
					link.Player(*event.GuildID)
				},
				func(playlist lavalink.Playlist) {

				},
				func(tracks []lavalink.Track) {

				},
				func() {

				},
				func(e lavalink.Exception) {

				},
			))
		},
	})
}
