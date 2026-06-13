package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// SendAnimation implements oas.Handler.
func (b *BotAPI) SendAnimation(ctx context.Context, req *oas.SendAnimation) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendAudio implements oas.Handler.
func (b *BotAPI) SendAudio(ctx context.Context, req *oas.SendAudio) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendDocument implements oas.Handler.
func (b *BotAPI) SendDocument(ctx context.Context, req *oas.SendDocument) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendMediaGroup implements oas.Handler.
func (b *BotAPI) SendMediaGroup(ctx context.Context, req *oas.SendMediaGroup) (*oas.ResultArrayOfMessage, error) {
	return nil, &NotImplementedError{}
}

// SendPhoto implements oas.Handler.
func (b *BotAPI) SendPhoto(ctx context.Context, req *oas.SendPhoto) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendSticker implements oas.Handler.
func (b *BotAPI) SendSticker(ctx context.Context, req *oas.SendSticker) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendVideo implements oas.Handler.
func (b *BotAPI) SendVideo(ctx context.Context, req *oas.SendVideo) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendVideoNote implements oas.Handler.
func (b *BotAPI) SendVideoNote(ctx context.Context, req *oas.SendVideoNote) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}

// SendVoice implements oas.Handler.
func (b *BotAPI) SendVoice(ctx context.Context, req *oas.SendVoice) (*oas.ResultMessage, error) {
	return nil, &NotImplementedError{}
}
