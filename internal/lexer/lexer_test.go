package lexer_test

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/TBPixel/tww-rando-twitch-bot/internal/lexer"
)

var keywords = []lexer.Ident{
	{
		Token: lexer.Keyword,
		Lit:   "!twwr",
	},
	{
		Token: lexer.Keyword + 1,
		Lit:   "example",
	},
}

type ErrReader struct {
	Err error
}

func (er ErrReader) Read(p []byte) (n int, err error) {
	return 0, er.Err
}

func TestLex(t *testing.T) {
	t.Run("should return an error if the first ident has a value less than keyword", func(t *testing.T) {
		_, err := lexer.New(strings.NewReader(""), []lexer.Ident{
			{
				Token: lexer.Keyword - 1,
				Lit:   "",
			},
		})

		want := fmt.Errorf("first ident token iota of %d < %d, did you forget to do `Keyword = iota + lexer.Keyword`", lexer.Keyword-1, lexer.Keyword)
		if err.Error() != want.Error() {
			t.Errorf("got '%v', want '%v'", err, want)
		}
	})

	t.Run("should return an io.EOF at end", func(t *testing.T) {
		lex, _ := lexer.New(strings.NewReader(""), keywords)

		_, _, err := lex.Lex()
		if err != io.EOF {
			t.Errorf("got %v, want %v", err, io.EOF)
		}
	})

	t.Run("should return an err if some unexpected error occurs", func(t *testing.T) {
		want := fmt.Errorf("a test error occurred")
		lex, _ := lexer.New(ErrReader{Err: want}, keywords)

		_, _, got := lex.Lex()
		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("should skip white space", func(t *testing.T) {
		lex, _ := lexer.New(strings.NewReader("  "), keywords)

		_, _, err := lex.Lex()
		if err != io.EOF {
			t.Errorf("got %v, want %v", err, io.EOF)
		}
	})

	t.Run("should return a number if found while parsing", func(t *testing.T) {
		want := "12345"
		lex, _ := lexer.New(strings.NewReader(want), keywords)

		token, got, _ := lex.Lex()
		if token != lexer.IDENT {
			t.Errorf("got %v, want %v", token, lexer.IDENT)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("will return an IDENT if an unknown identifier is found", func(t *testing.T) {
		want := "key"
		lex, _ := lexer.New(strings.NewReader(want), keywords)

		token, got, _ := lex.Lex()
		if token != lexer.IDENT {
			t.Errorf("got %v, want %v", token, lexer.IDENT)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("will return ident if an unknown word number combo is found", func(t *testing.T) {
		want := "s4"
		lex, _ := lexer.New(strings.NewReader(want), keywords)

		token, got, _ := lex.Lex()
		if token != lexer.IDENT {
			t.Errorf("got %v, want %v", token, lexer.IDENT)
		}

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("will return a keyword if a known identifier is found", func(t *testing.T) {
		for _, tok := range keywords {
			want := tok.Lit
			lex, _ := lexer.New(strings.NewReader(want), keywords)

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
		lex, _ := lexer.New(strings.NewReader("+"), keywords)

		token, _, _ := lex.Lex()
		if token != lexer.ILLEGAL {
			t.Errorf("got %v, want %v", token, lexer.ILLEGAL)
		}
	})
}

func TestLexAll(t *testing.T) {
	t.Run("will return a slice of twitch.idents in the order they were found", func(t *testing.T) {
		want := []lexer.Ident{
			{
				Token: lexer.Keyword,
				Lit:   "!twwr",
			},
			{
				Token: lexer.Keyword + 1,
				Lit:   "example",
			},
			{
				Token: lexer.IDENT,
				Lit:   "someident5",
			},
			{
				Token: lexer.IDENT,
				Lit:   "5someident",
			},
			{
				Token: lexer.IDENT,
				Lit:   "12345",
			},
		}
		lex, _ := lexer.New(strings.NewReader("!twwr example someident5 5someident 12345"), keywords)

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
