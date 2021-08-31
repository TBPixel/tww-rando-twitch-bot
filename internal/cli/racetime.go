package cli

import (
	"fmt"
	"log"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/app"
	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"
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
