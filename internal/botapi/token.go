package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/pool"
)

type tokenKey struct{}

// PropagateToken adds given token to context.
func PropagateToken(ctx context.Context, token pool.Token) context.Context {
	return context.WithValue(ctx, tokenKey{}, token)
}

// MustToken gets pool.Token from context if any.
//
// Panics otherwise.
func MustToken(ctx context.Context) pool.Token {
	return ctx.Value(tokenKey{}).(pool.Token)
}
