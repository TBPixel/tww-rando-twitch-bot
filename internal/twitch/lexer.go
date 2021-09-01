package twitch

import "github.com/TBPixel/tww-rando-twitch-bot/internal/lexer"

const (
	PLAY = iota + lexer.Keyword
	PREFIX
	SETTINGS
	RACE
	VS
	LEADERBOARD
	LINK
	ExamplePerma
	PERMA
	RESTREAM
	MULTI
	HELP
)

var Keywords = []lexer.Ident{
	{
		Token: PLAY,
		Lit:   "!play",
	},
	{
		Token: PREFIX,
		Lit:   "!twwr",
	},
	{
		Token: SETTINGS,
		Lit:   "settings",
	},
	{
		Token: RACE,
		Lit:   "race",
	},
	{
		Token: VS,
		Lit:   "vs",
	},
	{
		Token: LEADERBOARD,
		Lit:   "leaderboard",
	},
	{
		Token: LINK,
		Lit:   "link",
	},
	{
		Token: ExamplePerma,
		Lit:   "exampleperma",
	},
	{
		Token: PERMA,
		Lit:   "perma",
	},
	{
		Token: RESTREAM,
		Lit:   "restream",
	},
	{
		Token: MULTI,
		Lit:   "multi",
	},
	{
		Token: HELP,
		Lit:   "help",
	},
}
