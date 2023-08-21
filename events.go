package sponsorblock

import (
	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

const (
	EventTypeSegmentsLoaded = "SegmentsLoaded"
	EventTypeSegmentSkipped = "SegmentSkipped"
	EventTypeChaptersLoaded = "ChaptersLoaded"
	EventTypeChapterStarted = "ChapterStarted"
)

type SegmentsLoadedEvent struct {
	GuildID_ snowflake.ID `json:"guild_id"`
	Segments []Segment    `json:"segments"`
}

func (SegmentsLoadedEvent) Op() lavalink.Op {
	return lavalink.OpEvent
}

func (SegmentsLoadedEvent) Type() lavalink.EventType {
	return EventTypeSegmentsLoaded
}

func (e SegmentsLoadedEvent) GuildID() snowflake.ID {
	return e.GuildID_
}

type SegmentSkippedEvent struct {
	GuildID_ snowflake.ID `json:"guild_id"`
	Segment  Segment      `json:"segment"`
}

func (SegmentSkippedEvent) Op() lavalink.Op {
	return lavalink.OpEvent
}

func (SegmentSkippedEvent) Type() lavalink.EventType {
	return EventTypeSegmentSkipped
}

func (e SegmentSkippedEvent) GuildID() snowflake.ID {
	return e.GuildID_
}

type Segment struct {
	Category SegmentCategory   `json:"category"`
	Start    lavalink.Duration `json:"start"`
	End      lavalink.Duration `json:"end"`
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
	OnSegmentsLoaded(player disgolink.Player, event SegmentsLoadedEvent)
	OnSegmentSkipped(player disgolink.Player, event SegmentSkippedEvent)
}

type ChaptersLoadedEvent struct {
	GuildID_ snowflake.ID `json:"guild_id"`
	Chapters []Chapter    `json:"chapters"`
}

func (ChaptersLoadedEvent) Op() lavalink.Op {
	return lavalink.OpEvent
}

func (ChaptersLoadedEvent) Type() lavalink.EventType {
	return EventTypeChaptersLoaded
}

func (e ChaptersLoadedEvent) GuildID() snowflake.ID {
	return e.GuildID_
}

type Chapter struct {
	Name     string            `json:"name"`
	Start    lavalink.Duration `json:"start"`
	End      lavalink.Duration `json:"end"`
	Duration lavalink.Duration `json:"duration"`
}

type ChapterStartedEvent struct {
	GuildID_ snowflake.ID `json:"guild_id"`
	Chapter  Chapter      `json:"chapter"`
}

func (ChapterStartedEvent) Op() lavalink.Op {
	return lavalink.OpEvent
}

func (ChapterStartedEvent) Type() lavalink.EventType {
	return EventTypeChapterStarted
}

func (e ChapterStartedEvent) GuildID() snowflake.ID {
	return e.GuildID_
}
