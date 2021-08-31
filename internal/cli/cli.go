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
						Action:      twitchLogin(app),
					},
					{
						Name:        "follow",
						Description: "Instruct twwr bot to follow a channel when listening",
						Usage:       "Specify the channel you want the bot to join chat for",
						ArgsUsage:   "account_id channel",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "disable",
								Aliases: []string{"d"},
								Usage:   "instruct the bot to stop listening on your twitch channel",
								Value:   false,
							},
						},
						Action: twitchFollow(app),
					},
					{
						Name:        "listen",
						Description: "Listen for all active twitch channels",
						Action:      twitchListen(app),
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
						Action:      racetimeConnect(app),
					},
					{
						Name:        "chat",
						Description: "watch the chat of the given race",
						ArgsUsage:   "race",
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
