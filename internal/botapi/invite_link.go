package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// CreateChatInviteLink implements oas.Handler.
func (b *BotAPI) CreateChatInviteLink(ctx context.Context, req oas.CreateChatInviteLink) (oas.ResultChatInviteLink, error) {
	return oas.ResultChatInviteLink{}, &NotImplementedError{}
}

// EditChatInviteLink implements oas.Handler.
func (b *BotAPI) EditChatInviteLink(ctx context.Context, req oas.EditChatInviteLink) (oas.ResultChatInviteLink, error) {
	return oas.ResultChatInviteLink{}, &NotImplementedError{}
}

// ExportChatInviteLink implements oas.Handler.
func (b *BotAPI) ExportChatInviteLink(ctx context.Context, req oas.ExportChatInviteLink) (oas.ResultString, error) {
	return oas.ResultString{}, &NotImplementedError{}
}

// RevokeChatInviteLink implements oas.Handler.
func (b *BotAPI) RevokeChatInviteLink(ctx context.Context, req oas.RevokeChatInviteLink) (oas.ResultChatInviteLink, error) {
	return oas.ResultChatInviteLink{}, &NotImplementedError{}
}
