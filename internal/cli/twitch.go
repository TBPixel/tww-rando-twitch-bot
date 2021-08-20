package cli

import (
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
		log.Println(ctx.Context.Value("token"))
		return nil
	}
}
