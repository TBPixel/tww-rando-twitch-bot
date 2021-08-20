package cli

import (
	"fmt"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"

	"github.com/urfave/cli/v2"
)

func NewApp(app app.App) *cli.App {
	return &cli.App{
		Name:  "twwr",
		Usage: "Wind Waker Randomizer Twitch bot from the command line",
		Commands: []*cli.Command{
			{
				Name:        "twitch",
				Description: "twitch specific commands",
				Subcommands: []*cli.Command{
					{
						Name:        "login",
						Description: "Authorize twwr bot to act on your twitch channel",
						Usage:       "Redirect to an ouath2 access request",
						Action:      twitchLogin(app),
					},
				},
			},
			{
				Name:        "racetime",
				Description: "Racetime.gg specific commands",
				Subcommands: []*cli.Command{
					{
						Name:        "connect",
						Description: "Link your racetime account",
						Usage:       "Redirect to an oauth2 access request",
						Action:      racetimeConnect(app),
					},
					{
						Name:        "chat",
						Description: "watch the chat of the given race",
						Action: func(context *cli.Context) error {
							name := context.Args().First()
							if name == "" {
								return fmt.Errorf("race name is required")
							}

							bot, err := racetime.NewBot(app.Config.Racetime)
							if err != nil {
								return err
							}
							err = bot.Connect(context.Context, name)
							if err != nil {
								return err
							}
							return nil
						},
					},
				},
			},
		},
	}
}
