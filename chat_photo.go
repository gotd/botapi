package botapi

import (
	"context"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"
)

// SetChatPhoto sets a new profile photo for a chat. The photo must be a local
// upload (file_id and URL are not accepted by Telegram for this method).
func (b *Bot) SetChatPhoto(ctx context.Context, chat ChatID, photo InputFile) error {
	up, ok := photo.(*InputFileUpload)
	if !ok {
		return &Error{Code: 400, Description: "Bad Request: chat photo must be an uploaded file"}
	}

	uploaded, err := b.uploadInputFile(ctx, up)
	if err != nil {
		return err
	}

	input := &tg.InputChatUploadedPhoto{}
	input.SetFile(uploaded)

	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return err
	}

	switch c := p.(type) {
	case peers.Channel:
		if err := c.SetPhoto(ctx, input); err != nil {
			return asAPIError(err)
		}

		return nil
	case peers.Chat:
		if _, err := b.raw.MessagesEditChatPhoto(ctx, &tg.MessagesEditChatPhotoRequest{
			ChatID: c.ID(),
			Photo:  input,
		}); err != nil {
			return asAPIError(err)
		}

		return nil
	default:
		return errNotInPrivateChat()
	}
}

// DeleteChatPhoto removes the profile photo of a chat.
func (b *Bot) DeleteChatPhoto(ctx context.Context, chat ChatID) error {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return err
	}

	switch c := p.(type) {
	case peers.Channel:
		if err := c.DeletePhoto(ctx); err != nil {
			return asAPIError(err)
		}

		return nil
	case peers.Chat:
		if _, err := b.raw.MessagesEditChatPhoto(ctx, &tg.MessagesEditChatPhotoRequest{
			ChatID: c.ID(),
			Photo:  &tg.InputChatPhotoEmpty{},
		}); err != nil {
			return asAPIError(err)
		}

		return nil
	default:
		return errNotInPrivateChat()
	}
}
