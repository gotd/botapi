package pool

import (
	"strconv"
	"strings"

	"github.com/go-faster/errors"
)

// Token is a parsed BotFather token, e.g.
// 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11.
type Token struct {
	ID     int    // 123456
	Secret string // ABC-DEF1234ghIkl-zyx57W2v1u123ew11
}

// ParseToken parses a BotFather token.
func ParseToken(s string) (Token, error) {
	if s == "" {
		return Token{}, errors.New("blank token")
	}

	id, secret, ok := strings.Cut(s, ":")
	if !ok {
		return Token{}, errors.New("invalid token")
	}

	n, err := strconv.Atoi(id)
	if err != nil {
		return Token{}, errors.Wrap(err, "parse token id")
	}

	return Token{ID: n, Secret: secret}, nil
}

// String returns the token in its wire form.
func (t Token) String() string {
	return strconv.Itoa(t.ID) + ":" + t.Secret
}
