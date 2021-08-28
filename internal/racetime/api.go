package racetime

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"

	"github.com/fatih/structs"
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
	Goal      string `json:"goal"`
	NumRanked string `json:"num_ranked"`
	Rankings  struct {
		User         string `json:"user"`
		Place        string `json:"place"`
		PlaceOrdinal string `json:"place_ordinal"`
		Score        int    `json:"score"`
		TimesRaced   string `json:"times_raced"`
	}
}

type UserSearchParameters struct {
	Name string
	Scrim string
	Term string
}

// CategoryDetail fetches the race data of a specific racetime.gg category
func CategoryDetail(c config.Racetime, category string) (*CategoryResponse, error) {
	content, err := Get(c, fmt.Sprintf("%s/data", category), nil)

	var cr *CategoryResponse
	json.Unmarshal([]byte(content), cr)
	if err != nil {
		return nil, err
	}

	return cr, nil
}

func CategoryLeaderboards(c config.Racetime, category string) (*LeaderboardsResponse, error) {
	content, err := Get(c, fmt.Sprintf("%s/leaderboards/data", category), nil)

	var cl *LeaderboardsResponse
	json.Unmarshal([]byte(content), cl)
	if err != nil {
		return nil, err
	}

	return cl, nil
}

// RaceDetail fetches the race data of a specific racetime.gg race
func RaceDetail(c config.Racetime, category string, race string) (*RaceData, error) {
	content, err := Get(c, fmt.Sprintf("%s/%s/data", category, race), nil)

	var rd *RaceData
	json.Unmarshal([]byte(content), rd)
	if err != nil {
		return nil, err
	}

	return rd, nil
}

func PastUserRaces(c config.Racetime, user string, showEntrants bool, page int){
	content, err := Get(c, fmt.Sprintf("user/%s/races/data", user), nil)

	var pur *RaceData
	json.Unmarshal([]byte(content), pur)
	if err != nil {
	}
}

func UserSearch(c config.Racetime, user string){
	parameters := UserSearchParameters{
		Name: "colfra",
	}

	m := structs.Map(parameters)
	content, _ := Get(c, "user/search", m)
	var userSearchResults UserDataResponse
	json.Unmarshal([]byte(content), &userSearchResults)
	log.Println(userSearchResults.Results[0].ID)
}

func GetTest(){
	var test = `{"results": [{"id": "5BRGVMd30E368Lzv", "full_name": "colfra", "name": "colfra", "discriminator": null, "url": "/user/5BRGVMd30E368Lzv", "avatar": null, "pronouns": null, "flair": "staff", "twitch_name": null, "twitch_display_name": null, "twitch_channel": null, "can_moderate": false}]}`;
	jd := json.NewDecoder(strings.NewReader(test))
	var userResults UserDataResponse
	jd.Decode(&userResults)
	log.Println(userResults.Results[0].ID)
}

func Get(c config.Racetime, endpoint string, parameters map[string]interface{}) (string, error) {
	requestUrl, err := url.Parse(fmt.Sprintf("%s/%s", c.URL, endpoint))
	if (len(parameters) != 0){
		q := requestUrl.Query()
		for key, value := range parameters {
			if (value != ""){
				q.Add(strings.ToLower(key), value.(string))
			}
		}
		requestUrl.RawQuery = q.Encode()
	}

	resp, err := http.Get(requestUrl.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    content := string(bodyBytes)

	return content, nil
}