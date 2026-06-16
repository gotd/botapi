package main

import (
	"strings"
	"time"

	"github.com/gotd/td/telegram/message/rich"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi"
)

// registerStreaming wires the two draft-streaming flows — the pattern an AI bot
// uses to show a live preview while it generates output. A draft is broadcast
// with a *Draft method under a stable draftID; reusing that id across the stream
// replaces the same ephemeral (~30s) preview instead of starting a new one.
// Nothing persists until the matching non-draft send is called at the end.
func registerStreaming(bot *botapi.Bot) {
	bot.OnCommand("stream", "Stream a rich message, then persist it", streamRich)
	bot.OnCommand("streamtext", "Stream a plain-text message, then persist it", streamText)
}

// streamRich grows a rich message (Bot API 10.1) one page block at a time,
// previewing each step with SendRichMessageDraft and persisting the finished
// tree with SendRichMessage.
func streamRich(c *botapi.Context) error {
	chat, ok := c.Chat()
	if !ok {
		return nil
	}

	// A non-zero id identifying this draft; reuse it across the whole stream.
	draftID := time.Now().UnixNano()

	// The body grows block by block, as if generated token by token.
	body := []tg.PageBlockClass{rich.Heading2(rich.Plain("Streaming a rich message"))}

	for _, part := range []string{
		"First, the bot drafts an opening paragraph…",
		"…then adds detail: each SendRichMessageDraft call updates the same preview.",
		"Finally it wraps up. Nothing persists until the finished message is sent.",
	} {
		body = append(body, rich.Paragraph(rich.Plain(part)))

		// rich.New(...).Input() yields the same tg.InputRichMessageClass that
		// SendRichMessage takes; rich.HTML / rich.Markdown work here too.
		if err := c.Bot.SendRichMessageDraft(c, chat, draftID, rich.New(body...).Input()); err != nil {
			return err
		}

		// Stand in for generation latency. The preview is ephemeral (~30s).
		time.Sleep(800 * time.Millisecond)
	}

	// Generation finished: persist the complete message.
	body = append(body, rich.Footer(rich.Plain("sent with SendRichMessage")))

	_, err := c.Bot.SendRichMessage(c, chat, rich.New(body...).Input())

	return err
}

// streamText is the plain-text counterpart: it appends to a growing string,
// previewing each step with SendMessageDraft (an HTML draft) and persisting the
// final text with SendMessage. SendMessageDraft returns no message — the draft
// is an ephemeral preview, not a sent message.
func streamText(c *botapi.Context) error {
	chat, ok := c.Chat()
	if !ok {
		return nil
	}

	draftID := time.Now().UnixNano()

	var b strings.Builder

	for word := range strings.FieldsSeq("Streaming a plain text message word by word…") {
		if b.Len() > 0 {
			b.WriteByte(' ')
		}

		b.WriteString(word)

		// Reuse draftID so each call updates the same preview. A parse mode (or
		// explicit entities) formats the draft just like a real message.
		if err := c.Bot.SendMessageDraft(c, chat, draftID, b.String()+" ▌",
			botapi.WithDraftParseMode(botapi.ParseModeHTML)); err != nil {
			return err
		}

		time.Sleep(400 * time.Millisecond)
	}

	// Generation finished: persist the final text without the cursor.
	_, err := c.Bot.SendMessage(c, chat, b.String())

	return err
}
