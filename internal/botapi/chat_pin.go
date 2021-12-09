package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// PinChatMessage implements oas.Handler.
func (b *BotAPI) PinChatMessage(ctx context.Context, req oas.PinChatMessage) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// UnpinAllChatMessages implements oas.Handler.
func (b *BotAPI) UnpinAllChatMessages(ctx context.Context, req oas.UnpinAllChatMessages) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// UnpinChatMessage implements oas.Handler.
func (b *BotAPI) UnpinChatMessage(ctx context.Context, req oas.UnpinChatMessage) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
