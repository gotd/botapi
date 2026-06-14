package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// CreateInvoiceLink creates a link for an invoice and returns it. The created
// invoice can be paid by any user who follows the link.
func (b *Bot) CreateInvoiceLink(ctx context.Context, params InvoiceParams) (string, error) {
	res, err := b.raw.PaymentsExportInvoice(ctx, invoiceMedia(params))
	if err != nil {
		return "", asAPIError(err)
	}

	return res.URL, nil
}

// RefundStarPayment refunds a successful payment in Telegram Stars to the given
// user. chargeID is the telegram_payment_charge_id from the successful payment.
func (b *Bot) RefundStarPayment(ctx context.Context, userID int64, chargeID string) error {
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.PaymentsRefundStarsCharge(ctx, &tg.PaymentsRefundStarsChargeRequest{
		UserID:   user,
		ChargeID: chargeID,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
