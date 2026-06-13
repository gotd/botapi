// Command rich is a showcase of Telegram rich messages (Bot API 10.1) built with
// github.com/gotd/td/telegram/message/rich. Each command sends a rich message
// exercising a category of the allowed rich-text and page-block constructors:
//
//	/headings  headings 1-6 and a footer
//	/text      every allowed inline text style (bold, italic, spoiler, math, …)
//	/lists     bullet, checklist and ordered lists (incl. block items)
//	/blocks    paragraph, anchor, code, quotes, math, details, thinking, divider
//	/table     a table with header and data cells
//	/map       a map block
//	/media     photo/audio/image/custom-emoji/collage/slideshow (need real ids)
//
// Not every page block is valid in a bot-sent rich message. Per the official Bot
// API server (telegram-bot-api, td/td/telegram/WebPageBlock.cpp), these are
// Instant-View-page-only and the server rejects them with
// RICH_VALIDATE_CTOR_NOT_ALLOWED, so they are intentionally not used here:
// Title, Subtitle, Header, Subheader, Kicker, AuthorDate, Cover,
// RelatedArticles — and the inline auto-link styles (AutoURL/AutoEmail/
// AutoPhone). Video is omitted too.
//
// Run it with an MTProto app identity (https://my.telegram.org) and a BotFather
// token. To exercise /media's media-by-id blocks, also set any of PHOTO_ID,
// AUDIO_ID, DOCUMENT_ID, EMOJI_ID to real Telegram resource ids:
//
//	APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/rich
package main

import (
	"context"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gotd/log/logzap"
	"go.uber.org/zap"

	"github.com/gotd/td/telegram/message/rich"
	"github.com/gotd/td/tg"

	"github.com/gotd/botapi"
)

