package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// BanChatMember implements oas.Handler.
func (b *BotAPI) BanChatMember(ctx context.Context, req oas.BanChatMember) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// BanChatSenderChat implements oas.Handler.
func (b *BotAPI) BanChatSenderChat(ctx context.Context, req oas.BanChatSenderChat) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// GetChatAdministrators implements oas.Handler.
func (b *BotAPI) GetChatAdministrators(ctx context.Context, req oas.GetChatAdministrators) (oas.ResultArrayOfChatMember, error) {
	return oas.ResultArrayOfChatMember{}, &NotImplementedError{}
}

// GetChatMember implements oas.Handler.
func (b *BotAPI) GetChatMember(ctx context.Context, req oas.GetChatMember) (oas.ResultChatMember, error) {
	return oas.ResultChatMember{}, &NotImplementedError{}
}

// GetChatMemberCount implements oas.Handler.
func (b *BotAPI) GetChatMemberCount(ctx context.Context, req oas.GetChatMemberCount) (oas.ResultInt, error) {
	ch, err := b.resolveChatID(ctx, req.ChatID)
	if err != nil {
		return oas.ResultInt{}, err
	}
	return oas.ResultInt{
		Result: oas.NewOptInt(ch.ParticipantsCount()),
		Ok:     true,
	}, nil
}

// PromoteChatMember implements oas.Handler.
func (b *BotAPI) PromoteChatMember(ctx context.Context, req oas.PromoteChatMember) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// RestrictChatMember implements oas.Handler.
func (b *BotAPI) RestrictChatMember(ctx context.Context, req oas.RestrictChatMember) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// UnbanChatMember implements oas.Handler.
func (b *BotAPI) UnbanChatMember(ctx context.Context, req oas.UnbanChatMember) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}

// UnbanChatSenderChat implements oas.Handler.
func (b *BotAPI) UnbanChatSenderChat(ctx context.Context, req oas.UnbanChatSenderChat) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
