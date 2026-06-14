package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// hideJoinRequest approves or declines a pending join request for a chat.
func (b *Bot) hideJoinRequest(ctx context.Context, chat ChatID, userID int64, approved bool) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	if _, err := b.raw.MessagesHideChatJoinRequest(ctx, &tg.MessagesHideChatJoinRequestRequest{
		Approved: approved,
		Peer:     peer,
		UserID:   user,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// ApproveChatJoinRequest approves a chat join request. The bot must be an
// administrator with the can_invite_users right.
func (b *Bot) ApproveChatJoinRequest(ctx context.Context, chat ChatID, userID int64) error {
	return b.hideJoinRequest(ctx, chat, userID, true)
}

// DeclineChatJoinRequest declines a chat join request. The bot must be an
// administrator with the can_invite_users right.
func (b *Bot) DeclineChatJoinRequest(ctx context.Context, chat ChatID, userID int64) error {
	return b.hideJoinRequest(ctx, chat, userID, false)
}
