package botapi

import (
	"context"

	"github.com/go-faster/errors"

	"github.com/gotd/td/tg"
)

func (b *BotAPI) getStickerSet(ctx context.Context, input tg.InputStickerSetClass) (*tg.MessagesStickerSet, error) {
	// TODO(tdakkota): investigate GreatMinds hack
	//  See https://github.com/tdlib/telegram-bot-api/blob/6abdb73512110c2adfaa7145eb01e102e75b89f6/telegram-bot-api/Client.h#L69-L70
	result, err := b.raw.MessagesGetStickerSet(ctx, &tg.MessagesGetStickerSetRequest{
		Stickerset: input,
		Hash:       0,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get sticker_set")
	}
	// TODO(tdakkota): make cache
	switch result := result.(type) {
	case *tg.MessagesStickerSet:
		return result, nil
	default:
		return nil, errors.Errorf("unexpected type %T", result)
	}
}
