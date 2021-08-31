package twitch

import (
	"context"
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

// Bot remains connected to twitch IRC, watches
// chats and shares messages received through a channel
type Bot struct {
	client  *twitch.Client
	msgChan <-chan twitch.PrivateMessage
}

// NewBot creates a client connected to the twitch Bot server
func NewBot(conf config.Twitch, db *storage.DB) *Bot {
	client := twitch.NewClient(conf.Username, conf.IRCOAuth)
	msgChan := make(chan twitch.PrivateMessage)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msgChan <- message
	})

	return &Bot{client, msgChan}
}

// Join a twitch channel to listen for private messages
func (b *Bot) Join(channels ...string) {
	b.client.Join(channels...)
}

// Listen connects to the IRC server and awaits messages,
// handling any it sees as commands.
func (b *Bot) Listen(ctx context.Context) {
	go func() {
		err := b.client.Connect()
		if err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case <-ctx.Done():
			err := b.client.Disconnect()
			if err != nil {
				log.Println(err)
			}
			return
		case msg := <-b.msgChan:
			err := b.handleMessage(msg)
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func (b *Bot) handleMessage(message twitch.PrivateMessage) error {
	idents, err := parseBotCommands(message.Message)
	if err != nil {
		if err == ErrNotEnoughArgs {
			return nil
		}

		log.Printf("error parsing bot commonds: %s", err)
		return nil
	}

	switch idents[1].Token {
	case RACE:
		b.client.Say(message.Channel, "race")
	case VS:
		b.client.Say(message.Channel, "vs")
	case LEADERBOARD:
		b.client.Say(message.Channel, "leaderboard")
	case LINK:
		b.client.Say(message.Channel, "link")
	case ExamplePerma:
		b.client.Say(message.Channel, "exampleperma")
	case PERMA:
		b.client.Say(message.Channel, "perma")
	case RESTREAM:
		b.client.Say(message.Channel, "restream")
	case MULTI:
		b.client.Say(message.Channel, "multi")
	}

	return nil
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
