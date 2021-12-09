package pool

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/go-faster/errors"
)

// Token represents bot token, like 123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11
type Token struct {
	ID     int    // 123456
	Secret string // ABC-DEF1234ghIkl-zyx57W2v1u123ew11
}

func ParseToken(s string) (Token, error) {
	if s == "" {
		return Token{}, errors.New("blank")
	}
	parts := strings.Split(s, ":")
	if len(parts) != 2 {
		return Token{}, errors.New("invalid token")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return Token{}, err
	}
	return Token{
		ID:     id,
		Secret: parts[1],
	}, err
}

func (t Token) String() string {
	return fmt.Sprintf("%d:%s", t.ID, t.Secret)
}
