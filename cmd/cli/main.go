package main

import (
	"log"
	"os"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/cli"
)

func main() {
	application, err := app.New()
	if err != nil {
		log.Fatalln(err)
	}

	err = cli.NewApp(*application).Run(os.Args)
	if err != nil {
		log.Fatalln(err)
	}
}
