package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// StarAmount describes an amount of Telegram Stars.
type StarAmount struct {
	// Amount is the integer number of Telegram Stars, rounded to 0; can be
	// negative.
	Amount int `json:"amount"`
	// NanostarAmount is the number of 1/1000000000 shares of Telegram Stars,
	// rounded to 0; can be negative.
	NanostarAmount int `json:"nanostar_amount,omitempty"`
}

// GetMyStarBalance returns the current Telegram Stars balance of the bot.
func (b *Bot) GetMyStarBalance(ctx context.Context) (StarAmount, error) {
	res, err := b.raw.PaymentsGetStarsStatus(ctx, &tg.PaymentsGetStarsStatusRequest{
		Peer: &tg.InputPeerSelf{},
	})
	if err != nil {
		return StarAmount{}, asAPIError(err)
	}

	amount, ok := res.Balance.(*tg.StarsAmount)
	if !ok {
		return StarAmount{}, &Error{Code: 500, Description: "Internal Server Error: unexpected stars balance"}
	}

	return StarAmount{Amount: int(amount.Amount), NanostarAmount: amount.Nanos}, nil
}
