package racetime_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/racetime"
)

type ErrReader struct {
	Err error
}

func (er ErrReader) Read(p []byte) (n int, err error) {
	return 0, er.Err
}

func TestLex(t *testing.T) {
	t.Run("should return an io.EOF at end", func(t *testing.T) {
		lex := racetime.NewLexer(strings.NewReader(""))

		_, _, err := lex.Lex()
		if err != io.EOF {
			t.Errorf("got %v, want %v", err, io.EOF)
		}
	})

	t.Run("should return an err if some unexpected error occurs", func(t *testing.T) {
		want := fmt.Errorf("a test error occurred")
		lex := racetime.NewLexer(ErrReader{Err: want})

		_, _, got := lex.Lex()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should skip white space", func(t *testing.T) {
		lex := racetime.NewLexer(strings.NewReader("  "))

		_, _, err := lex.Lex()
		if err != io.EOF {
			t.Errorf("got %v, want %v", err, io.EOF)
		}
	})

	t.Run("should return a number if found while parsing", func(t *testing.T) {
		want := "12345"
		lex := racetime.NewLexer(strings.NewReader(want))

		token, got, _ := lex.Lex()
		if token != racetime.IDENT {
			t.Errorf("got %v, want %v", token, racetime.IDENT)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("will return an IDENT if an unknown identifier is found", func(t *testing.T) {
		want := "key"
		lex := racetime.NewLexer(strings.NewReader(want))

		token, got, _ := lex.Lex()
		if token != racetime.IDENT {
			t.Errorf("got %v, want %v", token, racetime.IDENT)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("will return ident if an unknown word number combo is found", func(t *testing.T) {
		want := "s4"
		lex := racetime.NewLexer(strings.NewReader(want))

		token, got, _ := lex.Lex()
		if token != racetime.IDENT {
			t.Errorf("got %v, want %v", token, racetime.IDENT)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("will return a keyword if a known identifier is found", func(t *testing.T) {
		keywords := []racetime.Ident{
			{
				Token: racetime.PREFIX,
				Lit:   "!twwr",
			},
			{
				Token: racetime.RACE,
				Lit:   "race",
			},
			{
				Token: racetime.VS,
				Lit:   "vs",
			},
			{
				Token: racetime.LEADERBOARD,
				Lit:   "leaderboard",
			},
			{
				Token: racetime.LINK,
				Lit:   "link",
			},
			{
				Token: racetime.EXAMPLE_PERMA,
				Lit:   "exampleperma",
			},
			{
				Token: racetime.PERMA,
				Lit:   "perma",
			},
			{
				Token: racetime.RESTREAM,
				Lit:   "restream",
			},
			{
				Token: racetime.MULTI,
				Lit:   "multi",
			},
		}

		for _, tok := range keywords {
			want := tok.Lit
			lex := racetime.NewLexer(strings.NewReader(want))

			token, got, _ := lex.Lex()
			if token != tok.Token {
				t.Errorf("got %v, want %v", token, tok.Token)
			}

			if got != want {
				t.Errorf("got %v, want %v", got, want)
			}
		}
	})

	t.Run("will return ILLEGAL if a non letter, number or punctuation is found", func(t *testing.T) {
		lex := racetime.NewLexer(strings.NewReader("+"))

		token, _, _ := lex.Lex()
		if token != racetime.ILLEGAL {
			t.Errorf("got %v, want %v", token, racetime.ILLEGAL)
		}
	})
}

func TestLexAll(t *testing.T) {
	t.Run("will return a slice of racetime.Idents in the order they were found", func(t *testing.T) {
		want := []racetime.Ident{
			{
				Token: racetime.PREFIX,
				Lit:   "!twwr",
			},
			{
				Token: racetime.RACE,
				Lit:   "race",
			},
			{
				Token: racetime.VS,
				Lit:   "vs",
			},
			{
				Token: racetime.IDENT,
				Lit:   "someident5",
			},
			{
				Token: racetime.IDENT,
				Lit:   "5someident",
			},
			{
				Token: racetime.IDENT,
				Lit:   "12345",
			},
		}
		lex := racetime.NewLexer(strings.NewReader("!twwr race vs someident5 5someident 12345"))

		got, _ := lex.LexAll()
		for i, ident := range got {
			if i > len(want) {
				t.Errorf("got %v identifiers, want %v", len(got), len(want))
			}

			if want[i] != ident {
				t.Errorf("got %v at position %v, want %v", ident, i, want[i])
			}
		}
	})
}
