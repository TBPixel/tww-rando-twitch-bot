package twitch

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/races"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"

	"github.com/pkg/errors"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/lexer"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/gempir/go-twitch-irc/v2"
)

var (
	ErrNotEnoughArgs = errors.New("not enough arguments to bot")
)

const (
	MultiTwitchURL = "https://multitwitch.tv"
	SeedHashPrefix = "Seed Hash:"
	Delimiter      = " | "
)

// Bot remains connected to twitch IRC, watches
// chats and shares messages received through a channel
type Bot struct {
	racetimeURL string
	db          *storage.DB
	client      *twitch.Client
	msgChan     <-chan twitch.PrivateMessage
	mut         sync.Mutex
	races       []racetime.RaceData
}

// NewBot creates a client connected to the twitch Bot server
func NewBot(conf config.Twitch, db *storage.DB, racetimeURL string) *Bot {
	client := twitch.NewClient(conf.Username, conf.IRCOAuth)
	msgChan := make(chan twitch.PrivateMessage)

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		msgChan <- message
	})

	return &Bot{
		racetimeURL: racetimeURL,
		db:          db,
		client:      client,
		msgChan:     msgChan,
		mut:         sync.Mutex{},
		races:       []racetime.RaceData{},
	}
}

// Join a twitch channel to listen for private messages
func (b *Bot) Join(channels ...string) {
	b.client.Join(channels...)
}

// Listen connects to the IRC server and awaits messages,
// handling any it sees as commands.
func (b *Bot) Listen(ctx context.Context, listener chan []racetime.RaceData) {
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
		case racesData := <-listener:
			b.mut.Lock()
			b.races = racesData
			b.mut.Unlock()
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

		log.Printf("error parsing bot command: %s", err)
		return nil
	}

	streamer, err := b.db.FindUser(storage.UserQuery{
		Field: storage.FieldTwitchID,
		Value: message.RoomID,
	})
	if err != nil {
		return err
	}
	if streamer == nil {
		return fmt.Errorf("unabled to find streamer with twitch id %s (name %s)", message.RoomID, message.Channel)
	}

	// ensure the bot can play marbles with Tanjo3 :widepeepoHappy:
	if idents[0].Token == PLAY && message.User.ID == streamer.TwitchID {
		b.client.Say(message.Channel, "!play")
		return nil
	}

	// skip if not a !twwr nor a recognized command
	if idents[0].Token != PREFIX {
		return nil
	}

	// TODO: Verify bot to allow whisper of help command
	if len(idents) == 1 {
		b.client.Say(message.Channel, handleHelpCommand())
		return nil
	}

	// Non-race commands
	switch idents[1].Token {
	case HELP:
		b.client.Say(message.Channel, handleHelpCommand())
		return nil
	case RESTREAM:
		// TODO: Restream command will check if the user has linked to a restream for this race
		// if one is set it will share that link
		//b.client.Say(message.Channel, "restream")
		return nil
	case LEADERBOARD:
		// TODO: Figure out what leaderboard will actually do
		//b.client.Say(message.Channel, "leaderboard")
		return nil
	}

	race := b.findRaceForUser(*streamer)
	if race == nil {
		b.client.Say(message.Channel, notCurrentlyInRace(*streamer))
		return nil
	}

	// race only commands
	switch idents[1].Token {
	case SETTINGS:
		b.client.Say(message.Channel, handleSettingsCommand(*streamer, *race))
	case RACE:
		b.client.Say(message.Channel, handleRaceCommand(*streamer, *race))
	case VS:
		b.client.Say(message.Channel, handleVsCommand(*streamer, *race))
	case LINK:
		b.client.Say(message.Channel, fmt.Sprintf("%s/%s", b.racetimeURL, race.Name))
	case ExamplePerma:
		b.client.Say(message.Channel, handleExamplePermaCommand(*streamer, *race))
	case MULTI:
		b.client.Say(message.Channel, handleMultiCommand(*streamer, *race))
	case PERMA:
		b.client.Say(message.Channel, handlePermaCommand(*streamer, *race))
	}

	return nil
}

