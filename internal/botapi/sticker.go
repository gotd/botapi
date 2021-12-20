package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// AddStickerToSet implements oas.Handler.
func (b *BotAPI) AddStickerToSet(ctx context.Context, req oas.AddStickerToSet) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// CreateNewStickerSet implements oas.Handler.
func (b *BotAPI) CreateNewStickerSet(ctx context.Context, req oas.CreateNewStickerSet) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// DeleteStickerFromSet implements oas.Handler.
func (b *BotAPI) DeleteStickerFromSet(ctx context.Context, req oas.DeleteStickerFromSet) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// GetStickerSet implements oas.Handler.
func (b *BotAPI) GetStickerSet(ctx context.Context, req oas.GetStickerSet) (oas.ResultStickerSet, error) {
	return oas.ResultStickerSet{}, &NotImplementedError{}
}

// SetStickerPositionInSet implements oas.Handler.
func (b *BotAPI) SetStickerPositionInSet(ctx context.Context, req oas.SetStickerPositionInSet) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// SetStickerSetThumb implements oas.Handler.
func (b *BotAPI) SetStickerSetThumb(ctx context.Context, req oas.SetStickerSetThumb) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
