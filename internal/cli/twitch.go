package cli

import (
	"fmt"
	"log"
	"strconv"
	"strings"

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

func twitchFollow(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		enabled := !ctx.Bool("disable")
		idStr := ctx.Args().First()
		if idStr == "" {
			return fmt.Errorf("missing required argument: id")
		}

		id, err := strconv.Atoi(idStr)
		if err != nil {
			return fmt.Errorf("id must be an unsigned integer")
		}

		user, err := app.DB.UpdateUser(uint64(id), storage.UserUpdate{
			ActiveInChannel: &enabled,
		})
		if err != nil {
			return err
		}

		log.Printf("twitch bot follow mode set to %v for channel %s", enabled, user.TwitchName)

		return nil
	}
}

func twitchListen(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
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
		app.Bot.Listen(ctx.Context)

		return nil
	}
}
