package config

import (
	"os"
)

// Open
func Open() App {
	return App{
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
	}
}

// App
type App struct {
	DB     DB
	Twitch Twitch
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
