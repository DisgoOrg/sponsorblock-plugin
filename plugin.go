package sponsorblock

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/disgoorg/disgolink/v2/disgolink"
	"github.com/disgoorg/disgolink/v2/lavalink"
	"github.com/disgoorg/snowflake/v2"
	"net/http"
	"time"
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
		},
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

type Plugin struct {
	eventPlugins []disgolink.EventPlugin
	httpClient   *http.Client
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

func (p *Plugin) GetCategories(ctx context.Context, node disgolink.Node, guildID snowflake.ID) (categories []string, err error) {
	path := fmt.Sprintf("%s/v3/sessions/%s/players/%d/sponsorblock/categories", node.Config().RestURL(), node.SessionID(), guildID)

	err = p.do(node, ctx, http.MethodGet, path, nil, &categories)
	return
}

func (p *Plugin) SetCategories(ctx context.Context, node disgolink.Node, guildID snowflake.ID, categories []string) error {
	path := fmt.Sprintf("%s/v3/sessions/%s/players/%d/sponsorblock/categories", node.Config().RestURL(), node.SessionID(), guildID)

	return p.do(node, ctx, http.MethodPut, path, categories, nil)
}

func (p *Plugin) DeleteCategories(ctx context.Context, node disgolink.Node, guildID snowflake.ID) error {
	path := fmt.Sprintf("%s/v3/sessions/%s/players/%d/sponsorblock/categories", node.Config().RestURL(), node.SessionID(), guildID)

	return p.do(node, ctx, http.MethodDelete, path, nil, nil)
}

func (p *Plugin) do(node disgolink.Node, ctx context.Context, method string, path string, rqBody any, rsBody any) error {

	var rqBuff *bytes.Buffer
	if rqBody != nil {
		rqBuff = &bytes.Buffer{}
		if err := json.NewEncoder(rqBuff).Encode(rqBody); err != nil {
			return err
		}
	}

	rq, err := http.NewRequestWithContext(ctx, method, path, rqBuff)
	if err != nil {
		return err
	}
	rq.Header.Set("Authorization", node.Config().Password)
	if rqBody != nil {
		rq.Header.Set("Content-Type", "application/json")
	}

	rs, err := p.httpClient.Do(rq)
	if err != nil {
		return err
	}
	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to set categories: %s", rs.Status)
	}

	if rsBody != nil {
		if err = json.NewDecoder(rs.Body).Decode(rsBody); err != nil {
			return err
		}
	}
	return nil
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
