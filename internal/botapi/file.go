package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// GetFile implements oas.Handler.
func (b *BotAPI) GetFile(ctx context.Context, req oas.GetFile) (oas.ResultFile, error) {
	return oas.ResultFile{}, &NotImplementedError{}
}

// UploadStickerFile implements oas.Handler.
func (b *BotAPI) UploadStickerFile(ctx context.Context, req oas.UploadStickerFile) (oas.ResultFile, error) {
	return oas.ResultFile{}, &NotImplementedError{}
}
