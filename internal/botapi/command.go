package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// GetMyCommands implements oas.Handler.
func (b *BotAPI) GetMyCommands(ctx context.Context, req oas.GetMyCommands) (oas.ResultArrayOfBotCommand, error) {
	return oas.ResultArrayOfBotCommand{}, &NotImplementedError{}
}

// SetMyCommands implements oas.Handler.
func (b *BotAPI) SetMyCommands(ctx context.Context, req oas.SetMyCommands) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// DeleteMyCommands implements oas.Handler.
func (b *BotAPI) DeleteMyCommands(ctx context.Context, req oas.DeleteMyCommands) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
