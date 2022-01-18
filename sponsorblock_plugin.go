package sponsorblock

import (
	"encoding/json"
	"github.com/DisgoOrg/disgolink/lavalink"
)

var _ lavalink.EventExtensions = (*Plugin)(nil)

func New() *Plugin {
	return &Plugin{
		eventExtensions: []lavalink.EventExtension{
			SegmentsLoadedHandler{},
			SegmentSkippedHandler{},
		},
	}
}

type Plugin struct {
	eventExtensions []lavalink.EventExtension
}

func (p *Plugin) EventExtensions() []lavalink.EventExtension {
	return p.eventExtensions
}

var _ lavalink.EventExtension = (*SegmentsLoadedHandler)(nil)

type SegmentsLoadedHandler struct{}

func (h SegmentsLoadedHandler) Event() lavalink.EventType {
	return "SegmentsLoaded"
}
func (h SegmentsLoadedHandler) OnEvent(node lavalink.Node, data []byte) {
	var e SegmentsLoadedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		node.Lavalink().Logger().Error("Failed to unmarshal SegmentsLoaded Event", err)
		return
	}

	player := node.Lavalink().ExistingPlayer(e.GuildID)
	player.EmitEvent(func(l interface{}) {
		if listener, ok := l.(SegmentEventListener); ok {
			listener.OnSegmentsLoaded(player, e.Segments)
		}
	})
}

var _ lavalink.EventExtension = (*SegmentSkippedHandler)(nil)

type SegmentSkippedHandler struct{}

func (h SegmentSkippedHandler) Event() lavalink.EventType {
	return "SegmentSkipped"
}
func (h SegmentSkippedHandler) OnEvent(node lavalink.Node, data []byte) {
	var e SegmentSkippedEvent
	if err := json.Unmarshal(data, &e); err != nil {
		node.Lavalink().Logger().Error("Failed to unmarshal SegmentSkipped Event", err)
		return
	}

	player := node.Lavalink().ExistingPlayer(e.GuildID)
	player.EmitEvent(func(l interface{}) {
		if listener, ok := l.(SegmentEventListener); ok {
			listener.OnSegmentSkipped(player, e.Segment)
		}
	})
}
