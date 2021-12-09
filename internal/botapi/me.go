package botapi

import (
	"context"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi/internal/oas"
)

func convertUser(user *tg.User) oas.User {
	return oas.User{
		ID:                      int(user.ID),
		IsBot:                   user.Bot,
		FirstName:               user.FirstName,
		LastName:                optString(user.GetLastName),
		Username:                optString(user.GetUsername),
		LanguageCode:            optString(user.GetLangCode),
		CanJoinGroups:           optBool(user.BotNochats),
		CanReadAllGroupMessages: optBool(user.BotChatHistory),
		SupportsInlineQueries:   optBool(user.BotInlinePlaceholder == ""),
	}
}

// GetMe implements oas.Handler.
func (b *BotAPI) GetMe(ctx context.Context) (oas.ResultUser, error) {
	var self *tg.User
	if err := b.do(ctx, func(client *telegram.Client) (err error) {
		self, err = client.Self(ctx)
		return err
	}); err != nil {
		return oas.ResultUser{}, err
	}

	return oas.ResultUser{
		Result: oas.NewOptUser(convertUser(self)),
		Ok:     true,
	}, nil
}

// Close implements oas.Handler.
func (b *BotAPI) Close(ctx context.Context) (oas.Result, error) {
	b.pool.Kill(MustToken(ctx))
	return resultOK(true), nil
}

// LogOut implements oas.Handler.
func (b *BotAPI) LogOut(ctx context.Context) (oas.Result, error) {
	var r bool
	if err := b.do(ctx, func(client *telegram.Client) (err error) {
		r, err = client.API().AuthLogOut(ctx)
		return err
	}); err != nil {
		return oas.Result{}, err
	}

	return resultOK(r), &NotImplementedError{}
}
