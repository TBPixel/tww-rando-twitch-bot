package twitch

import (
	"log"
	"strings"

	"github.com/pkg/errors"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/lexer"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/gempir/go-twitch-irc/v2"
)

var (
	ErrNotEnoughArgs = errors.New("not enough arguments to bot")
)

// Bot
type Bot struct {
	client *twitch.Client
}

// NewBot creates a client connected to the twitch IRC server
func NewBot(conf config.Twitch, db *storage.DB) *Bot {
	client := twitch.NewClient(conf.Username, conf.OAuth)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		idents, err := parseBotCommands(message.Message)
		if err != nil {
			if err == ErrNotEnoughArgs {
				return
			}

			log.Printf("error parsing bot commonds: %s", err)
			return
		}

		switch idents[1].Token {
		case RACE:
			client.Say(message.Channel, "race")
		case VS:
			client.Say(message.Channel, "vs")
		case LEADERBOARD:
			client.Say(message.Channel, "leaderboard")
		case LINK:
			client.Say(message.Channel, "link")
		case ExamplePerma:
			client.Say(message.Channel, "exampleperma")
		case PERMA:
			client.Say(message.Channel, "perma")
		case RESTREAM:
			client.Say(message.Channel, "restream")
		case MULTI:
			client.Say(message.Channel, "multi")
		default:
			return
		}
	})

	return &Bot{client}
}

// Connect to the twitch irc server
func (b *Bot) Connect() error {
	return b.client.Connect()
}

// Join an IRC for a twitch channel
func (b *Bot) Join(channels ...string) {
	b.client.Join(channels...)
}

// Disconnect from the twitch IRC server
func (b *Bot) Disconnect() error {
	return b.client.Disconnect()
}

func parseBotCommands(input string) ([]lexer.Ident, error) {
	lex, err := lexer.New(strings.NewReader(input), Keywords)
	if err != nil {
		return nil, err
	}

	args, err := lex.LexAll()
	if err != nil {
		return nil, err
	}

	if len(args) < 2 {
		return nil, ErrNotEnoughArgs
	}

	// skip if not a !twwr nor a recognized command
	if args[0].Token != PREFIX {
		return nil, ErrNotEnoughArgs
	}

	return args, nil
}
