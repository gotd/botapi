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

// RichMessage is the structured content of a rich message, expressed as native
// MTProto page blocks.
//
// The HTTP Bot API accepts a rich message only as an HTML or Markdown string and
// parses it into this structure server-side. botapi is MTProto-native, so it
// takes the page-block tree directly: no content parser is needed and the full
// block vocabulary (paragraphs, headers, lists, embedded media, …) is available.
// Build blocks with the tg.PageBlock* types; Bot.Raw exposes the same client.
type RichMessage struct {
	// Blocks are the page blocks that make up the message body.
	Blocks []tg.PageBlockClass
	// Photos are photos referenced by the blocks.
	Photos []tg.PhotoClass
	// Documents are documents referenced by the blocks.
	Documents []tg.DocumentClass
	// RTL renders the message right-to-left.
	RTL bool
	// Part marks this as a partial segment of a longer streamed message.
	Part bool
}

// toTg converts the rich message into its MTProto representation.
func (r RichMessage) toTg() tg.RichMessage {
	out := tg.RichMessage{
		Blocks:    r.Blocks,
		Photos:    r.Photos,
		Documents: r.Documents,
	}
	if r.RTL {
		out.SetRtl(true)
	}

	if r.Part {
		out.SetPart(true)
	}

	return out
}

// SendRichMessageDraft streams a partial rich message to a chat while the message
// is being generated, so its members see an ephemeral preview. draftID is a
// unique client-side identifier for the draft. Once generation finishes, persist
// the result with a regular send.
func (b *Bot) SendRichMessageDraft(ctx context.Context, chat ChatID, draftID int64, message RichMessage, opts ...DraftOption) error {
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
		Action: &tg.SendMessageRichMessageDraftAction{
			RandomID:    draftID,
			RichMessage: message.toTg(),
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
