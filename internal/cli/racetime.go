package cli

import (
	"fmt"
	"log"
	"strconv"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/storage"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"
	"github.com/urfave/cli/v2"
)

func racetimeConnect(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		// This is to workaround being in a cli context.
		// In an actual app we'd keep the user authenticated
		// and link them via a session
		idStr := ctx.Args().First()
		if idStr == "" {
			return fmt.Errorf("missing required argument: account_id")
		}
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return err
		}

		user, err := app.DB.FindUser(storage.UserQuery{
			Field: storage.FieldID,
			Value: uint64(id),
		})
		if err != nil {
			return err
		}

		authorizeUser(
			ctx,
			app.Config,
			fmt.Sprintf("%s/%s", app.Config.Racetime.URL, racetime.AuthURL),
			fmt.Sprintf("%s/%s", app.Config.Racetime.URL, racetime.TokenURL),
			app.Config.Racetime.ClientID,
			app.Config.Racetime.ClientSecret,
			app.Config.Racetime.RedirectURL,
			[]string{"read"},
			racetimeTokenParserFunc)

		tkn, ok := ctx.Context.Value("token").(RacetimeAccessTokenContents)
		if !ok {
			return fmt.Errorf("failed to parse racetime access token contents")
		}

		user, err = app.DB.UpdateUser(user.ID, storage.UserUpdate{
			RacetimeID: &tkn.ID,
		})

		log.Printf("%+v\n", user)

		return nil
	}
}

func racetimeUserSearch(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		name := ctx.Args().First()
		if name == "" {
			return fmt.Errorf("missing required argument: name")
		}

		user, err := racetime.UserSearch(app.Config.Racetime, name)
		if err != nil {
			return err
		}

		log.Printf("%+v\n", user)
		return nil
	}
}

func racetimePastUserRaces(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		showEntrants := ctx.Bool("show_entrances")
		page := ctx.Uint("page")

		id := ctx.Args().First()
		if id == "" {
			return fmt.Errorf("missing required argument: id")
		}

		races, err := racetime.PastUserRaces(app.Config.Racetime, id, showEntrants, page)
		if err != nil {
			return err
		}

		log.Printf("%+v\n", races)

		return nil
	}
}

func racetimeCategoryDetail(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		category := ctx.Args().First()
		if category == "" {
			return fmt.Errorf("missing required argument: category")
		}

		detail, err := racetime.CategoryDetail(app.Config.Racetime, category)
		if err != nil {
			return err
		}

		log.Printf("%+v\n", detail)

		return nil
	}
}

func racetimeCategoryLeaderboards(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		category := ctx.Args().First()
		if category == "" {
			return fmt.Errorf("missing required argument: category")
		}

		leaderboards, err := racetime.CategoryLeaderboards(app.Config.Racetime, category)
		if err != nil {
			return err
		}

		log.Printf("%+v\n", leaderboards)

		return nil
	}
}

func racetimeRaceDetail(app app.App) cli.ActionFunc {
	return func(ctx *cli.Context) error {
		if ctx.NArg() < 2 {
			return fmt.Errorf("missing required arguments: category race")
		}

		category, race := ctx.Args().Get(0), ctx.Args().Get(1)

		raceDetail, err := racetime.RaceDetail(app.Config.Racetime, category, race)
		if err != nil {
			return err
		}

		log.Printf("%+v\n", raceDetail)

		return nil
	}
}
