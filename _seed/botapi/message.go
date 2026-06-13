package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// CopyMessage implements oas.Handler.
func (b *BotAPI) CopyMessage(ctx context.Context, req *oas.CopyMessage) (*oas.ResultMessageId, error) {
	return nil, &NotImplementedError{}
}

// DeleteMessage implements oas.Handler.
func (b *BotAPI) DeleteMessage(ctx context.Context, req *oas.DeleteMessage) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}

// EditMessageCaption implements oas.Handler.
func (b *BotAPI) EditMessageCaption(ctx context.Context, req *oas.EditMessageCaption) (*oas.ResultMessageOrBoolean, error) {
	return nil, &NotImplementedError{}
}

// EditMessageMedia implements oas.Handler.
func (b *BotAPI) EditMessageMedia(ctx context.Context, req *oas.EditMessageMedia) (*oas.ResultMessageOrBoolean, error) {
	return nil, &NotImplementedError{}
}

// EditMessageReplyMarkup implements oas.Handler.
func (b *BotAPI) EditMessageReplyMarkup(ctx context.Context, req *oas.EditMessageReplyMarkup) (*oas.ResultMessageOrBoolean, error) {
	return nil, &NotImplementedError{}
}

// EditMessageText implements oas.Handler.
func (b *BotAPI) EditMessageText(ctx context.Context, req *oas.EditMessageText) (*oas.ResultMessageOrBoolean, error) {
	return nil, &NotImplementedError{}
}

// ForwardMessage implements oas.Handler.
func (b *BotAPI) ForwardMessage(ctx context.Context, req *oas.ForwardMessage) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}