func main() {
	log, _ := zap.NewProduction()
	defer func() { _ = log.Sync() }()

	appID, err := strconv.Atoi(os.Getenv("APP_ID"))
	if err != nil {
		log.Fatal("APP_ID must be a number (see https://my.telegram.org)", zap.Error(err))
	}

	bot, err := botapi.New(os.Getenv("BOT_TOKEN"), botapi.Options{
		AppID:   appID,
		AppHash: os.Getenv("APP_HASH"),
		Logger:  logzap.New(log),
	})
	if err != nil {
		log.Fatal("Create bot", zap.Error(err))
	}
	bot.Use(botapi.Recover(), botapi.Logging())

	bot.OnCommand("start", "List the rich-message demos", func(c *botapi.Context) error {
		return send(c,
			rich.Heading1(rich.Plain("Rich message showcase")),
			rich.Paragraph(rich.Plain("Try these commands to see the page-block constructors:")),
			rich.List(
				rich.ListItem(rich.Fixed(rich.Plain("/headings"))),
				rich.ListItem(rich.Fixed(rich.Plain("/text"))),
				rich.ListItem(rich.Fixed(rich.Plain("/lists"))),
				rich.ListItem(rich.Fixed(rich.Plain("/blocks"))),
				rich.ListItem(rich.Fixed(rich.Plain("/table"))),
				rich.ListItem(rich.Fixed(rich.Plain("/map"))),
				rich.ListItem(rich.Fixed(rich.Plain("/media"))),
			),
		)
	})

	bot.OnCommand("headings", "Headings and a footer", func(c *botapi.Context) error {
		return send(c,
			rich.Heading1(rich.Plain("Heading 1")),
			rich.Heading2(rich.Plain("Heading 2")),
			rich.Heading3(rich.Plain("Heading 3")),
			rich.Heading4(rich.Plain("Heading 4")),
			rich.Heading5(rich.Plain("Heading 5")),
			rich.Heading6(rich.Plain("Heading 6")),
			rich.Heading(3, rich.Plain("Heading via Heading(level)")),
			rich.Footer(rich.Plain("— the footer")),
		)
	})

	bot.OnCommand("text", "Every allowed inline text style", func(c *botapi.Context) error {
		return send(c,
			rich.Paragraph(rich.Concat(
				rich.Bold(rich.Plain("bold")), sp(), rich.Italic(rich.Plain("italic")), sp(),
				rich.Underline(rich.Plain("underline")), sp(), rich.Strike(rich.Plain("strike")), sp(),
				rich.Spoiler(rich.Plain("spoiler")), sp(), rich.Marked(rich.Plain("marked")), sp(),
				rich.Fixed(rich.Plain("fixed")),
			)),
			rich.Paragraph(rich.Concat(
				rich.Plain("x"), rich.Subscript(rich.Plain("i")),
				rich.Plain(" and e"), rich.Superscript(rich.Plain("iπ")),
				rich.Plain(", inline math "), rich.Math("a^2+b^2=c^2"),
			)),
			rich.Paragraph(rich.Concat(
				rich.Mention(rich.Plain("@durov")), sp(),
				rich.MentionName(rich.Plain("Ada"), 1), sp(),
				rich.Hashtag(rich.Plain("#golang")), sp(),
				rich.Cashtag(rich.Plain("$TON")), sp(),
				rich.BotCommand(rich.Plain("/start")), sp(),
				rich.BankCard(rich.Plain("4111111111111111")),
			)),
			rich.Paragraph(rich.Concat(
				rich.Email(rich.Plain("hi@example.com"), "hi@example.com"), sp(),
				rich.Phone(rich.Plain("+1 555 0100"), "+15550100"), sp(),
				rich.URL(rich.Plain("gotd/td"), "https://github.com/gotd/td", 0), sp(),
				rich.AnchorLink(rich.Plain("jump to anchor"), "more"),
			)),
			rich.Paragraph(rich.Concat(
				rich.Plain("date: "),
				rich.Date(rich.Plain("now"), int(time.Now().Unix()), rich.DateFlags{LongDate: true, ShortTime: true}),
				rich.Plain(", empty: ["), rich.Empty(), rich.Plain("]"),
			)),
		)
	})

	bot.OnCommand("lists", "Bullet, checklist and ordered lists", func(c *botapi.Context) error {
		return send(c,
			rich.Heading3(rich.Plain("Unordered")),
			rich.List(
				rich.ListItem(rich.Plain("a plain item")),
				rich.CheckListItem(true, rich.Plain("a checked item")),
				rich.CheckListItem(false, rich.Plain("an unchecked item")),
				rich.ListItemBlocks(
					rich.Paragraph(rich.Plain("an item made of blocks:")),
					rich.Preformatted(rich.Plain("nested := true"), "go"),
				),
			),
			rich.Heading3(rich.Plain("Ordered")),
			rich.OrderedList(
				rich.OrderedListItem("1.", rich.Plain("first")),
				rich.OrderedListItem("2.", rich.Plain("second")),
				rich.OrderedListItemBlocks("3.", rich.Paragraph(rich.Plain("third, with blocks"))),
			),
		)
	})

	bot.OnCommand("blocks", "Anchor, code, quotes, math, details, divider", func(c *botapi.Context) error {
		return send(c,
			rich.AnchorBlock("more"),
			rich.Paragraph(rich.Concat(
				rich.Anchor(rich.Plain("An inline anchor"), "inline-anchor"),
				rich.Plain(" sits inside the text."),
			)),
			rich.Preformatted(rich.Plain("func main() {\n\tfmt.Println(\"hi\")\n}"), "go"),
			rich.Blockquote(rich.Plain("A block quotation."), rich.Plain("— a credit")),
			rich.BlockquoteBlocks(rich.Plain("— a multi-block quote"),
				rich.Paragraph(rich.Plain("First quoted paragraph.")),
				rich.Paragraph(rich.Plain("Second quoted paragraph.")),
			),
			rich.Pullquote(rich.Plain("A pull quote stands out."), rich.Empty()),
			rich.MathBlock("\\int_0^1 x^2\\,dx = \\tfrac{1}{3}"),
			rich.Details(true, rich.Plain("Expandable details"),
				rich.Paragraph(rich.Plain("Hidden until expanded.")),
			),
			rich.Thinking(rich.Plain("…thinking block…")),
			rich.Divider(),
			rich.Paragraph(rich.Plain("(below the divider)")),
		)
	})

	bot.OnCommand("table", "A table", func(c *botapi.Context) error {
		return send(c,
			rich.Table(rich.Plain("High scores"),
				rich.Row(rich.HeaderCell(rich.Plain("Name")), rich.HeaderCell(rich.Plain("Score"))),
				rich.Row(rich.Cell(rich.Plain("Ada")), rich.Cell(rich.Plain("100"))),
				rich.Row(rich.Cell(rich.Plain("Bob")), rich.Cell(rich.Plain("80"))),
			),
		)
	})

	bot.OnCommand("map", "A map block", func(c *botapi.Context) error {
		return send(c,
			rich.Paragraph(rich.Plain("Red Square, Moscow:")),
			rich.Map(&tg.InputGeoPoint{Lat: 55.7539, Long: 37.6208}, 15, 600, 400, caption("Red Square")),
		)
	})

	bot.OnCommand("media", "Media blocks (set *_ID env vars for real media)", func(c *botapi.Context) error {
		return send(c, mediaBlocks()...)
	})

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	log.Info("Starting rich bot")
	if err := bot.Run(ctx); err != nil {
		log.Fatal("Run", zap.Error(err))
	}
}

