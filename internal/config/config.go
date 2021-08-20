package config

import (
	"os"
	"strings"
)

// Open
func Open() App {
	return App{
		Host: os.Getenv("APP_HOST"),
		DB: DB{
			Path:          os.Getenv("DB_PATH"),
			EnableLogging: false,
		},
		Twitch: Twitch{
			Username:     os.Getenv("TWITCH_USERNAME"),
			OAuth:        os.Getenv("TWITCH_OAUTH"),
			ClientID:     os.Getenv("TWITCH_CLIENT_ID"),
			ClientSecret: os.Getenv("TWITCH_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("TWITCH_REDIRECT_URL"),
		},
		Racetime: newRacetime(),
	}
}

// App
type App struct {
	Host     string
	DB       DB
	Twitch   Twitch
	Racetime Racetime
}

// DB
type DB struct {
	Path          string
	EnableLogging bool
}

// Twitch
type Twitch struct {
	Username     string
	OAuth        string
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

type Racetime struct {
	URL          string
	WSSchema     string
	ClientID     string
	ClientSecret string
}

func newRacetime() Racetime {
	wsSchema := "wss"
	if strings.Contains(os.Getenv("RACETIME_URL"), "local") {
		wsSchema = "ws"
	}
	return Racetime{
		URL:          os.Getenv("RACETIME_URL"),
		WSSchema:     wsSchema,
		ClientID:     os.Getenv("RACETIME_CLIENT_ID"),
		ClientSecret: os.Getenv("RACETIME_CLIENT_SECRET"),
	}
}
