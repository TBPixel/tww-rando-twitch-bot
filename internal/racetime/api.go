package racetime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
)

type RaceData struct {
	Name   string `json:"name"`
	Status struct {
		Value        string `json:"value"`
		VerboseValue string `json:"verbose_value"`
		HelpText     string `json:"help_text"`
	} `json:"status"`
	URL     string `json:"url"`
	DataURL string `json:"data_url"`
	Goal    struct {
		Name   string `json:"name"`
		Custom bool   `json:"custom"`
	} `json:"goal"`
	Info                  string    `json:"info"`
	EntrantsCount         int       `json:"entrants_count"`
	EntrantsCountFinished int       `json:"entrants_count_finished"`
	EntrantsCountInactive int       `json:"entrants_count_inactive"`
	OpenedAt              time.Time `json:"opened_at"`
	StartedAt             time.Time `json:"started_at"`
	TimeLimit             string    `json:"time_limit"`
	Category              *struct {
		Name      string `json:"name"`
		ShortName string `json:"short_name"`
		Slug      string `json:"slug"`
		URL       string `json:"url"`
		DataURL   string `json:"data_url"`
		Image     string `json:"image"`
	} `json:"category,omitempty"`
}

type UserData struct {
	ID                string      `json:"id"`
	FullName          string      `json:"full_name"`
	Name              string      `json:"name"`
	Discriminator     interface{} `json:"discriminator"`
	URL               string      `json:"url"`
	Avatar            string      `json:"avatar"`
	Pronouns          string      `json:"pronouns"`
	Flair             string      `json:"flair"`
	TwitchName        string      `json:"twitch_name"`
	TwitchDisplayName string      `json:"twitch_display_name"`
	TwitchChannel     string      `json:"twitch_channel"`
	CanModerate       bool        `json:"can_moderate"`
}

type CategoryResponse struct {
	Name              string     `json:"name"`
	ShortName         string     `json:"short_name"`
	Slug              string     `json:"slug"`
	URL               string     `json:"url"`
	DataURL           string     `json:"data_url"`
	Image             string     `json:"image"`
	Info              string     `json:"info"`
	StreamingRequired bool       `json:"streaming_required"`
	Owners            []UserData `json:"owners"`
	Moderators        []UserData `json:"moderators"`
	Goals             []string   `json:"goals"`
	CurrentRaces      []RaceData `json:"current_races"`
}

// CategoryDetail fetches the race data of a specific racetime.gg category
func CategoryDetail(c config.Racetime, category string) (*CategoryResponse, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s/data", c.URL, category))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var cr *CategoryResponse
	err = json.NewDecoder(resp.Body).Decode(cr)
	if err != nil {
		return nil, err
	}

	return cr, nil
}
