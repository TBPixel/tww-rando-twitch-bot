package lexer

import (
	"bufio"
	"fmt"
	"io"
	"strings"
	"unicode"
)

const (
	ERROR = iota
	ILLEGAL
	EOF
	IDENT
	Keyword // iota of keywords. lexer users should append this to their iota as Keyword = iota + lexer.Keyword
)

// Token represents a literary Token for query purposes
type Token int

// Ident packages an identifier iota and it's literary string representation
type Ident struct {
	Token Token
	Lit   string
}

// idents groups a list of idents together
type idents []Ident

// Find searches a list of idents for the literary string and returns
// the first match found, or nil otherwise
func (idents idents) Find(lit string) *Ident {
	for _, id := range idents {
		if id.Lit != lit {
			continue
		}

		return &id
	}

	return nil
}

// Lexer provides a simple lexing iteration over a given reader.
type Lexer struct {
	reader *bufio.Reader
	idents idents
	tokens []string
}

// New returns a new lexer which reads from the given io.Reader
// and searches for the list of idents
//
// lex := lexer.New(reader, keywords)
// Token, lit, err := lex.Lex()
// if err != nil && err != io.EOF {
//   return err
// }
func New(reader io.Reader, idents []Ident) (*Lexer, error) {
	if len(idents) > 0 && idents[0].Token < Keyword {
		return nil, fmt.Errorf("first ident Token iota of %d < %d, did you forget to do `Keyword = iota + lexer.Keyword`", idents[0].Token, Keyword)
	}

	return &Lexer{
		reader: bufio.NewReader(reader),
		idents: idents,
	}, nil
}

// Lex scans the input for the next Token. It returns the position of the Token,
// the Token's type, and the literal value.
func (l *Lexer) Lex() (Token, string, error) {
	// keep looping until we return a Token
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

			ident := l.idents.Find(strings.ToLower(lit))
			if ident == nil {
				return IDENT, lit, nil
			}

			return ident.Token, lit, nil
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

// backup unreads the most recently read rune
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
