package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// AcceptedGiftTypes describes the types of gifts that can be gifted to a user or
// a chat.
type AcceptedGiftTypes struct {
	// UnlimitedGifts reports whether unlimited regular gifts are accepted.
	UnlimitedGifts bool `json:"unlimited_gifts"`
	// LimitedGifts reports whether limited regular gifts are accepted.
	LimitedGifts bool `json:"limited_gifts"`
	// UniqueGifts reports whether unique (collectible) gifts are accepted.
	UniqueGifts bool `json:"unique_gifts"`
	// PremiumSubscription reports whether a Telegram Premium subscription is
	// accepted as a gift.
	PremiumSubscription bool `json:"premium_subscription"`
	// GiftsFromChannels reports whether gifts sent on behalf of channels are
	// accepted.
	GiftsFromChannels bool `json:"gifts_from_channels"`
}

// disallowed maps the accepted gift types onto the MTProto disallowed-gifts
// settings, which is expressed as the inverse (a set of disallow flags).
func (a AcceptedGiftTypes) disallowed() tg.DisallowedGiftsSettings {
	return tg.DisallowedGiftsSettings{
		DisallowUnlimitedStargifts:    !a.UnlimitedGifts,
		DisallowLimitedStargifts:      !a.LimitedGifts,
		DisallowUniqueStargifts:       !a.UniqueGifts,
		DisallowPremiumGifts:          !a.PremiumSubscription,
		DisallowStargiftsFromChannels: !a.GiftsFromChannels,
	}
}

// SetBusinessAccountGiftSettings changes the privacy settings for incoming gifts
// of a managed business account. Requires the can_change_gift_settings business
// bot right.
//
// The account's other global privacy settings are read first and preserved, so
// only the gift button and accepted gift types change.
func (b *Bot) SetBusinessAccountGiftSettings(
	ctx context.Context, businessConnectionID string, showGiftButton bool, acceptedGiftTypes AcceptedGiftTypes,
) error {
	var current tg.GlobalPrivacySettings

	err := b.invokeBusiness(ctx, businessConnectionID, &tg.AccountGetGlobalPrivacySettingsRequest{}, &current)
	if err != nil {
		return asAPIError(err)
	}

	current.SetDisplayGiftsButton(showGiftButton)
	current.SetDisallowedGifts(acceptedGiftTypes.disallowed())

	var updated tg.GlobalPrivacySettings

	err = b.invokeBusiness(ctx, businessConnectionID, &tg.AccountSetGlobalPrivacySettingsRequest{
		Settings: current,
	}, &updated)
	if err != nil {
		return asAPIError(err)
	}

	return nil
}
