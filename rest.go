package sponsorblock

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/disgoorg/disgolink/v3/disgolink"
	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/json"
	"github.com/disgoorg/snowflake/v2"
)

// GetCategories gets the categories to skip for the guild.
func GetCategories(client disgolink.RestClient, sessionID string, guildID snowflake.ID) ([]SegmentCategory, error) {
	rq, err := http.NewRequest(http.MethodGet, fmt.Sprintf("/v4/sessions/%s/players/%s/sponsorblock/categories", sessionID, guildID), nil)
	if err != nil {
		return nil, err
	}
	rs, err := client.Do(rq)
	if err != nil {
		return nil, err
	}

	defer rs.Body.Close()

	if rs.StatusCode == http.StatusNotFound {
		return nil, nil
	} else if rs.StatusCode != http.StatusOK {
		var lavalinkError lavalink.Error
		if err = json.NewDecoder(rs.Body).Decode(&lavalinkError); err != nil {
			return nil, err
		}
		return nil, lavalinkError
	}

	var result []SegmentCategory
	if err = json.NewDecoder(rs.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// SetCategories sets the categories to skip for the guild.
func SetCategories(client disgolink.RestClient, sessionID string, guildID snowflake.ID, categories []SegmentCategory) error {
	buff := new(bytes.Buffer)
	if err := json.NewEncoder(buff).Encode(categories); err != nil {
		return err
	}
	rq, err := http.NewRequest(http.MethodPut, fmt.Sprintf("/v4/sessions/%s/players/%s/sponsorblock/categories", sessionID, guildID), buff)
	if err != nil {
		return err
	}

	rs, err := client.Do(rq)
	if err != nil {
		return err
	}

	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		var lavalinkError lavalink.Error
		if err = json.NewDecoder(rs.Body).Decode(&lavalinkError); err != nil {
			return err
		}
		return lavalinkError
	}

	return nil
}

// DeleteCategories deletes the categories to skip for the guild.
func DeleteCategories(client disgolink.RestClient, sessionID string, guildID snowflake.ID) error {
	rq, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("/v4/sessions/%s/players/%s/sponsorblock/categories", sessionID, guildID), nil)
	if err != nil {
		return err
	}

	rs, err := client.Do(rq)
	if err != nil {
		return err
	}

	defer rs.Body.Close()

	if rs.StatusCode != http.StatusOK {
		var lavalinkError lavalink.Error
		if err = json.NewDecoder(rs.Body).Decode(&lavalinkError); err != nil {
			return err
		}
		return lavalinkError
	}

	return nil
}
