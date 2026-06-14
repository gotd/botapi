package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// AnswerShippingQuery replies to a shipping query. When ok is true, provide the
// available shipping options; otherwise pass an error message via
// WithShippingError describing why the order can't be shipped.
func (b *Bot) AnswerShippingQuery(ctx context.Context, shippingQueryID string, ok bool, opts ...ShippingAnswerOption) error {
	var cfg shippingAnswerConfig

	for _, o := range opts {
		o(&cfg)
	}

	queryID, err := strconv.ParseInt(shippingQueryID, 10, 64)
	if err != nil {
		return &Error{Code: 400, Description: "Bad Request: invalid shipping query id"}
	}

	req := &tg.MessagesSetBotShippingResultsRequest{QueryID: queryID}
	if ok {
		req.ShippingOptions = shippingOptionsToTg(cfg.options)
	} else {
		req.Error = cfg.errorMessage
	}

	if _, err := b.raw.MessagesSetBotShippingResults(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// ShippingAnswerOption configures an AnswerShippingQuery call.
type ShippingAnswerOption func(*shippingAnswerConfig)

type shippingAnswerConfig struct {
	options      []ShippingOption
	errorMessage string
}

// WithShippingOptions sets the available shipping options (required when ok).
func WithShippingOptions(options ...ShippingOption) ShippingAnswerOption {
	return func(c *shippingAnswerConfig) { c.options = options }
}

// WithShippingError sets the human-readable reason shipping is impossible (used
// when ok is false).
func WithShippingError(message string) ShippingAnswerOption {
	return func(c *shippingAnswerConfig) { c.errorMessage = message }
}

// AnswerPreCheckoutQuery replies to a pre-checkout query. When ok is false,
// provide a human-readable reason via WithPreCheckoutError.
func (b *Bot) AnswerPreCheckoutQuery(ctx context.Context, preCheckoutQueryID string, ok bool, opts ...PreCheckoutAnswerOption) error {
	var cfg preCheckoutAnswerConfig

	for _, o := range opts {
		o(&cfg)
	}

	queryID, err := strconv.ParseInt(preCheckoutQueryID, 10, 64)
	if err != nil {
		return &Error{Code: 400, Description: "Bad Request: invalid pre-checkout query id"}
	}

	req := &tg.MessagesSetBotPrecheckoutResultsRequest{
		Success: ok,
		QueryID: queryID,
	}
	if !ok {
		req.Error = cfg.errorMessage
	}

	if _, err := b.raw.MessagesSetBotPrecheckoutResults(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// PreCheckoutAnswerOption configures an AnswerPreCheckoutQuery call.
type PreCheckoutAnswerOption func(*preCheckoutAnswerConfig)

type preCheckoutAnswerConfig struct {
	errorMessage string
}

// WithPreCheckoutError sets the reason the checkout can't proceed (used when ok
// is false).
func WithPreCheckoutError(message string) PreCheckoutAnswerOption {
	return func(c *preCheckoutAnswerConfig) { c.errorMessage = message }
}

// shippingOptionsToTg converts Bot API shipping options to MTProto.
func shippingOptionsToTg(options []ShippingOption) []tg.ShippingOption {
	out := make([]tg.ShippingOption, 0, len(options))
	for _, o := range options {
		out = append(out, tg.ShippingOption{ID: o.ID, Title: o.Title, Prices: pricesToTg(o.Prices)})
	}

	return out
}
