package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// editText is the shared path for editing a message's text or caption.
func (b *Bot) editText(ctx context.Context, chat ChatID, messageID int, text string, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	styled, err := styledText(text, cfg.parseMode, b.peers.UserResolveHook(ctx))
	if err != nil {
		return nil, err
	}

	builder := &b.sender.To(peer).Builder
	if cfg.disableWebPreview {
		builder = builder.NoWebpage()
	}

	if cfg.markup != nil {
		mkp, err := replyMarkupToTg(cfg.markup)
		if err != nil {
			return nil, err
		}

		builder = builder.Markup(mkp)
	}

	resp, err := builder.Edit(messageID).StyledText(ctx, styled...)

	return b.sentMessage(ctx, peer, resp, err)
}

// EditMessageText edits the text of a message and returns the edited message.
func (b *Bot) EditMessageText(ctx context.Context, chat ChatID, messageID int, text string, opts ...SendOption) (*Message, error) {
	return b.editText(ctx, chat, messageID, text, opts...)
}

// EditMessageCaption edits the caption of a media message and returns the
// edited message.
func (b *Bot) EditMessageCaption(ctx context.Context, chat ChatID, messageID int, caption string, opts ...SendOption) (*Message, error) {
	return b.editText(ctx, chat, messageID, caption, opts...)
}

// EditMessageReplyMarkup edits only the reply markup of a message, leaving its
// text and media unchanged. A nil markup removes the keyboard.
func (b *Bot) EditMessageReplyMarkup(ctx context.Context, chat ChatID, messageID int, markup ReplyMarkup) (*Message, error) {
	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	req := &tg.MessagesEditMessageRequest{Peer: peer, ID: messageID}

	if markup != nil {
		mkp, err := replyMarkupToTg(markup)
		if err != nil {
			return nil, err
		}

		req.SetReplyMarkup(mkp)
	}

	resp, err := b.raw.MessagesEditMessage(ctx, req)

	return b.sentMessage(ctx, peer, resp, err)
}

// editedMessageFromResp extracts the edited message from an edit response.
// Edits return UpdateEditMessage/UpdateEditChannelMessage, which the generic
// send unpacker does not handle.
func editedMessageFromResp(resp tg.UpdatesClass) (tg.MessageClass, bool) {
	var updates []tg.UpdateClass

	switch v := resp.(type) {
	case *tg.Updates:
		updates = v.Updates
	case *tg.UpdatesCombined:
		updates = v.Updates
	case *tg.UpdateShort:
		updates = []tg.UpdateClass{v.Update}
	default:
		return nil, false
	}

	for _, u := range updates {
		switch upd := u.(type) {
		case *tg.UpdateEditMessage:
			return upd.Message, true
		case *tg.UpdateEditChannelMessage:
			return upd.Message, true
		}
	}

	return nil, false
}
