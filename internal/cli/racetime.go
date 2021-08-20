package cli

import (
	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"
	"github.com/urfave/cli/v2"
)

func racetimeConnect(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		//authorizeUser(
		//	ctx,
		//	app.Config.Twitch.ClientID,
		//	app.Config.Twitch.ClientSecret,
		//	app.Config.Twitch.RedirectURL)
		//log.Println(ctx.Context.Value("token"))
		return nil
	}
}
