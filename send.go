package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/telegram/message/entity"
	"github.com/gotd/td/telegram/message/html"
	"github.com/gotd/td/telegram/message/markdown"
	"github.com/gotd/td/telegram/message/styling"
	"github.com/gotd/td/telegram/message/unpack"
	"github.com/gotd/td/tg"
)

// SendOption configures an outgoing send. Options are shared across the send
// methods; pass any combination.
type SendOption func(*sendConfig)

type sendConfig struct {
	disableWebPreview bool
	silent            bool
	protect           bool
	replyTo           int
	markup            ReplyMarkup
	parseMode         ParseMode
	businessConn      string
}

// DisableWebPagePreview disables the link preview for messages with links.
func DisableWebPagePreview() SendOption { return func(c *sendConfig) { c.disableWebPreview = true } }

// Silent sends the message without a notification sound.
func Silent() SendOption { return func(c *sendConfig) { c.silent = true } }

// ProtectContent prevents the message content from being forwarded or saved.
func ProtectContent() SendOption { return func(c *sendConfig) { c.protect = true } }

// ReplyTo makes the message a reply to the message with the given id.
func ReplyTo(messageID int) SendOption { return func(c *sendConfig) { c.replyTo = messageID } }

// WithReplyMarkup attaches an inline/reply keyboard (or removes one).
func WithReplyMarkup(m ReplyMarkup) SendOption { return func(c *sendConfig) { c.markup = m } }

// WithParseMode selects the formatting mode for the text or caption.
func WithParseMode(m ParseMode) SendOption { return func(c *sendConfig) { c.parseMode = m } }

// WithBusinessConnection sends the message on behalf of the business account
// identified by the given connection id, instead of as the bot itself.
//
// It applies to the text and media send methods (SendMessage, SendPhoto,
// SendDocument, SendVideo, SendAudio, SendMediaGroup, …). Files are uploaded in
// the bot's own session; only the resulting send is performed on behalf of the
// account.
func WithBusinessConnection(connectionID string) SendOption {
	return func(c *sendConfig) { c.businessConn = connectionID }
}

// senderFor returns the message sender the config selects: the bot's own sender,
// or a business-scoped sender when a business connection is set.
func (b *Bot) senderFor(cfg sendConfig) *message.Sender {
	if cfg.businessConn == "" {
		return b.sender
	}

	return b.businessSender(cfg.businessConn)
}

// styledText turns text + parse mode into styling options.
//
// The switch over ParseMode is exhaustive (exhaustive lint).
func styledText(text string, mode ParseMode, resolver entity.UserResolver) ([]styling.StyledTextOption, error) {
	switch mode {
	case ParseModeNone:
		return []styling.StyledTextOption{styling.Plain(text)}, nil
	case ParseModeHTML:
		return []styling.StyledTextOption{html.String(resolver, text)}, nil
	case ParseModeMarkdownV2, ParseModeMarkdown:
		return []styling.StyledTextOption{markdown.String(resolver, text)}, nil
	default:
		return nil, &Error{Code: 400, Description: "Bad Request: unsupported parse mode"}
	}
}

// applySendConfig applies the common builder options to a message builder.
func (b *Bot) applySendConfig(builder *message.Builder, cfg sendConfig) (*message.Builder, error) {
	if cfg.disableWebPreview {
		builder = builder.NoWebpage()
	}

	if cfg.silent {
		builder = builder.Silent()
	}

	if cfg.protect {
		builder = builder.NoForwards()
	}

	if cfg.replyTo != 0 {
		builder = builder.Reply(cfg.replyTo)
	}

	if cfg.markup != nil {
		mkp, err := replyMarkupToTg(cfg.markup)
		if err != nil {
			return nil, err
		}

		builder = builder.Markup(mkp)
	}

	return builder, nil
}

// sentMessage unpacks a send or edit response into a Bot API Message,
// backfilling the peer id when the server omitted it.
func (b *Bot) sentMessage(ctx context.Context, peer tg.InputPeerClass, resp tg.UpdatesClass, sendErr error) (*Message, error) {
	m, err := unpack.MessageClass(resp, sendErr)
	if err != nil {
		// unpack.MessageClass only handles new-message updates. Edits return
		// UpdateEditMessage/UpdateEditChannelMessage, which it rejects; extract
		// those here so EditMessage* don't report a failure for a successful edit.
		if edited, ok := editedMessageFromResp(resp); ok {
			m = edited
		} else {
			return nil, asAPIError(err)
		}
	}

	msg, ok := m.(*tg.Message)
	if !ok {
		return &Message{}, nil
	}

	if msg.PeerID == nil {
		switch p := peer.(type) {
		case *tg.InputPeerChat:
			msg.PeerID = &tg.PeerChat{ChatID: p.ChatID}
		case *tg.InputPeerUser:
			msg.PeerID = &tg.PeerUser{UserID: p.UserID}
		case *tg.InputPeerChannel:
			msg.PeerID = &tg.PeerChannel{ChannelID: p.ChannelID}
		}
	}

	return b.convertMessage(ctx, msg)
}

// SendMessage sends a text message to a chat and returns the sent message.
func (b *Bot) SendMessage(ctx context.Context, chat ChatID, text string, opts ...SendOption) (*Message, error) {
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

	builder := &b.senderFor(cfg).To(peer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.StyledText(ctx, styled...)

	return b.sentMessage(ctx, peer, resp, err)
}
