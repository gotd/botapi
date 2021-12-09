package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// SendAnimation implements oas.Handler.
func (b *BotAPI) SendAnimation(ctx context.Context, req oas.SendAnimation) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendAudio implements oas.Handler.
func (b *BotAPI) SendAudio(ctx context.Context, req oas.SendAudio) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendChatAction implements oas.Handler.
func (b *BotAPI) SendChatAction(ctx context.Context, req oas.SendChatAction) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// SendContact implements oas.Handler.
func (b *BotAPI) SendContact(ctx context.Context, req oas.SendContact) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendDice implements oas.Handler.
func (b *BotAPI) SendDice(ctx context.Context, req oas.SendDice) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendDocument implements oas.Handler.
func (b *BotAPI) SendDocument(ctx context.Context, req oas.SendDocument) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendGame implements oas.Handler.
func (b *BotAPI) SendGame(ctx context.Context, req oas.SendGame) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendInvoice implements oas.Handler.
func (b *BotAPI) SendInvoice(ctx context.Context, req oas.SendInvoice) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendLocation implements oas.Handler.
func (b *BotAPI) SendLocation(ctx context.Context, req oas.SendLocation) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendMediaGroup implements oas.Handler.
func (b *BotAPI) SendMediaGroup(ctx context.Context, req oas.SendMediaGroup) (oas.ResultArrayOfMessage, error) {
	return oas.ResultArrayOfMessage{}, &NotImplementedError{}
}

// SendMessage implements oas.Handler.
func (b *BotAPI) SendMessage(ctx context.Context, req oas.SendMessage) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendPhoto implements oas.Handler.
func (b *BotAPI) SendPhoto(ctx context.Context, req oas.SendPhoto) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendPoll implements oas.Handler.
func (b *BotAPI) SendPoll(ctx context.Context, req oas.SendPoll) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendSticker implements oas.Handler.
func (b *BotAPI) SendSticker(ctx context.Context, req oas.SendSticker) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVenue implements oas.Handler.
func (b *BotAPI) SendVenue(ctx context.Context, req oas.SendVenue) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVideo implements oas.Handler.
func (b *BotAPI) SendVideo(ctx context.Context, req oas.SendVideo) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVideoNote implements oas.Handler.
func (b *BotAPI) SendVideoNote(ctx context.Context, req oas.SendVideoNote) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}

// SendVoice implements oas.Handler.
func (b *BotAPI) SendVoice(ctx context.Context, req oas.SendVoice) (oas.ResultMessage, error) {
	return oas.ResultMessage{}, &NotImplementedError{}
}
