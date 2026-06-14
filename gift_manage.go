package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// ConvertGiftToStars converts a regular gift owned by the connected business
// account back to Telegram Stars. Requires the can_convert_gifts_to_stars
// business bot right.
func (b *Bot) ConvertGiftToStars(ctx context.Context, businessConnectionID, ownedGiftID string) error {
	gift, err := ownedGiftToTg(ownedGiftID)
	if err != nil {
		return err
	}

	var res tg.BoolBox

	if err := b.invokeBusiness(ctx, businessConnectionID, &tg.PaymentsConvertStarGiftRequest{
		Stargift: gift,
	}, &res); err != nil {
		return asAPIError(err)
	}

	return nil
}

// UpgradeGift upgrades a regular gift owned by the connected business account to
// a unique (collectible) gift. Requires the can_transfer_and_upgrade_gifts
// business bot right. keepOriginalDetails keeps the original text, sender and
// receiver in the upgraded gift.
//
// Only gifts whose upgrade is already paid for are supported; a paid upgrade
// (the optional star_count of the Bot API method) is not yet implemented.
func (b *Bot) UpgradeGift(ctx context.Context, businessConnectionID, ownedGiftID string, keepOriginalDetails bool) error {
	gift, err := ownedGiftToTg(ownedGiftID)
	if err != nil {
		return err
	}

	var res tg.UpdatesBox

	if err := b.invokeBusiness(ctx, businessConnectionID, &tg.PaymentsUpgradeStarGiftRequest{
		KeepOriginalDetails: keepOriginalDetails,
		Stargift:            gift,
	}, &res); err != nil {
		return asAPIError(err)
	}

	return nil
}

// TransferGift transfers a unique gift owned by the connected business account
// to another user or channel chat. Requires the can_transfer_and_upgrade_gifts
// business bot right.
//
// Only free transfers are supported; a paid transfer (the optional star_count
// of the Bot API method) is not yet implemented.
func (b *Bot) TransferGift(ctx context.Context, businessConnectionID, ownedGiftID string, newOwner ChatID) error {
	gift, err := ownedGiftToTg(ownedGiftID)
	if err != nil {
		return err
	}

	peer, err := b.resolveInputPeer(ctx, newOwner)
	if err != nil {
		return err
	}

	var res tg.UpdatesBox

	if err := b.invokeBusiness(ctx, businessConnectionID, &tg.PaymentsTransferStarGiftRequest{
		Stargift: gift,
		ToID:     peer,
	}, &res); err != nil {
		return asAPIError(err)
	}

	return nil
}
