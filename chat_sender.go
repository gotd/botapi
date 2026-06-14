package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// BanChatSenderChat bans a channel chat from posting in a supergroup or channel
// on behalf of any of its sender chats. The bot must be an administrator with
// the appropriate rights. The ban is permanent until lifted with
// UnbanChatSenderChat.
func (b *Bot) BanChatSenderChat(ctx context.Context, chat ChatID, senderChatID int64) error {
	return b.editSenderChatBan(ctx, chat, senderChatID, tg.ChatBannedRights{ViewMessages: true})
}

// UnbanChatSenderChat lifts a previously set sender-chat ban in a supergroup or
// channel.
func (b *Bot) UnbanChatSenderChat(ctx context.Context, chat ChatID, senderChatID int64) error {
	return b.editSenderChatBan(ctx, chat, senderChatID, tg.ChatBannedRights{})
}

// editSenderChatban sets the banned rights of a sender chat (channel or user
// peer) within a supergroup or channel. Empty rights clear the ban.
func (b *Bot) editSenderChatBan(ctx context.Context, chat ChatID, senderChatID int64, rights tg.ChatBannedRights) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}

	sender, err := b.resolveInputPeer(ctx, ID(senderChatID))
	if err != nil {
		return err
	}

	if _, err := b.raw.ChannelsEditBanned(ctx, &tg.ChannelsEditBannedRequest{
		Channel:      channel,
		Participant:  sender,
		BannedRights: rights,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}
