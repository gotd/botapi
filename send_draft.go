package botapi

import (
	"context"

	"github.com/gotd/td/tg"
)

// DraftOption configures SendMessageDraft.
type DraftOption func(*draftConfig)

type draftConfig struct {
	parseMode       ParseMode
	entities        []MessageEntity
	messageThreadID int
}

// WithDraftParseMode selects the formatting mode for the draft text.
func WithDraftParseMode(mode ParseMode) DraftOption {
	return func(c *draftConfig) { c.parseMode = mode }
}

// WithDraftEntities sets explicit entities for the draft text; they take
// precedence over a parse mode.
func WithDraftEntities(entities []MessageEntity) DraftOption {
	return func(c *draftConfig) { c.entities = entities }
}

// WithDraftMessageThread targets a specific forum topic / thread.
func WithDraftMessageThread(messageThreadID int) DraftOption {
	return func(c *draftConfig) { c.messageThreadID = messageThreadID }
}

// SendMessageDraft broadcasts a message draft to a chat, so its members see what
// the bot is preparing. draftID is a unique client-side identifier for the draft.
func (b *Bot) SendMessageDraft(ctx context.Context, chat ChatID, draftID int64, text string, opts ...DraftOption) error {
	var cfg draftConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	textWithEntities, err := b.textWithEntities(ctx, text, cfg.parseMode, cfg.entities)
	if err != nil {
		return err
	}

	req := &tg.MessagesSetTypingRequest{
		Peer: peer,
		Action: &tg.SendMessageTextDraftAction{
			RandomID: draftID,
			Text:     textWithEntities,
		},
	}
	if cfg.messageThreadID != 0 {
		req.SetTopMsgID(cfg.messageThreadID)
	}

	if _, err := b.raw.MessagesSetTyping(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}

// SendRichMessageDraft streams a partial rich message (Bot API 10.1) to a chat
// while the message is being generated, so its members see an ephemeral preview.
// draftID is a unique client-side identifier for the draft; reuse it across the
// stream and finalize with SendRichMessage once generation completes.
//
// Build msg the same way as for SendRichMessage — from native page blocks with
// github.com/gotd/td/telegram/message/rich (rich.New(...).Input()), or from a
// whole HTML/Markdown document (rich.HTML / rich.Markdown), which Telegram's
// servers parse.
func (b *Bot) SendRichMessageDraft(
	ctx context.Context, chat ChatID, draftID int64, msg tg.InputRichMessageClass, opts ...DraftOption,
) error {
	var cfg draftConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return err
	}

	req := &tg.MessagesSetTypingRequest{
		Peer: peer,
		Action: &tg.InputSendMessageRichMessageDraftAction{
			RandomID:    draftID,
			RichMessage: msg,
		},
	}
	if cfg.messageThreadID != 0 {
		req.SetTopMsgID(cfg.messageThreadID)
	}

	if _, err := b.raw.MessagesSetTyping(ctx, req); err != nil {
		return asAPIError(err)
	}

	return nil
}
