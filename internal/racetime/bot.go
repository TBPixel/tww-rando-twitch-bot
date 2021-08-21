package racetime

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
)

type Bot struct {
	token  TokenSet
	config config.Racetime
}

func NewBot(c config.Racetime) (*Bot, error) {
	resp, err := http.PostForm(fmt.Sprintf("%s/o/token", c.URL), url.Values{
		"client_id":     []string{c.ClientID},
		"client_secret": []string{c.ClientSecret},
		"grant_type":    []string{"client_credentials"},
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token TokenSet
	err = json.NewDecoder(resp.Body).Decode(&token)

	return &Bot{
		token:  token,
		config: c,
	}, nil
}
