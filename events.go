package sponsorblock

import (
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"time"
)

const (
	EventTypeSegmentsLoaded = "SegmentsLoaded"
	EventTypeSegmentSkipped = "SegmentSkipped"
)

type SegmentsLoadedEvent struct {
	GuildID_ string    `json:"guild_id"`
	Segments []Segment `json:"segments"`
}

func (e SegmentsLoadedEvent) Type() lavalink.EventType {
	return EventTypeSegmentsLoaded
}

func (e SegmentsLoadedEvent) GuildID() snowflake.ID {
	return snowflake.MustParse(e.GuildID_)
}

type SegmentSkippedEvent struct {
	GuildID_ string  `json:"guild_id"`
	Segment  Segment `json:"segment"`
}

func (e SegmentSkippedEvent) Type() lavalink.EventType {
	return EventTypeSegmentSkipped
}

func (e SegmentSkippedEvent) GuildID() snowflake.ID {
	return snowflake.MustParse(e.GuildID_)
}

type Segment struct {
	Category SegmentCategory `json:"category"`
	Start    time.Duration   `json:"start"`
	End      time.Duration   `json:"end"`
}

type SegmentCategory string

//goland:noinspection GoUnusedConst
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
	OnSegmentsLoaded(player disgolink.Player, event SegmentsLoadedEvent)
	OnSegmentSkipped(player disgolink.Player, event SegmentSkippedEvent)
}
