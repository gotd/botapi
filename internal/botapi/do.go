package botapi

import (
	"context"

	"github.com/gotd/td/telegram"
)

func (b *BotAPI) do(ctx context.Context, cb func(client *telegram.Client) error) error {
	return b.pool.Do(ctx, MustToken(ctx), cb)
}
