package sponsorblock

import (
	"encoding/json"
	"log/slog"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
)

var (
	_ disgolink.EventPlugins = (*Plugin)(nil)
	_ disgolink.Plugin       = (*Plugin)(nil)
)

func New() *Plugin {
	return NewWithLogger(slog.Default())
}

func NewWithLogger(logger *slog.Logger) *Plugin {
	return &Plugin{
		eventPlugins: []disgolink.EventPlugin{
			&segmentsLoadedHandler{
				logger: logger,
			},
			&segmentSkippedHandler{
				logger: logger,
			},
			&chaptersLoadedHandler{
				logger: logger,
			},
			&chapterStartedHandler{
				logger: logger,
			},
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

type segmentsLoadedHandler struct {
	logger *slog.Logger
}

func (h *segmentsLoadedHandler) Event() lavalink.EventType {
	return EventTypeSegmentsLoaded
}
func (h *segmentsLoadedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e SegmentsLoadedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Error("Failed to unmarshal SegmentsLoaded Event", slog.Any("err", err))
		return
	}

	player.Lavalink().EmitEvent(player, e)
}

var _ disgolink.EventPlugin = (*segmentSkippedHandler)(nil)

type segmentSkippedHandler struct {
	logger *slog.Logger
}

func (h *segmentSkippedHandler) Event() lavalink.EventType {
	return EventTypeSegmentSkipped
}

func (h *segmentSkippedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e SegmentSkippedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Error("Failed to unmarshal SegmentSkipped Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}

type chaptersLoadedHandler struct {
	logger *slog.Logger
}

func (h *chaptersLoadedHandler) Event() lavalink.EventType {
	return EventTypeChaptersLoaded
}

func (h *chaptersLoadedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e ChaptersLoadedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Error("Failed to unmarshal ChaptersLoaded Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}

var _ disgolink.EventPlugin = (*chapterStartedHandler)(nil)

type chapterStartedHandler struct {
	logger *slog.Logger
}

func (h *chapterStartedHandler) Event() lavalink.EventType {
	return EventTypeChapterStarted
}

func (h *chapterStartedHandler) OnEventInvocation(player disgolink.Player, data []byte) {
	var e ChapterStartedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		h.logger.Error("Failed to unmarshal ChapterStarted Event", err)
		return
	}

	player.Lavalink().EmitEvent(player, e)
}
