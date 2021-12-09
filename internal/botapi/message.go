package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// CopyMessage implements oas.Handler.
func (b *BotAPI) CopyMessage(ctx context.Context, req oas.CopyMessage) (oas.ResultMessageId, error) {
	return oas.ResultMessageId{}, &NotImplementedError{}
}

// DeleteMessage implements oas.Handler.
func (b *BotAPI) DeleteMessage(ctx context.Context, req oas.DeleteMessage) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// EditMessageCaption implements oas.Handler.
func (b *BotAPI) EditMessageCaption(ctx context.Context, req oas.EditMessageCaption) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// EditMessageLiveLocation implements oas.Handler.
func (b *BotAPI) EditMessageLiveLocation(ctx context.Context, req oas.EditMessageLiveLocation) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// EditMessageMedia implements oas.Handler.
func (b *BotAPI) EditMessageMedia(ctx context.Context, req oas.EditMessageMedia) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// EditMessageReplyMarkup implements oas.Handler.
func (b *BotAPI) EditMessageReplyMarkup(ctx context.Context, req oas.EditMessageReplyMarkup) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// EditMessageText implements oas.Handler.
func (b *BotAPI) EditMessageText(ctx context.Context, req oas.EditMessageText) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// ForwardMessage implements oas.Handler.
func (b *BotAPI) ForwardMessage(ctx context.Context, req oas.ForwardMessage) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// StopMessageLiveLocation implements oas.Handler.
func (b *BotAPI) StopMessageLiveLocation(ctx context.Context, req oas.StopMessageLiveLocation) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// StopPoll implements oas.Handler.
func (b *BotAPI) StopPoll(ctx context.Context, req oas.StopPoll) (oas.ResultPoll, error) {
	return oas.ResultPoll{}, &NotImplementedError{}
}
