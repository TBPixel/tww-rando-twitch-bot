package cli

import (
	"fmt"
	"log"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/twitch"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"

	"github.com/urfave/cli/v2"
)

func twitchLogin(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		authorizeUser(
			ctx,
			twitch.AuthURL,
			twitch.TokenURL,
			app.Config.Twitch.ClientID,
			app.Config.Twitch.ClientSecret,
			app.Config.Twitch.RedirectURL)
		//log.Println(ctx.Context.Value("token"))
		return nil
	}
}

func twitchChat(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		channel := ctx.Args().First()
		if channel == "" {
			return fmt.Errorf("twitch channel name is required")
		}
		bot := twitch.NewBot(app.Config.Twitch, app.DB)
		bot.Join(channel)

		log.Printf("connecting to irc and watching %s\n", channel)

		err := bot.Connect()
		if err != nil {
			return err
		}
		defer bot.Disconnect()

		return nil
	}
}
