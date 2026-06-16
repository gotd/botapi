package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// topicCreateService extracts the service message that a forum-topic creation
// produced from the RPC's Updates response.
func topicCreateService(resp tg.UpdatesClass) (*tg.MessageService, bool) {
	var ups []tg.UpdateClass

	switch u := resp.(type) {
	case *tg.Updates:
		ups = u.Updates
	case *tg.UpdatesCombined:
		ups = u.Updates
	default:
		return nil, false
	}

	for _, up := range ups {
		var msg tg.MessageClass

		switch m := up.(type) {
		case *tg.UpdateNewChannelMessage:
			msg = m.Message
		case *tg.UpdateNewMessage:
			msg = m.Message
		default:
			continue
		}

		if svc, ok := msg.(*tg.MessageService); ok {
			return svc, true
		}
	}

	return nil, false
}

// generalForumTopicID is the topic id of the special "General" topic that every
// forum supergroup has. The general-topic methods operate on it.
const generalForumTopicID = 1

// ForumTopic is the result of CreateForumTopic: a topic in a forum supergroup.
type ForumTopic struct {
	// MessageThreadID is the unique identifier of the forum topic.
	MessageThreadID int `json:"message_thread_id"`
	// Name is the topic name.
	Name string `json:"name"`
	// IconColor is the color of the topic icon in RGB format.
	IconColor int `json:"icon_color"`
	// IconCustomEmojiID is the unique identifier of the custom emoji shown as the
	// topic icon, if any.
	IconCustomEmojiID string `json:"icon_custom_emoji_id,omitempty"`
}

// forumTopicConfig accumulates the optional parameters shared by CreateForumTopic
// and EditForumTopic.
type forumTopicConfig struct {
	name         string
	nameSet      bool
	iconColor    int
	iconColorSet bool
	iconEmoji    string
	iconEmojiSet bool
}

// ForumTopicOption configures CreateForumTopic / EditForumTopic.
type ForumTopicOption func(*forumTopicConfig)

// WithForumTopicName sets the topic name. Used by EditForumTopic; CreateForumTopic
// takes the name as a positional argument.
func WithForumTopicName(name string) ForumTopicOption {
	return func(c *forumTopicConfig) {
		c.name = name
		c.nameSet = true
	}
}

// WithForumTopicIconColor sets the color of the topic icon in RGB format. One of
// 0x6FB9F0, 0xFFD67E, 0xCB86DB, 0x8EEE98, 0xFF93B2, or 0xFB6F5F. Only honored on
// creation.
func WithForumTopicIconColor(color int) ForumTopicOption {
	return func(c *forumTopicConfig) {
		c.iconColor = color
		c.iconColorSet = true
	}
}

// WithForumTopicIconCustomEmojiID sets the custom emoji shown as the topic icon.
// An empty string removes the icon (EditForumTopic only).
func WithForumTopicIconCustomEmojiID(id string) ForumTopicOption {
	return func(c *forumTopicConfig) {
		c.iconEmoji = id
		c.iconEmojiSet = true
	}
}

// parseCustomEmojiID turns a Bot API custom_emoji_id string into the MTProto
// document id. An empty string maps to 0 (no/removed icon).
func parseCustomEmojiID(id string) (int64, error) {
	if id == "" {
		return 0, nil
	}

	v, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return 0, &Error{Code: 400, Description: "Bad Request: invalid icon_custom_emoji_id"}
	}

	return v, nil
}

// CreateForumTopic creates a topic in a forum supergroup. The bot must be an
// administrator with the can_manage_topics right. Returns the created topic.
func (b *Bot) CreateForumTopic(ctx context.Context, chat ChatID, name string, opts ...ForumTopicOption) (*ForumTopic, error) {
	var cfg forumTopicConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	randomID, err := randInt64()
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesCreateForumTopicRequest{
		Peer:     peer,
		Title:    name,
		RandomID: randomID,
	}

	if cfg.iconColorSet {
		req.SetIconColor(cfg.iconColor)
	}

	if cfg.iconEmojiSet {
		emojiID, err := parseCustomEmojiID(cfg.iconEmoji)
		if err != nil {
			return nil, err
		}

		req.SetIconEmojiID(emojiID)
	}

	resp, err := b.raw.MessagesCreateForumTopic(ctx, req)
	if err != nil {
		return nil, asAPIError(err)
	}

	// The topic id is the id of the service message that created it; its
	// MessageActionTopicCreate carries the resolved name/color/icon.
	topic := &ForumTopic{Name: name, IconColor: cfg.iconColor, IconCustomEmojiID: cfg.iconEmoji}

	if svc, ok := topicCreateService(resp); ok {
		topic.MessageThreadID = svc.ID

		if act, ok := svc.Action.(*tg.MessageActionTopicCreate); ok {
			topic.Name = act.Title
			topic.IconColor = act.IconColor

			if act.IconEmojiID != 0 {
				topic.IconCustomEmojiID = strconv.FormatInt(act.IconEmojiID, 10)
			}
		}
	}

	return topic, nil
}

