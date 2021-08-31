package app

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/twitch"
	"github.com/joho/godotenv"
)

type App struct {
	TwitchClient *twitch.ApiClient
	DB           *storage.DB
	Bot          *twitch.Bot
	Config       config.App
}

func New() (*App, error) {
	// opt to ignore dotenv loader as it is a development convenience
	_ = godotenv.Load()

	conf := config.Open()

	db, err := storage.Open(conf.DB)
	if err != nil {
		return nil, err
	}

	bot := twitch.NewBot(conf.Twitch, db)
	ttvClient, err := twitch.NewApiClient(conf.Twitch)
	if err != nil {
		return nil, err
	}

	return &App{
		TwitchClient: ttvClient,
		DB:           db,
		Bot:          bot,
		Config:       conf,
	}, nil
}
