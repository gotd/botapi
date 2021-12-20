package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// GetUserProfilePhotos implements oas.Handler.
func (b *BotAPI) GetUserProfilePhotos(ctx context.Context, req oas.GetUserProfilePhotos) (oas.ResultUserProfilePhotos, error) {
	return oas.ResultUserProfilePhotos{}, &NotImplementedError{}
}
