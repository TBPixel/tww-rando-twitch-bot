package twitch

import "github.com/TBPixel/tww-rando-twitch-bot/internal/lexer"

const (
	PREFIX = iota + lexer.Keyword
	RACE
	VS
	LEADERBOARD
	LINK
	ExamplePerma
	PERMA
	RESTREAM
	MULTI
)

var Keywords = []lexer.Ident{
	{
		Token: PREFIX,
		Lit:   "!twwr",
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
}
