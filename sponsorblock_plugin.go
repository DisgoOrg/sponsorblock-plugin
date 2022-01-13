package sponsorblock

import (
	"github.com/DisgoOrg/disgolink/lavalink"
)

var _ lavalink.Plugin = (*Plugin)(nil)

func New() *Plugin {
	return &Plugin{}
}

type Plugin struct {
}

func (p *Plugin) Name() string {
	return Name
}

func (p *Plugin) Version() string {
	return Version
}

func (p *Plugin) Event() lavalink.EventType {
	return "SegmentsLoaded"
}
func (p *Plugin) OnEventInvocation(node lavalink.Node, data []byte) {

}
