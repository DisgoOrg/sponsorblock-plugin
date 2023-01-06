package main

import (
	"context"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"os"

	"github.com/disgoorg/log"
)

var (
	token   = os.Getenv("bot_token")
	guildID = snowflake.GetEnv("guild_id")
)

func main() {
	log.SetLevel(log.LevelDebug)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client, err := disgo.New(token,
		bot.WithGatewayConfigOpts(
			gateway.WithIntents(gateway.IntentGuilds, gateway.IntentGuildVoiceStates),
		),
		bot.WithCacheConfigOpts(
			cache.WithCacheFlags(cache.FlagVoiceStates),
		),
	)
	if err != nil {
		panic(err)
	}
	lavalinkClient := disgolink.New(client.ApplicationID(),
		disgolink.WithPlugins(sponsorblock.New()),
	)
	client.AddEventListeners(
		bot.NewListenerFunc(func(event *events.GuildVoiceStateUpdate) {
			lavalinkClient.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
		}),
		bot.NewListenerFunc(func(event *events.VoiceServerUpdate) {
			lavalinkClient.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
		}),
		bot.NewListenerFunc(func(event *events.ApplicationCommandInteractionCreate) {
			data := event.SlashCommandInteractionData()
			if data.CommandName() != "play" {
				return
			}
			lavalinkClient.BestNode().LoadTracks(context.TODO(), data.String("query"), disgolink.NewResultHandler(
				func(track lavalink.Track) {

					lavalinkClient.Player(*event.GuildID()).Update(context.TODO(), lavalink.WithTrack(track))
				},
				func(playlist lavalink.Playlist) {

				},
				func(tracks []lavalink.Track) {

				},
				func() {

				},
				func(e error) {

				},
			))
		}),
	)

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), guildID, []discord.ApplicationCommandCreate{
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
		},
	}); err != nil {
		panic(err)
		return
	}
}