// editForumTopic is the shared implementation behind the (re)naming, icon,
// close/reopen and hide/unhide operations on a topic.
func (b *Bot) editForumTopic(ctx context.Context, chat ChatID, topicID int, apply func(*tg.MessagesEditForumTopicRequest) error) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	req := &tg.MessagesEditForumTopicRequest{Peer: peer, TopicID: topicID}

	if err := apply(req); err != nil {
		return err
	}

	if _, err := b.raw.MessagesEditForumTopic(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// EditForumTopic edits the name and/or icon of a topic in a forum supergroup. The
// bot must be an administrator with the can_manage_topics right.
func (b *Bot) EditForumTopic(ctx context.Context, chat ChatID, messageThreadID int, opts ...ForumTopicOption) error {
	var cfg forumTopicConfig

	for _, o := range opts {
		o(&cfg)
	}

	return b.editForumTopic(ctx, chat, messageThreadID, func(req *tg.MessagesEditForumTopicRequest) error {
		if cfg.nameSet {
			req.SetTitle(cfg.name)
		}

		if cfg.iconEmojiSet {
			emojiID, err := parseCustomEmojiID(cfg.iconEmoji)
			if err != nil {
				return err
			}

			req.SetIconEmojiID(emojiID)
		}

		return nil
	})
}

// CloseForumTopic closes an open topic in a forum supergroup.
func (b *Bot) CloseForumTopic(ctx context.Context, chat ChatID, messageThreadID int) error {
	return b.editForumTopic(ctx, chat, messageThreadID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetClosed(true)

		return nil
	})
}

// ReopenForumTopic reopens a closed topic in a forum supergroup.
func (b *Bot) ReopenForumTopic(ctx context.Context, chat ChatID, messageThreadID int) error {
	return b.editForumTopic(ctx, chat, messageThreadID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetClosed(false)

		return nil
	})
}

// DeleteForumTopic deletes a topic in a forum supergroup together with all its
// messages. The bot must be an administrator with the can_delete_messages right.
func (b *Bot) DeleteForumTopic(ctx context.Context, chat ChatID, messageThreadID int) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	if _, err := b.raw.MessagesDeleteTopicHistory(ctx, &tg.MessagesDeleteTopicHistoryRequest{
		Peer:     peer,
		TopMsgID: messageThreadID,
	}); err != nil {
		return asAPIError(err)
	}

	return nil
}

// UnpinAllForumTopicMessages clears the list of pinned messages in a forum topic.
func (b *Bot) UnpinAllForumTopicMessages(ctx context.Context, chat ChatID, messageThreadID int) error {
	return b.unpinAllTopicMessages(ctx, chat, messageThreadID)
}

// unpinAllTopicMessages unpins every message in the given topic.
func (b *Bot) unpinAllTopicMessages(ctx context.Context, chat ChatID, topicID int) error {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	req := &tg.MessagesUnpinAllMessagesRequest{Peer: peer}
	req.SetTopMsgID(topicID)

	if _, err := b.raw.MessagesUnpinAllMessages(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// EditGeneralForumTopic edits the name of the "General" topic in a forum
// supergroup. The bot must be an administrator with the can_manage_topics right.
func (b *Bot) EditGeneralForumTopic(ctx context.Context, chat ChatID, name string) error {
	return b.editForumTopic(ctx, chat, generalForumTopicID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetTitle(name)

		return nil
	})
}

// CloseGeneralForumTopic closes the "General" topic in a forum supergroup.
func (b *Bot) CloseGeneralForumTopic(ctx context.Context, chat ChatID) error {
	return b.editForumTopic(ctx, chat, generalForumTopicID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetClosed(true)

		return nil
	})
}

// ReopenGeneralForumTopic reopens the "General" topic in a forum supergroup.
func (b *Bot) ReopenGeneralForumTopic(ctx context.Context, chat ChatID) error {
	return b.editForumTopic(ctx, chat, generalForumTopicID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetClosed(false)

		return nil
	})
}

// HideGeneralForumTopic hides the "General" topic in a forum supergroup. The
// topic is automatically closed if it was open.
func (b *Bot) HideGeneralForumTopic(ctx context.Context, chat ChatID) error {
	return b.editForumTopic(ctx, chat, generalForumTopicID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetHidden(true)

		return nil
	})
}

// UnhideGeneralForumTopic unhides the "General" topic in a forum supergroup.
func (b *Bot) UnhideGeneralForumTopic(ctx context.Context, chat ChatID) error {
	return b.editForumTopic(ctx, chat, generalForumTopicID, func(req *tg.MessagesEditForumTopicRequest) error {
		req.SetHidden(false)

		return nil
	})
}

// UnpinAllGeneralForumTopicMessages clears the list of pinned messages in the
// "General" topic of a forum supergroup.
func (b *Bot) UnpinAllGeneralForumTopicMessages(ctx context.Context, chat ChatID) error {
	return b.unpinAllTopicMessages(ctx, chat, generalForumTopicID)
}

// GetForumTopicIconStickers returns custom emoji stickers that can be used as a
// forum topic icon by any user.
func (b *Bot) GetForumTopicIconStickers(ctx context.Context) ([]Sticker, error) {
	res, err := b.raw.MessagesGetStickerSet(ctx, &tg.MessagesGetStickerSetRequest{
		Stickerset: &tg.InputStickerSetEmojiDefaultTopicIcons{},
		Hash:       0,
	})
	if err != nil {
		return nil, asAPIError(err)
	}

	set, ok := res.(*tg.MessagesStickerSet)
	if !ok {
		return nil, &Error{Code: 400, Description: "Bad Request: sticker set not found"}
	}

	out := make([]Sticker, 0, len(set.Documents))

	for _, d := range set.Documents {
		if doc, ok := d.(*tg.Document); ok {
			out = append(out, stickerFromDocument(doc, set.Set.ShortName, StickerCustomEmoji))
		}
	}

	return out, nil
}
