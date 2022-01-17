package sponsorblock

import (
	"encoding/json"
	"github.com/DisgoOrg/disgolink/lavalink"
)

var (
	_ lavalink.EventExtensions = (*Plugin)(nil)
)

func New() *Plugin {
	plugin := &Plugin{}
	plugin.eventExtensions = []lavalink.EventExtension{
		&SegmentsLoadedHandler{Plugin: plugin},
		&SegmentSkippedHandler{Plugin: plugin},
	}
	return plugin
}

type Plugin struct {
	eventExtensions []lavalink.EventExtension
}

func (p *Plugin) EventExtensions() []lavalink.EventExtension {
	return p.eventExtensions
}

type SegmentsLoadedHandler struct {
	*Plugin
}

func (h *SegmentsLoadedHandler) Event() lavalink.EventType {
	return "SegmentsLoaded"
}
func (h *SegmentsLoadedHandler) OnEventInvocation(node lavalink.Node, data []byte) {
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

type SegmentSkippedHandler struct {
	*Plugin
}

func (h *SegmentSkippedHandler) Event() lavalink.EventType {
	return "SegmentSkipped"
}
func (h *SegmentSkippedHandler) OnEventInvocation(node lavalink.Node, data []byte) {
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
