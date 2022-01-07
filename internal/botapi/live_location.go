package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// EditMessageLiveLocation implements oas.Handler.
func (b *BotAPI) EditMessageLiveLocation(ctx context.Context, req oas.EditMessageLiveLocation) (oas.ResultMessageOrBoolean, error) {
	return oas.ResultMessageOrBoolean{}, &NotImplementedError{}
}

// StopMessageLiveLocation implements oas.Handler.
func (b *BotAPI) StopMessageLiveLocation(ctx context.Context, req oas.StopMessageLiveLocation) (oas.ResultMessageOrBoolean, error) {
	return oas.ResultMessageOrBoolean{}, &NotImplementedError{}
}
