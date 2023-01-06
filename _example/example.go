package main

import (
	"context"
	"fmt"
	"github.com/disgoorg/disgo"
	"github.com/disgoorg/disgo/bot"
	"github.com/disgoorg/disgo/cache"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/gateway"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"github.com/disgoorg/sponsorblock-plugin"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/disgoorg/log"
)

var (
	Token   = os.Getenv("TOKEN")
	GuildID = snowflake.GetEnv("GUILD_ID")

	NodeName      = os.Getenv("NODE_NAME")
	NodeAddress   = os.Getenv("NODE_ADDRESS")
	NodePassword  = os.Getenv("NODE_PASSWORD")
	NodeSecure, _ = strconv.ParseBool(os.Getenv("NODE_SECURE"))
)

func main() {
	log.SetLevel(log.LevelInfo)
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	client, err := disgo.New(Token,
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

	sponsorblockPlugin := sponsorblock.New()
	lavalinkClient := disgolink.New(client.ApplicationID(),
		disgolink.WithPlugins(sponsorblockPlugin),
		disgolink.WithListenerFunc(func(player disgolink.Player, event sponsorblock.SegmentSkippedEvent) {
			fmt.Printf("skipped segment: %v\n", event)
		}),
		disgolink.WithListenerFunc(func(player disgolink.Player, event sponsorblock.SegmentsLoadedEvent) {
			fmt.Printf("loaded segments: %v\n", event)
		}),
	)
	client.AddEventListeners(
		bot.NewListenerFunc(func(event *events.GuildVoiceStateUpdate) {
			if event.Member.User.ID != client.ID() {
				return
			}
			lavalinkClient.OnVoiceStateUpdate(context.TODO(), event.VoiceState.GuildID, event.VoiceState.ChannelID, event.VoiceState.SessionID)
		}),
		bot.NewListenerFunc(func(event *events.VoiceServerUpdate) {
			lavalinkClient.OnVoiceServerUpdate(context.TODO(), event.GuildID, event.Token, *event.Endpoint)
		}),
		bot.NewListenerFunc(func(event *events.ApplicationCommandInteractionCreate) {
			data := event.SlashCommandInteractionData()
			if data.CommandName() == "play" {

				lavalinkClient.BestNode().LoadTracks(context.TODO(), data.String("query"), disgolink.NewResultHandler(
					func(track lavalink.Track) {
						voiceState, ok := client.Caches().VoiceStates().Get(*event.GuildID(), event.User().ID)
						if !ok {
							_ = event.CreateMessage(discord.MessageCreate{
								Content: "you need to be in a voice channel",
							})
							return
						}
						_ = client.Connect(context.TODO(), *event.GuildID(), *voiceState.ChannelID)

						player := lavalinkClient.Player(*event.GuildID())
						if err = sponsorblockPlugin.SetCategories(context.TODO(), player.Node(), *event.GuildID(), []string{"sponsor", "selfpromo", "interaction", "intro", "outro", "preview", "music_offtopic", "filler"}); err != nil {
							log.Error("error setting categories: ", err)
						}
						if err = player.Update(context.TODO(), lavalink.WithTrack(track)); err != nil {
							log.Error("error updating player: ", err)
						}
						if err = event.CreateMessage(discord.MessageCreate{
							Content: "playing: " + track.Info.Title,
						}); err != nil {
							log.Error("error creating message: ", err)
						}
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
			} else if data.CommandName() == "seek" {
				lavalinkClient.Player(*event.GuildID()).Update(context.TODO(), lavalink.WithPosition(lavalink.Duration(data.Int("position"))*lavalink.Second))
			}
		}),
	)

	if _, err = client.Rest().SetGuildCommands(client.ApplicationID(), GuildID, []discord.ApplicationCommandCreate{
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
		discord.SlashCommandCreate{
			Name:        "seek",
			Description: "seeks music",
			Options: []discord.ApplicationCommandOption{
				discord.ApplicationCommandOptionInt{
					Name:        "position",
					Description: "position in seconds",
					Required:    true,
				},
			},
		},
	}); err != nil {
		panic(err)
		return
	}

	if err = client.OpenGateway(context.TODO()); err != nil {
		panic(err)
	}

	defer client.Close(context.TODO())

	_, err = lavalinkClient.AddNode(context.TODO(), disgolink.NodeConfig{
		Name:     NodeName,
		Address:  NodeAddress,
		Password: NodePassword,
		Secure:   NodeSecure,
	})

	log.Info("Example is now running. Press CTRL-C to exit.")
	s := make(chan os.Signal, 1)
	signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-s
}
