package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// AnswerCallbackQueryOption configures an AnswerCallbackQuery call.
type AnswerCallbackQueryOption func(*answerCallbackConfig)

type answerCallbackConfig struct {
	text      string
	showAlert bool
	url       string
	cacheTime int
}

// WithCallbackText sets the notification text shown to the user (0-200
// characters). By default nothing is shown.
func WithCallbackText(text string) AnswerCallbackQueryOption {
	return func(c *answerCallbackConfig) { c.text = text }
}

// WithCallbackAlert shows the notification as an alert dialog instead of a
// top-of-screen toast.
func WithCallbackAlert() AnswerCallbackQueryOption {
	return func(c *answerCallbackConfig) { c.showAlert = true }
}

// WithCallbackURL sets a URL to open. Telegram restricts which URLs are
// accepted (game URLs, or t.me/your_bot?start= links).
func WithCallbackURL(url string) AnswerCallbackQueryOption {
	return func(c *answerCallbackConfig) { c.url = url }
}

// WithCallbackCacheTime sets, in seconds, how long the result may be cached
// client-side.
func WithCallbackCacheTime(seconds int) AnswerCallbackQueryOption {
	return func(c *answerCallbackConfig) { c.cacheTime = seconds }
}

// AnswerCallbackQuery responds to a callback query sent from an inline keyboard.
// The callbackQueryID is the CallbackQuery.ID from the update.
func (b *Bot) AnswerCallbackQuery(ctx context.Context, callbackQueryID string, opts ...AnswerCallbackQueryOption) error {
	var cfg answerCallbackConfig
	for _, o := range opts {
		o(&cfg)
	}

	queryID, err := strconv.ParseInt(callbackQueryID, 10, 64)
	if err != nil {
		return &Error{Code: 400, Description: "Bad Request: invalid callback query id"}
	}

	req := &tg.MessagesSetBotCallbackAnswerRequest{
		Alert:     cfg.showAlert,
		QueryID:   queryID,
		CacheTime: cfg.cacheTime,
	}
	if cfg.text != "" {
		req.SetMessage(cfg.text)
	}
	if cfg.url != "" {
		req.SetURL(cfg.url)
	}

	if _, err := b.raw.MessagesSetBotCallbackAnswer(ctx, req); err != nil {
		return asAPIError(err)
	}
	return nil
}
