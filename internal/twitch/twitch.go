package twitch

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"

	"github.com/gempir/go-twitch-irc/v2"
)

const (
	URL      = "https://twitch.tv"
	TokenURL = "https://id.twitch.tv/oauth2/token"
	AuthURL  = "https://id.twitch.tv/oauth2/authorize"
)

// Client
type Client struct {
	client *twitch.Client
}

// NewClient creates a client connected to the twitch IRC server
func NewClient(conf config.Twitch) *Client {
	client := twitch.NewClient(conf.Username, conf.OAuth)

	return &Client{client}
}

// Connect to the twitch irc server
func (c *Client) Connect() error {
	return c.client.Connect()
}

// Disconnect from the twitch IRC server
func (c *Client) Disconnect() error {
	return c.client.Disconnect()
}