func (b *Bot) findRaceForUser(user storage.User) (race *racetime.RaceData) {
	b.mut.Lock()
	defer b.mut.Unlock()

	for _, r := range b.races {
		for _, entrant := range r.Entrants {
			if entrant.User.ID == user.RacetimeID {
				return &r
			}
		}
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

	if len(args) == 0 {
		return nil, ErrNotEnoughArgs
	}

	return args, nil
}

func extractPreset(race racetime.RaceData) string {
	for _, p := range races.Presets() {
		if strings.Contains(strings.ToLower(race.Info), strings.ToLower(p)) {
			return p
		}
	}

	return ""
}

func notCurrentlyInRace(streamer storage.User) string {
	return fmt.Sprintf("%s is not currently in a race", streamer.TwitchDisplayName)
}

func customCategory(streamer storage.User) string {
	return fmt.Sprintf("%s is playing a custom race category", streamer.TwitchDisplayName)
}

func handleSettingsCommand(streamer storage.User, race racetime.RaceData) string {
	if race.Goal.Name != races.Standard && race.Goal.Name != races.SpoilerLog {
		return customCategory(streamer)
	}

	ex := races.ExamplePermaByPreset(extractPreset(race))
	if ex == nil {
		return customCategory(streamer)
	}

	return fmt.Sprintf("%s: %s", ex.Preset, ex.Description)
}

func handleRaceCommand(streamer storage.User, race racetime.RaceData) string {
	if race.Goal.Name != races.Standard && race.Goal.Name != races.SpoilerLog {
		return customCategory(streamer)
	}

	ex := races.ExamplePermaByPreset(extractPreset(race))
	if ex == nil {
		return customCategory(streamer)
	}

	return fmt.Sprintf("%s is playing %s (!twwr settings)", streamer.TwitchDisplayName, ex.Preset)
}

func handleExamplePermaCommand(streamer storage.User, race racetime.RaceData) string {
	if race.Goal.Name != races.Standard && race.Goal.Name != races.SpoilerLog {
		return customCategory(streamer)
	}

	ex := races.ExamplePermaByPreset(extractPreset(race))
	if ex == nil {
		return customCategory(streamer)
	}

	return fmt.Sprintf("example permalink: %s", ex.Perma)
}

func handleVsCommand(streamer storage.User, race racetime.RaceData) string {
	var entrants []string
	for _, u := range race.Entrants {
		// skip the streamer
		if u.User.ID == streamer.RacetimeID {
			continue
		}

		if u.User.TwitchName != "" {
			entrants = append(entrants, u.User.TwitchDisplayName)
		} else {
			entrants = append(entrants, u.User.Name)
		}
	}

	if len(entrants) == 0 {
		return fmt.Sprintf("There are currently no other entrants in race with %s", streamer.TwitchDisplayName)
	}

	return fmt.Sprintf("%s is currently racing against: %s", streamer.TwitchDisplayName, strings.Join(entrants, ", "))
}

func handleMultiCommand(streamer storage.User, race racetime.RaceData) string {
	var entrants []string
	for _, u := range race.Entrants {
		// skip users without a twitch account
		if u.User.TwitchName == "" {
			continue
		}

		entrants = append(entrants, u.User.TwitchName)
	}

	if len(entrants) == 0 {
		return fmt.Sprintf("There are currently no other entrants in race with %s", streamer.TwitchDisplayName)
	}

	return fmt.Sprintf("%s/%s", MultiTwitchURL, strings.Join(entrants, "/"))
}

func handlePermaCommand(streamer storage.User, race racetime.RaceData) string {
	if race.Goal.Name != races.Standard && race.Goal.Name != races.SpoilerLog {
		return customCategory(streamer)
	}

	hashStartIndex := strings.Index(race.Info, SeedHashPrefix)
	if hashStartIndex == -1 {
		return "Permalink has not yet been generated or cannot be found"
	}

	seedEndIndex := hashStartIndex - len(Delimiter)
	seedStartIndex := strings.LastIndex(race.Info[:seedEndIndex], Delimiter) + len(Delimiter)
	if seedStartIndex == -1 {
		return "Permalink has not yet been generated or cannot be found"
	}

	return race.Info[seedStartIndex:seedEndIndex]
}

func handleHelpCommand() string {
	// list of active commands
	commands := []string{
		"settings",
		"race",
		"vs",
		"link",
		"exampleperma",
		"perma",
		"multi",
		"help",
	}

	return fmt.Sprintf("commands: %s", strings.Join(commands, ", "))
}
