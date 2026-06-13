package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// errNotInPrivateChat is returned by chat-management methods invoked on a
// private (user) chat where they do not apply.
func errNotInPrivateChat() *Error {
	return &Error{Code: 400, Description: "Bad Request: method is not available in private chats"}
}

// chatAdmin is the subset of peer operations shared by basic groups and
// channels/supergroups. peers.Chat and peers.Channel both implement it; users
// do not.
type chatAdmin interface {
	SetTitle(ctx context.Context, title string) error
	SetDescription(ctx context.Context, about string) error
	Leave(ctx context.Context) error
}

// resolveChatAdmin resolves a ChatID to the shared admin operations, rejecting
// private chats.
func (b *Bot) resolveChatAdmin(ctx context.Context, chat ChatID) (chatAdmin, error) {
	p, err := b.resolvePeer(ctx, chat)
	if err != nil {
		return nil, err
	}
	a, ok := p.(chatAdmin)
	if !ok {
		return nil, errNotInPrivateChat()
	}
	return a, nil
}

// SetChatTitle changes the title of a chat. The bot must be an administrator
// with the appropriate rights.
func (b *Bot) SetChatTitle(ctx context.Context, chat ChatID, title string) error {
	a, err := b.resolveChatAdmin(ctx, chat)
	if err != nil {
		return err
	}
	if err := a.SetTitle(ctx, title); err != nil {
		return asAPIError(err)
	}
	return nil
}

// SetChatDescription changes the description of a chat.
func (b *Bot) SetChatDescription(ctx context.Context, chat ChatID, description string) error {
	a, err := b.resolveChatAdmin(ctx, chat)
	if err != nil {
		return err
	}
	if err := a.SetDescription(ctx, description); err != nil {
		return asAPIError(err)
	}
	return nil
}

// LeaveChat makes the bot leave a group, supergroup or channel.
func (b *Bot) LeaveChat(ctx context.Context, chat ChatID) error {
	a, err := b.resolveChatAdmin(ctx, chat)
	if err != nil {
		return err
	}
	if err := a.Leave(ctx); err != nil {
		return asAPIError(err)
	}
	return nil
}

// PinChatMessage pins a message in a chat. By default it notifies members; pass
// Silent to pin quietly.
func (b *Bot) PinChatMessage(ctx context.Context, chat ChatID, messageID int, opts ...SendOption) error {
	var cfg sendConfig
	for _, o := range opts {
		o(&cfg)
	}
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}
	if _, err := b.raw.MessagesUpdatePinnedMessage(ctx, &tg.MessagesUpdatePinnedMessageRequest{
		Peer:   peer,
		ID:     messageID,
		Silent: cfg.silent,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// UnpinChatMessage unpins the message with the given id in a chat.
func (b *Bot) UnpinChatMessage(ctx context.Context, chat ChatID, messageID int) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}
	if _, err := b.raw.MessagesUpdatePinnedMessage(ctx, &tg.MessagesUpdatePinnedMessageRequest{
		Peer:  peer,
		ID:    messageID,
		Unpin: true,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// UnpinAllChatMessages clears the list of pinned messages in a chat.
func (b *Bot) UnpinAllChatMessages(ctx context.Context, chat ChatID) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}
	if _, err := b.raw.MessagesUnpinAllMessages(ctx, &tg.MessagesUnpinAllMessagesRequest{
		Peer: peer,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// SetChatPermissions sets the default permissions for all non-administrator
// members of a supergroup. permissions is an allow-list.
func (b *Bot) SetChatPermissions(ctx context.Context, chat ChatID, permissions ChatPermissions) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}
	if _, err := b.raw.MessagesEditChatDefaultBannedRights(ctx, &tg.MessagesEditChatDefaultBannedRightsRequest{
		Peer:         peer,
		BannedRights: permissions.toBannedRights(0),
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// SetChatStickerSet sets the group sticker set for a supergroup.
func (b *Bot) SetChatStickerSet(ctx context.Context, chat ChatID, stickerSetName string) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}
	if _, err := b.raw.ChannelsSetStickers(ctx, &tg.ChannelsSetStickersRequest{
		Channel:    channel,
		Stickerset: &tg.InputStickerSetShortName{ShortName: stickerSetName},
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// DeleteChatStickerSet removes the group sticker set from a supergroup.
func (b *Bot) DeleteChatStickerSet(ctx context.Context, chat ChatID) error {
	channel, err := b.resolveChannel(ctx, chat)
	if err != nil {
		return err
	}
	if _, err := b.raw.ChannelsSetStickers(ctx, &tg.ChannelsSetStickersRequest{
		Channel:    channel,
		Stickerset: &tg.InputStickerSetEmpty{},
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}
