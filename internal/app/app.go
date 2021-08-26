package app

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/twitch"
	"github.com/joho/godotenv"
)

type App struct {
	DB     *storage.DB
	Bot    *twitch.Bot
	Config config.App
}

func Run() (*App, error) {
	// opt to ignore dotenv loader as it is a development convenience
	_ = godotenv.Load()

	conf := config.Open()

	db, err := storage.Open(conf.DB)
	if err != nil {
		return nil, err
	}

	bot := twitch.NewBot(conf.Twitch, db)

	return &App{
		DB:     db,
		Bot:    bot,
		Config: conf,
	}, nil
}
