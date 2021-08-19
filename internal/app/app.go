package app

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/config"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/twitch"
	"github.com/joho/godotenv"
)

type App struct {
	DB     *storage.DB
	IRC    *twitch.Client
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

	irc := twitch.NewClient(conf.Twitch)
	go func() {
		err = irc.Connect()
	}()

	return &App{
		DB:     db,
		IRC:    irc,
		Config: conf,
	}, err
}
