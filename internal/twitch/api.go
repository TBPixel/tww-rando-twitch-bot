package twitch

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
)

const (
	api         = "https://api.twitch.tv/helix"
	auth        = "https://id.twitch.tv/oauth2/token"
	clientGrant = "client_credentials"
)

type User struct {
	ID              string    `json:"id"`
	Login           string    `json:"login"`
	DisplayName     string    `json:"display_name"`
	ProfileImageURL string    `json:"profile_image_url"`
	CreatedAt       time.Time `json:"created_at"`
}

type ApiClient struct {
	config      config.Twitch
	client      http.Client
	accessToken string
}

func NewApiClient(config config.Twitch) (*ApiClient, error) {
	client := http.Client{}

	accessToken, err := newAccessToken(config)
	if err != nil {
		return nil, err
	}

	return &ApiClient{
		config,
		client,
		accessToken,
	}, nil
}

func (c *ApiClient) GetUsers(channels []string) ([]User, error) {
	req, err := c.req("GET", "users")
	if err != nil {
		return nil, err
	}

	query := req.URL.Query()
	query.Set("login", strings.Join(channels, ","))
	req.URL.RawQuery = query.Encode()

	body, err := c.fetch(req)
	if err != nil {
		return nil, err
	}
	defer body.Close()

	type getUserResponse struct {
		Data []User `json:"data"`
	}
	var payload getUserResponse
	err = json.NewDecoder(body).Decode(&payload)
	if err != nil {
		return nil, err
	}

	return payload.Data, nil
}

func (c *ApiClient) GetUser(channel string) (*User, error) {
	users, err := c.GetUsers([]string{channel})
	if err != nil {
		return nil, err
	}

	if len(users) == 0 {
		return nil, fmt.Errorf("no channel with the login %s could be found", channel)
	}

	return &users[0], nil
}

func (c *ApiClient) req(method, path string) (*http.Request, error) {
	uri := fmt.Sprintf("%s/%s", api, path)

	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Client-Id", c.config.ClientID)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.accessToken))

	return req, nil
}

func (c *ApiClient) fetch(req *http.Request) (io.ReadCloser, error) {
	res, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == http.StatusUnauthorized {
		accessToken, err := newAccessToken(c.config)
		if err != nil {
			return nil, err
		}

		c.accessToken = accessToken
		return c.fetch(req)
	}

	if res.StatusCode != http.StatusOK {
		return res.Body, fmt.Errorf("status code %v from twitch api", res.StatusCode)
	}

	return res.Body, nil
}

func newAccessToken(config config.Twitch) (string, error) {
	query := url.Values{
		"client_id":     []string{config.ClientID},
		"client_secret": []string{config.ClientSecret},
		"grant_type":    []string{clientGrant},
	}
	u, _ := url.Parse(auth)
	u.RawQuery = query.Encode()

	res, err := http.Post(u.String(), "", nil)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	type token struct {
		AccessToken  string   `json:"access_token"`
		RefreshToken string   `json:"refresh_token"`
		ExpiresIn    int      `json:"expires_in"`
		Scope        []string `json:"scope"`
		TokenType    string   `json:"token_type"`
	}
	var t token
	err = json.NewDecoder(res.Body).Decode(&t)
	if err != nil {
		return "", err
	}

	return t.AccessToken, nil
}
