package sponsorblock

import (
	"github.com/DisgoOrg/disgolink/lavalink"
)

var _ lavalink.OpCommand = (*PlayCommand)(nil)

type PlayCommand struct {
	lavalink.PlayCommand
	SkipSegments []SegmentCategory `json:"skip_segments"`
}
