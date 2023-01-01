package botapi

import (
	"context"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func convertRawToBotAPIUser(user *tg.User) oas.User {
	return oas.User{
		ID:                      user.ID,
		IsBot:                   user.Bot,
		FirstName:               user.FirstName,
		LastName:                optString(user.GetLastName),
		Username:                optString(user.GetUsername),
		LanguageCode:            optString(user.GetLangCode),
		CanJoinGroups:           oas.NewOptBool(!user.BotNochats),
		CanReadAllGroupMessages: oas.NewOptBool(user.BotChatHistory),
		SupportsInlineQueries:   oas.NewOptBool(user.BotInlinePlaceholder != ""),
	}
}

func convertToBotAPIUser(u peers.User) oas.User {
	return convertRawToBotAPIUser(u.Raw())
}

// GetMe implements oas.Handler.
func (b *BotAPI) GetMe(ctx context.Context) (*oas.ResultUser, error) {
	me, err := b.peers.Self(ctx)
	if err != nil {
		return nil, err
	}

	return &oas.ResultUser{
		Result: oas.NewOptUser(convertToBotAPIUser(me)),
		Ok:     true,
	}, nil
}

// Close implements oas.Handler.
func (b *BotAPI) Close(ctx context.Context) (*oas.Result, error) {
	// FIXME(tdakkota): kill BotAPI.
	return resultOK(true), nil
}

// LogOut implements oas.Handler.
func (b *BotAPI) LogOut(ctx context.Context) (*oas.Result, error) {
	if _, err := b.raw.AuthLogOut(ctx); err != nil {
		return nil, err
	}

	return resultOK(true), nil
}
