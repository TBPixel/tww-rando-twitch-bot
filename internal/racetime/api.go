package racetime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
)

type PaginatedRaces struct {
	Count    uint       `json:"count"`
	NumPages uint       `json:"num_pages"`
	Races    []RaceData `json:"races"`
}

type Entrant struct {
	User   UserData
	Status struct {
		Value        string `json:"value"`
		VerboseValue string `json:"verbose_value"`
		HelpText     string `json:"help_text"`
	} `json:"status"`
	FinishTime     string        `json:"finish_time"`
	FinishedAt     time.Time     `json:"finished_at"`
	Place          int           `json:"place"`
	PlaceOrdinal   string        `json:"place_ordinal"`
	Score          int           `json:"score"`
	ScoreChange    int           `json:"score_change"`
	Comment        interface{}   `json:"comment"`
	HasComment     bool          `json:"has_comment"`
	StreamLive     bool          `json:"stream_live"`
	StreamOverride bool          `json:"stream_override"`
	Actions        []interface{} `json:"actions"`
}

type RaceData struct {
	Name   string `json:"name"`
	Slug   string `json:"slug"`
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
	Entrants              []Entrant `json:"entrants"`
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

type UserDataResponse struct {
	Results []UserData `json:"results"`
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

type LeaderboardsResponse struct {
	Leaderboards []struct {
		Goal      string `json:"goal"`
		NumRanked int    `json:"num_ranked"`
		Rankings  []struct {
			User         UserData `json:"user"`
			Place        int      `json:"place"`
			PlaceOrdinal string   `json:"place_ordinal"`
			Score        int      `json:"score"`
			TimesRaced   int      `json:"times_raced"`
		}
	} `json:"leaderboards"`
}

type UserSearchParameters struct {
	Name  string
	Scrim string
	Term  string
}

// CategoryDetail fetches the race data of a specific racetime.gg category
func CategoryDetail(c config.Racetime, category string) (*CategoryResponse, error) {
	req, err := req(c, "GET", fmt.Sprintf("%s/data", category))
	if err != nil {
		return nil, err
	}

	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("bad request")
	}

	var cr CategoryResponse
	err = json.NewDecoder(res.Body).Decode(&cr)
	if err != nil {
		return nil, err
	}

	return &cr, nil
}

func CategoryLeaderboards(c config.Racetime, category string) (*LeaderboardsResponse, error) {
	req, err := req(c, "GET", fmt.Sprintf("%s/leaderboards/data", category))
	if err != nil {
		return nil, err
	}

	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("bad request")
	}

	var cl LeaderboardsResponse
	err = json.NewDecoder(res.Body).Decode(&cl)
	if err != nil {
		return nil, err
	}

	return &cl, nil
}

// RaceDetail fetches the race data of a specific racetime.gg race
func RaceDetail(c config.Racetime, category string, race string) (*RaceData, error) {
	req, err := req(c, "GET", fmt.Sprintf("%s/%s/data", category, race))
	if err != nil {
		return nil, err
	}

	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("bad request")
	}

	var rd RaceData
	err = json.NewDecoder(res.Body).Decode(&rd)
	if err != nil {
		return nil, err
	}

	return &rd, nil
}

func PastUserRaces(c config.Racetime, user string, showEntrants bool, page uint) (*PaginatedRaces, error) {
	req, err := req(c, "GET", fmt.Sprintf("user/%s/races/data", user))
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("page", strconv.Itoa(int(page)))
	if showEntrants {
		query.Set("show_entrants", "true")
	}

	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("bad request")
	}

	var pur PaginatedRaces
	err = json.NewDecoder(res.Body).Decode(&pur)
	if err != nil {
		return nil, err
	}

	return &pur, nil
}

func UserSearch(c config.Racetime, name string) ([]UserData, error) {
	req, err := req(c, "GET", "user/search")
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("name", name)
	req.URL.RawQuery = query.Encode()

	res, err := fetch(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("bad request")
	}

	type userSearchResults struct {
		Results []UserData `json:"results"`
	}

	var results userSearchResults
	err = json.NewDecoder(res.Body).Decode(&results)
	if err != nil {
		return nil, err
	}

	return results.Results, nil
}

func fetch(req *http.Request) (*http.Response, error) {
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func req(c config.Racetime, method, path string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/%s", c.URL, path)

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	return req, nil
}
