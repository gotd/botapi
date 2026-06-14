package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message/rich"
	"github.com/gotd/td/tg"
)

// SendRichMessage sends a rich message (Bot API 10.1): structured content —
// headings, paragraphs, lists, tables, block quotes, media, math and more — as
// a tree of page blocks rather than a flat string with entity ranges.
//
// Build the content with github.com/gotd/td/telegram/message/rich, for example:
//
//	msg := rich.New(
//		rich.Heading1(rich.Plain("Title")),
//		rich.Paragraph(rich.Bold(rich.Plain("Hello"))),
//	).Input()
//	bot.SendRichMessage(ctx, chat, msg)
//
// For whole-document HTML or Markdown, prefer SendRichHTML / SendRichMarkdown.
// The usual SendOptions (reply, silent, protect, reply markup) apply.
func (b *Bot) SendRichMessage(ctx context.Context, chat ChatID, msg tg.InputRichMessageClass, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	builder := &b.sender.To(peer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.RichMessage(ctx, msg)

	return b.sentMessage(ctx, peer, resp, err)
}

// SendRichHTML sends a rich message whose content is the given HTML document,
// parsed by Telegram's servers.
func (b *Bot) SendRichHTML(ctx context.Context, chat ChatID, html string, opts ...SendOption) (*Message, error) {
	return b.SendRichMessage(ctx, chat, rich.HTML(html), opts...)
}

// SendRichMarkdown sends a rich message whose content is the given Markdown
// document, parsed by Telegram's servers.
func (b *Bot) SendRichMarkdown(ctx context.Context, chat ChatID, markdown string, opts ...SendOption) (*Message, error) {
	return b.SendRichMessage(ctx, chat, rich.Markdown(markdown), opts...)
}
