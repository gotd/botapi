package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// EmojiStatusOption configures a SetUserEmojiStatus call.
type EmojiStatusOption func(*emojiStatusConfig)

type emojiStatusConfig struct {
	until int
}

// WithEmojiStatusExpiration sets the Unix time when the emoji status will
// expire. By default the status does not expire.
func WithEmojiStatusExpiration(unixTime int) EmojiStatusOption {
	return func(c *emojiStatusConfig) { c.until = unixTime }
}

// SetUserEmojiStatus sets the emoji status of a user that previously allowed the
// bot to manage it. An empty customEmojiID removes the current emoji status.
func (b *Bot) SetUserEmojiStatus(ctx context.Context, userID int64, customEmojiID string, opts ...EmojiStatusOption) error {
	var cfg emojiStatusConfig

	for _, o := range opts {
		o(&cfg)
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	var status tg.EmojiStatusClass = &tg.EmojiStatusEmpty{}

	if customEmojiID != "" {
		id, err := strconv.ParseInt(customEmojiID, 10, 64)
		if err != nil {
			return errInvalidCustomEmojiID()
		}

		s := &tg.EmojiStatus{DocumentID: id}
		if cfg.until != 0 {
			s.SetUntil(cfg.until)
		}

		status = s
	}

	if _, err := b.raw.BotsUpdateUserEmojiStatus(ctx, &tg.BotsUpdateUserEmojiStatusRequest{
		UserID:      user,
		EmojiStatus: status,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
