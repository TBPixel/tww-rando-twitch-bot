package racetime

import (
	"bufio"
	"io"
	"strings"
	"unicode"
)

const (
	// Tokens
	ERROR = iota
	ILLEGAL
	EOF
	IDENT
	// Keywords
	PREFIX
	RACE
	VS
	LEADERBOARD
	LINK
	EXAMPLE_PERMA
	PERMA
	RESTREAM
	MULTI
)

type Token int

type Ident struct {
	Token Token
	Lit   string
}

var tokens = []string{
	// tokens
	ERROR:   "ERROR",
	ILLEGAL: "ILLEGAL",
	EOF:     "EOF",
	IDENT:   "IDENT",
	// keywords
	PREFIX:        "!twwr",
	RACE:          "race",
	VS:            "vs",
	LEADERBOARD:   "leaderboard",
	LINK:          "link",
	EXAMPLE_PERMA: "exampleperma",
	PERMA:         "perma",
	RESTREAM:      "restream",
	MULTI:         "multi",
}

func (t Token) String() string {
	return tokens[t]
}

type Lexer struct {
	reader *bufio.Reader
}

func NewLexer(reader io.Reader) *Lexer {
	return &Lexer{
		reader: bufio.NewReader(reader),
	}
}

// Lex scans the input for the next token. It returns the position of the token,
// the token's type, and the literal value.
func (l *Lexer) Lex() (Token, string, error) {
	// keep looping until we return a token
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				return EOF, "", io.EOF
			}

			return ERROR, "", err
		}

		if unicode.IsSpace(r) {
			continue // nothing to do here, just move on
		}

		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) {
			// backup and let lexIdent rescan the beginning of the ident
			err := l.backup()
			if err != nil {
				return ERROR, "", nil
			}

			lit, err := l.lexIdent()
			if err != nil {
				return ERROR, "", err
			}

			switch strings.ToLower(lit) {
			case "!twwr":
				return PREFIX, lit, nil
			case "race":
				return RACE, lit, nil
			case "vs":
				return VS, lit, nil
			case "leaderboard":
				return LEADERBOARD, lit, nil
			case "link":
				return LINK, lit, nil
			case "exampleperma":
				return EXAMPLE_PERMA, lit, nil
			case "perma":
				return PERMA, lit, nil
			case "restream":
				return RESTREAM, lit, nil
			case "multi":
				return MULTI, lit, nil
			default:
				return IDENT, lit, nil
			}
		}

		return ILLEGAL, string(r), nil
	}
}

// LexAll scans the input for all tokens, returning an array when io.EOF is reached.
// if an unexpected error occurs, that is returned instead
func (l *Lexer) LexAll() ([]Ident, error) {
	var idents []Ident

	for {
		token, lit, err := l.Lex()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		idents = append(idents, Ident{
			Token: token,
			Lit:   lit,
		})
	}

	return idents, nil
}

func (l *Lexer) backup() error {
	return l.reader.UnreadRune()
}

// lexIdent scans the input until the end of an identifier and then returns the
// literal.
func (l *Lexer) lexIdent() (string, error) {
	var lit string
	for {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			if err == io.EOF {
				// at the end of the identifier
				return lit, nil
			}

			return lit, err
		}

		if unicode.IsLetter(r) || unicode.IsNumber(r) || unicode.IsPunct(r) {
			lit = lit + string(r)
		} else {
			// scanned something not in the identifier
			err := l.backup()
			if err != nil {
				return lit, err
			}

			return lit, nil
		}
	}
}
