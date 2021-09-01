package cli

import (
	"fmt"
	"log"
	"strings"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/races"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"
	"github.com/urfave/cli/v2"
)

func twwrBot(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		category := ctx.Args().First()
		if category == "" {
			return fmt.Errorf("missing required argument: category")
		}

		monitor := races.NewMonitor(app.Config.Racetime, category)
		listener := monitor.AddListener()
		log.Printf("racetime monitor watching all races in %s", category)
		go monitor.Listen(ctx.Context)

		users, err := app.DB.FindUsers(storage.UserQuery{
			Field: storage.FieldActiveInChannel,
			Value: true,
		})
		if err != nil {
			return err
		}
		var channels []string
		for _, u := range users {
			channels = append(channels, u.TwitchName)
		}

		app.Bot.Join(channels...)

		log.Printf("bot listening to all active twitch channels: %s", strings.Join(channels, ", "))
		app.Bot.Listen(ctx.Context, listener)

		return nil
	}
}