// mediaBlocks builds the media-by-id page blocks. The blocks that reference a
// Telegram resource are included only when the matching env var holds a real id,
// so the message always sends; the constructors are exercised regardless.
func mediaBlocks() []tg.PageBlockClass {
	blocks := []tg.PageBlockClass{rich.Heading2(rich.Plain("Media blocks"))}

	if id := envID("PHOTO_ID"); id != 0 {
		blocks = append(blocks,
			rich.Photo(id, caption("A photo")),
			rich.Collage(caption("A collage"),
				rich.Photo(id, emptyCaption()), rich.Photo(id, emptyCaption())),
			rich.Slideshow(caption("A slideshow"),
				rich.Photo(id, emptyCaption()), rich.Photo(id, emptyCaption())),
		)
	}
	if id := envID("AUDIO_ID"); id != 0 {
		blocks = append(blocks, rich.Audio(id, caption("An audio")))
	}
	if id := envID("DOCUMENT_ID"); id != 0 {
		blocks = append(blocks, rich.Paragraph(rich.Concat(
			rich.Plain("inline image: "), rich.Image(id, 24, 24),
		)))
	}
	if id := envID("EMOJI_ID"); id != 0 {
		blocks = append(blocks, rich.Paragraph(rich.Concat(
			rich.Plain("custom emoji: "), rich.CustomEmoji(id, "👍"),
		)))
	}

	// A map always renders, so /media works even without any media ids.
	blocks = append(blocks,
		rich.Map(&tg.InputGeoPoint{Lat: 55.7539, Long: 37.6208}, 15, 600, 400, caption("Red Square")),
	)
	return blocks
}

// send assembles the blocks into a rich message and sends it to the chat.
func send(c *botapi.Context, blocks ...tg.PageBlockClass) error {
	chat, ok := c.Chat()
	if !ok {
		return nil
	}
	_, err := c.Bot.SendRichMessage(c, chat, rich.New(blocks...).Input())
	return err
}

// sp is a single space text node, for separating inline runs.
func sp() tg.RichTextClass { return rich.Plain(" ") }

// caption builds a page caption with the given text and no credit.
func caption(text string) tg.PageCaption { return rich.Caption(rich.Plain(text), rich.Empty()) }

// emptyCaption is a caption with no text or credit.
func emptyCaption() tg.PageCaption { return rich.Caption(rich.Empty(), rich.Empty()) }

// envID parses an int64 id from an environment variable, or 0 if unset/invalid.
func envID(name string) int64 {
	id, _ := strconv.ParseInt(os.Getenv(name), 10, 64)
	return id
}
