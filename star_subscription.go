package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// EditUserStarSubscription cancels or re-enables a user's Telegram Stars
// subscription to the bot's paid content. telegramPaymentChargeID identifies the
// subscription; isCanceled cancels it when true and re-enables it when false.
func (b *Bot) EditUserStarSubscription(ctx context.Context, userID int64, telegramPaymentChargeID string, isCanceled bool) error {
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	req := &tg.PaymentsChangeStarsSubscriptionRequest{
		Peer:           inputPeerFromUser(user),
		SubscriptionID: telegramPaymentChargeID,
	}
	req.SetCanceled(isCanceled)

	if _, err := b.raw.PaymentsChangeStarsSubscription(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}
