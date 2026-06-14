package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// ManagedBotAccessSettings describes the access settings of a managed bot.
type ManagedBotAccessSettings struct {
	// IsAccessRestricted reports whether access to the managed bot is restricted
	// to an allow-list of users (AddedUserIDs).
	IsAccessRestricted bool `json:"is_access_restricted,omitempty"`
	// AddedUserIDs are the users explicitly granted access to the managed bot.
	AddedUserIDs []int64 `json:"added_user_ids,omitempty"`
}

// GetManagedBotAccessSettings returns the access settings of a managed bot.
func (b *Bot) GetManagedBotAccessSettings(ctx context.Context, userID int64) (*ManagedBotAccessSettings, error) {
	bot, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	res, err := b.raw.BotsGetAccessSettings(ctx, bot)
	if err != nil {
		return nil, asAPIError(err)
	}

	out := &ManagedBotAccessSettings{IsAccessRestricted: res.Restricted}

	for _, u := range res.AddUsers {
		if user, ok := u.(*tg.User); ok {
			out.AddedUserIDs = append(out.AddedUserIDs, user.ID)
		}
	}

	return out, nil
}

// SetManagedBotAccessSettings changes the access settings of a managed bot. When
// restricted is true, only the users in addUserIDs may access the bot.
func (b *Bot) SetManagedBotAccessSettings(ctx context.Context, userID int64, restricted bool, addUserIDs []int64) error {
	bot, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	addUsers := make([]tg.InputUserClass, 0, len(addUserIDs))

	for _, id := range addUserIDs {
		u, err := b.resolveInputUser(ctx, id)
		if err != nil {
			return err
		}

		addUsers = append(addUsers, u)
	}

	req := &tg.BotsEditAccessSettingsRequest{
		Restricted: restricted,
		Bot:        bot,
		AddUsers:   addUsers,
	}

	if _, err := b.raw.BotsEditAccessSettings(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// GetManagedBotToken returns the bot token of a managed bot.
func (b *Bot) GetManagedBotToken(ctx context.Context, userID int64) (string, error) {
	return b.exportManagedBotToken(ctx, userID, false)
}

// ReplaceManagedBotToken revokes the current token of a managed bot and returns
// a freshly generated one. The previous token stops working.
func (b *Bot) ReplaceManagedBotToken(ctx context.Context, userID int64) (string, error) {
	return b.exportManagedBotToken(ctx, userID, true)
}

// exportManagedBotToken exports (revoke=false) or regenerates (revoke=true) the
// token of a managed bot.
func (b *Bot) exportManagedBotToken(ctx context.Context, userID int64, revoke bool) (string, error) {
	bot, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return "", err
	}

	res, err := b.raw.BotsExportBotToken(ctx, &tg.BotsExportBotTokenRequest{
		Bot:    bot,
		Revoke: revoke,
	})
	if err != nil {
		return "", asAPIError(err)
	}

	return res.Token, nil
}
