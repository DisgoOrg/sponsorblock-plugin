package sponsorblock

import (
	"github.com/DisgoOrg/disgolink/lavalink"
	"time"
)

type SegmentsLoadedEvent struct {
	GuildID  string    `json:"guild_id"`
	Segments []Segment `json:"segments"`
}

type SegmentSkippedEvent struct {
	GuildID string  `json:"guild_id"`
	Segment Segment `json:"segment"`
}

type Segment struct {
	Category SegmentCategory `json:"category"`
	Start    time.Duration   `json:"start"`
	End      time.Duration   `json:"end"`
}

type SegmentCategory string

const (
	SegmentCategorySponsor       = "sponsor"
	SegmentCategorySelfpromo     = "selfpromo"
	SegmentCategoryInteraction   = "interaction"
	SegmentCategoryIntro         = "intro"
	SegmentCategoryOutro         = "outro"
	SegmentCategoryPreview       = "preview"
	SegmentCategoryMusicOfftopic = "music_offtopic"
	SegmentCategoryFiller        = "filler"
)

type SegmentEventListener interface {
	OnSegmentsLoaded(player lavalink.Player, segments []Segment)
	OnSegmentSkipped(player lavalink.Player, segment Segment)
}

type SegmentEventAdapter struct{}

func (a SegmentEventAdapter) OnSegmentsLoaded(player lavalink.Player, segments []Segment) {}
func (a SegmentEventAdapter) OnSegmentSkipped(player lavalink.Player, segment Segment)    {}
