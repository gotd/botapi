package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// giftConfig holds the optional parameters of SendGift.
type giftConfig struct {
	text          string
	parseMode     ParseMode
	entities      []MessageEntity
	payForUpgrade bool
}

// GiftOption customizes SendGift.
type GiftOption func(*giftConfig)

// WithGiftText attaches a text message to the gift.
func WithGiftText(text string) GiftOption {
	return func(c *giftConfig) { c.text = text }
}

// WithGiftParseMode sets the parse mode used for the gift text.
func WithGiftParseMode(mode ParseMode) GiftOption {
	return func(c *giftConfig) { c.parseMode = mode }
}

// WithGiftEntities sets explicit entities for the gift text, overriding the
// parse mode.
func WithGiftEntities(entities []MessageEntity) GiftOption {
	return func(c *giftConfig) { c.entities = entities }
}

// WithGiftPayForUpgrade pays for the gift to be eligible for an upgrade to a
// unique (collectible) gift by the receiver.
func WithGiftPayForUpgrade() GiftOption {
	return func(c *giftConfig) { c.payForUpgrade = true }
}

// SendGift sends a gift, identified by an id from GetAvailableGifts, to a user
// or channel chat. The gift can't be converted to Telegram Stars by the
// receiver.
func (b *Bot) SendGift(ctx context.Context, target ChatID, giftID string, opts ...GiftOption) error {
	var cfg giftConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	id, err := strconv.ParseInt(giftID, 10, 64)
	if err != nil {
		return &Error{Code: 400, Description: "Bad Request: invalid gift_id"}
	}

	peer, err := b.resolveInputPeer(ctx, target)
	if err != nil {
		return err
	}

	invoice := &tg.InputInvoiceStarGift{
		Peer:           peer,
		GiftID:         id,
		IncludeUpgrade: cfg.payForUpgrade,
	}

	if cfg.text != "" {
		text, entities, err := b.giftMessage(ctx, cfg)
		if err != nil {
			return err
		}

		invoice.SetMessage(tg.TextWithEntities{Text: text, Entities: entities})
	}

	return b.payStarsForm(ctx, invoice)
}

// GiftPremiumSubscription gifts a Telegram Premium subscription to a user for
// the given number of months. The Telegram Stars cost is determined by the
// payment form for the chosen duration. WithGiftPayForUpgrade has no effect
// here.
func (b *Bot) GiftPremiumSubscription(ctx context.Context, userID int64, months int, opts ...GiftOption) error {
	var cfg giftConfig

	for _, opt := range opts {
		opt(&cfg)
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	invoice := &tg.InputInvoicePremiumGiftStars{
		UserID: user,
		Months: months,
	}

	if cfg.text != "" {
		text, entities, err := b.giftMessage(ctx, cfg)
		if err != nil {
			return err
		}

		invoice.SetMessage(tg.TextWithEntities{Text: text, Entities: entities})
	}

	return b.payStarsForm(ctx, invoice)
}

// giftMessage resolves the gift text into a (text, entities) pair: explicit
// entities take precedence over the parse mode.
func (b *Bot) giftMessage(ctx context.Context, cfg giftConfig) (string, []tg.MessageEntityClass, error) {
	if len(cfg.entities) > 0 {
		return cfg.text, entitiesToTg(cfg.entities), nil
	}

	return b.styledMessage(ctx, cfg.text, cfg.parseMode)
}

// payStarsForm runs the two-step Telegram Stars payment flow for an invoice:
// fetch the payment form, then submit it.
func (b *Bot) payStarsForm(ctx context.Context, invoice tg.InputInvoiceClass) error {
	form, err := b.raw.PaymentsGetPaymentForm(ctx, &tg.PaymentsGetPaymentFormRequest{
		Invoice: invoice,
	})
	if err != nil {
		return asAPIError(err)
	}

	formID, err := starsFormID(form)
	if err != nil {
		return err
	}

	if _, err := b.raw.PaymentsSendStarsForm(ctx, &tg.PaymentsSendStarsFormRequest{
		FormID:  formID,
		Invoice: invoice,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// starsFormID extracts the form id from a Telegram Stars payment form. Star
// gifts return a starGift form; subscriptions and other star purchases return a
// plain stars form.
func starsFormID(form tg.PaymentsPaymentFormClass) (int64, error) {
	switch f := form.(type) {
	case *tg.PaymentsPaymentFormStarGift:
		return f.FormID, nil
	case *tg.PaymentsPaymentFormStars:
		return f.FormID, nil
	default:
		return 0, &Error{Code: 500, Description: "Internal Server Error: unexpected payment form"}
	}
}
