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

// AnswerInlineQueryOption configures an AnswerInlineQuery call.
type AnswerInlineQueryOption func(*answerInlineConfig)

type answerInlineConfig struct {
	cacheTime     int
	isPersonal    bool
	nextOffset    string
	switchPMText  string
	switchPMParam string
}

// WithInlineCacheTime sets the maximum time, in seconds, the result may be
// cached on the server.
func WithInlineCacheTime(seconds int) AnswerInlineQueryOption {
	return func(c *answerInlineConfig) { c.cacheTime = seconds }
}

// WithInlinePersonal marks the results as personal so they are not cached for
// other users querying the same thing.
func WithInlinePersonal() AnswerInlineQueryOption {
	return func(c *answerInlineConfig) { c.isPersonal = true }
}

// WithInlineNextOffset sets the offset a client returns to request the next page
// of results. An empty offset (the default) means no more results.
func WithInlineNextOffset(offset string) AnswerInlineQueryOption {
	return func(c *answerInlineConfig) { c.nextOffset = offset }
}

// WithInlineSwitchPM shows a button above the results that switches the user to
// a private chat with the bot, passing startParam as the /start parameter.
func WithInlineSwitchPM(text, startParam string) AnswerInlineQueryOption {
	return func(c *answerInlineConfig) {
		c.switchPMText = text
		c.switchPMParam = startParam
	}
}

// AnswerInlineQuery sends the results of an inline query. inlineQueryID is the
// InlineQuery.ID from the update.
func (b *Bot) AnswerInlineQuery(
	ctx context.Context, inlineQueryID string, results []InlineQueryResult, opts ...AnswerInlineQueryOption,
) error {
	var cfg answerInlineConfig
	for _, o := range opts {
		o(&cfg)
	}

	queryID, err := strconv.ParseInt(inlineQueryID, 10, 64)
	if err != nil {
		return &Error{Code: 400, Description: "Bad Request: invalid inline query id"}
	}

	tgResults := make([]tg.InputBotInlineResultClass, 0, len(results))
	for _, r := range results {
		if r == nil {
			return &Error{Code: 400, Description: "Bad Request: inline query result is nil"}
		}
		converted, err := r.toTg(ctx, b)
		if err != nil {
			return err
		}
		tgResults = append(tgResults, converted)
	}

	req := &tg.MessagesSetInlineBotResultsRequest{
		QueryID:    queryID,
		Results:    tgResults,
		CacheTime:  cfg.cacheTime,
		Private:    cfg.isPersonal,
		NextOffset: cfg.nextOffset,
	}
	if cfg.switchPMText != "" {
		req.SetSwitchPm(tg.InlineBotSwitchPM{
			Text:       cfg.switchPMText,
			StartParam: cfg.switchPMParam,
		})
	}

	if _, err := b.raw.MessagesSetInlineBotResults(ctx, req); err != nil {
		return asAPIError(err)
	}
	return nil
}
