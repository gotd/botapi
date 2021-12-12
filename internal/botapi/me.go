package botapi

import (
	"context"

	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func convertToUser(user *tg.User) oas.User {
	return oas.User{
		ID:                      user.ID,
		IsBot:                   user.Bot,
		FirstName:               user.FirstName,
		LastName:                optString(user.GetLastName),
		Username:                optString(user.GetUsername),
		LanguageCode:            optString(user.GetLangCode),
		CanJoinGroups:           oas.NewOptBool(user.BotNochats),
		CanReadAllGroupMessages: oas.NewOptBool(user.BotChatHistory),
		SupportsInlineQueries:   oas.NewOptBool(user.BotInlinePlaceholder == ""),
	}
}

// GetMe implements oas.Handler.
func (b *BotAPI) GetMe(ctx context.Context) (oas.ResultUser, error) {
	me, err := b.client.Self(ctx)
	if err != nil {
		return oas.ResultUser{}, err
	}
	b.updateSelf(me)

	return oas.ResultUser{
		Result: oas.NewOptUser(convertToUser(me)),
		Ok:     true,
	}, nil
}

// Close implements oas.Handler.
func (b *BotAPI) Close(ctx context.Context) (oas.Result, error) {
	// FIXME(tdakkota): kill BotAPI.
	return resultOK(true), nil
}

// LogOut implements oas.Handler.
func (b *BotAPI) LogOut(ctx context.Context) (oas.Result, error) {
	r, err := b.raw.AuthLogOut(ctx)
	if err != nil {
		return oas.Result{}, err
	}

	return resultOK(r), nil
}
