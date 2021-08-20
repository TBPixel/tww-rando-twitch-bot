package twitch

import (
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
)

type TokenSet struct {
	AccessToken  string        `json:"access_token"`
	RefreshToken string        `json:"refresh_token"`
	ExpiresIn    int           `json:"expires_in"`
	Scope        []interface{} `json:"scope"`
	TokenType    string        `json:"token_type"`
}

func OAuthClientConnection(c config.Twitch) (*TokenSet, error) {
	resp, err := http.PostForm(TokenURL, url.Values{
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
		"grant_type":    []string{"client_credentials"},
		"scope":         []string{},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var t TokenSet
	err = json.NewDecoder(resp.Body).Decode(&t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
