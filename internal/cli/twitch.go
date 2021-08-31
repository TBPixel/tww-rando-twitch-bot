package cli

import (
	"fmt"
	"log"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"

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
		token, ok := ctx.Context.Value("token").(TwitchAccessTokenContents)
		if !ok {
			return fmt.Errorf("failed to get token from twitch api")
		}
		ttvUser, err := app.TwitchClient.GetUser(token.PreferredUsername)
		if err != nil {
			return err
		}

		user, err := app.DB.FindUser(storage.UserQuery{
			Field: storage.FieldTwitchID,
			Value: ttvUser.ID,
		})
		if err != nil {
			if err != storage.ErrNotFound {
				return err
			}

			user, err = app.DB.CreateUser(ttvUser.ID, ttvUser.Login, ttvUser.DisplayName, ttvUser.ProfileImageURL)
			if err != nil {
				return err
			}
			log.Printf("new user created with id %d for ttv %s", user.ID, user.TwitchName)
		}

		log.Printf("%+v\n", user)

		return nil
	}
}

func twitchChat(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		channel := ctx.Args().First()
		if channel == "" {
			return fmt.Errorf("twitch channel name is required")
		}
		app.Bot.Join(channel)

		log.Printf("connecting to irc and watching %s\n", channel)

		app.Bot.Listen(ctx.Context)

		return nil
	}
}
