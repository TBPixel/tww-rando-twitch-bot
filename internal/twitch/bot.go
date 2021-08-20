package twitch

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/gempir/go-twitch-irc/v2"
)

// Bot
type Bot struct {
	client *twitch.Client
}

// NewBot creates a client connected to the twitch IRC server
func NewBot(conf config.Twitch) *Bot {
	client := twitch.NewClient(conf.Username, conf.OAuth)

	return &Bot{client}
}

// Connect to the twitch irc server
func (c *Bot) Connect() error {
	return c.client.Connect()
}

// Join an IRC for a twitch channel
func (c *Bot) Join(channels ...string) {
	c.client.Join(channels...)
}

// Disconnect from the twitch IRC server
func (c *Bot) Disconnect() error {
	return c.client.Disconnect()
}
