package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// DeleteWebhook implements oas.Handler.
func (b *BotAPI) DeleteWebhook(ctx context.Context, req oas.OptDeleteWebhook) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}

// GetWebhookInfo implements oas.Handler.
func (b *BotAPI) GetWebhookInfo(ctx context.Context) (*oas.ResultWebhookInfo, error) {
	return nil, &NotImplementedError{}
}

// SetWebhook implements oas.Handler.
func (b *BotAPI) SetWebhook(ctx context.Context, req *oas.SetWebhook) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}
