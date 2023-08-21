package sponsorblock

import (
	"encoding/json"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var (
	_ disgolink.EventPlugins = (*Plugin)(nil)
	_ disgolink.Plugin       = (*Plugin)(nil)
)

func New() *Plugin {
	return &Plugin{
		eventPlugins: []disgolink.EventPlugin{
			&segmentsLoadedHandler{},
			&segmentSkippedHandler{},
			&chaptersLoadedHandler{},
			&chapterStartedHandler{},
		},
	}
}

type Plugin struct {
	eventPlugins []disgolink.EventPlugin
}

func (p *Plugin) EventPlugins() []disgolink.EventPlugin {
	return p.eventPlugins
}

func (p *Plugin) Name() string {
	return "sponsorblock"
}

func (p *Plugin) Version() string {
	return "1.0.0"
}

var _ disgolink.EventPlugin = (*segmentsLoadedHandler)(nil)

type segmentsLoadedHandler struct{}

func (h *segmentsLoadedHandler) Event() lavalink.EventType {
	return EventTypeSegmentsLoaded
}
func (h *segmentsLoadedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e SegmentsLoadedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		player.Lavalink().Logger().Error("Failed to unmarshal SegmentsLoaded Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}

var _ disgolink.EventPlugin = (*segmentSkippedHandler)(nil)

type segmentSkippedHandler struct{}

func (h *segmentSkippedHandler) Event() lavalink.EventType {
	return EventTypeSegmentSkipped
}

func (h *segmentSkippedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e SegmentSkippedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		player.Lavalink().Logger().Error("Failed to unmarshal SegmentSkipped Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}

type chaptersLoadedHandler struct{}

func (h *chaptersLoadedHandler) Event() lavalink.EventType {
	return EventTypeChaptersLoaded
}

func (h *chaptersLoadedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e ChaptersLoadedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		player.Lavalink().Logger().Error("Failed to unmarshal ChaptersLoaded Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}

var _ disgolink.EventPlugin = (*chapterStartedHandler)(nil)

type chapterStartedHandler struct{}

func (h *chapterStartedHandler) Event() lavalink.EventType {
	return EventTypeChapterStarted
}

func (h *chapterStartedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e ChapterStartedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		player.Lavalink().Logger().Error("Failed to unmarshal ChapterStarted Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}
