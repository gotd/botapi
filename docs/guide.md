# botapi guide

A practical tour of `github.com/gotd/botapi` — a Telegram Bot API library
implemented over MTProto (via `gotd/td`) rather than HTTP to `api.telegram.org`.
For the why, see [`architecture.md`](./architecture.md); for status, see
[`roadmap.md`](./roadmap.md).

## Contents

- [Getting started](#getting-started)
- [Targeting chats](#targeting-chats)
- [Sending messages](#sending-messages)
- [Formatting](#formatting)
- [Keyboards](#keyboards)
- [Sending media](#sending-media)
- [Other sends](#other-sends)
- [Receiving updates](#receiving-updates)
- [Predicates](#predicates)
- [Middleware](#middleware)
- [Groups](#groups)
- [Commands](#commands)
- [Callback & inline queries](#callback--inline-queries)
- [Editing, forwarding, deleting](#editing-forwarding-deleting)
- [Files](#files)
- [Chat management](#chat-management)
- [Errors](#errors)
- [Resilience: flood-wait & rate limiting](#resilience-flood-wait--rate-limiting)
- [Persistence](#persistence)
- [Running many bots](#running-many-bots)
- [The escape hatch](#the-escape-hatch)

## Getting started

You need two things:

1. An **MTProto app identity** — `AppID` and `AppHash` from
   <https://my.telegram.org>. These identify the *application*, not the bot, and
   are required even for bots.
2. A **bot token** from [@BotFather](https://t.me/BotFather).

```go
bot, err := botapi.New(token, botapi.Options{AppID: appID, AppHash: appHash})
if err != nil {
	return err
}

bot.OnCommand("start", "Start the bot", func(c *botapi.Context) error {
	_, err := c.Reply("Hello!")
	return err
})

// Run connects, authorizes as a bot and serves updates until ctx is canceled.
return bot.Run(ctx)
```

`New` does no network I/O; register your handlers, then call `Run`.

## Targeting chats

Outgoing methods take a `ChatID`, a sealed union you build with `ID` (numeric)
or `Username`:

```go
botapi.ID(123456789)        // a numeric chat id
botapi.Username("@channel") // an @username (leading @ optional)
```

## Sending messages

Send methods hang off `*Bot`, take a `context.Context` first, a `ChatID`, and a
variadic of shared `SendOption`s:

```go
msg, err := bot.SendMessage(ctx, botapi.ID(chatID), "hi",
	botapi.ReplyTo(replyID),
	botapi.Silent(),
	botapi.DisableWebPagePreview(),
)
```

Common options (all `SendOption`): `ReplyTo`, `Silent`, `ProtectContent`,
`DisableWebPagePreview`, `WithReplyMarkup`, `WithParseMode`.

Inside a handler the `Context` shortcuts are usually enough:

```go
c.Send("text")               // send to the update's chat
c.Reply("text")              // reply to the incoming message
```

## Formatting

Pass `WithParseMode` with `ParseModeHTML`, `ParseModeMarkdownV2`, or the legacy
`ParseModeMarkdown`:

```go
bot.SendMessage(ctx, chat, "<b>bold</b> <i>italic</i>",
	botapi.WithParseMode(botapi.ParseModeHTML))

bot.SendMessage(ctx, chat, "*bold* _italic_ ||spoiler||",
	botapi.WithParseMode(botapi.ParseModeMarkdownV2))
```

## Keyboards

`ReplyMarkup` is a sealed union: `*InlineKeyboardMarkup`,
`*ReplyKeyboardMarkup`, `*ReplyKeyboardRemove`, `*ForceReply`. Build inline
keyboards with the helpers:

```go
kb := botapi.InlineKeyboard(
	[]botapi.InlineKeyboardButton{
		botapi.InlineButtonData("👍", "vote:up"),
		botapi.InlineButtonData("👎", "vote:down"),
	},
	[]botapi.InlineKeyboardButton{
		botapi.InlineButtonURL("source", "https://github.com/gotd/td"),
	},
)
bot.SendMessage(ctx, chat, "Vote:", botapi.WithReplyMarkup(kb))
```

Reply (custom) keyboards use `ReplyKeyboardMarkup` with `Button`,
`ButtonContact`, `ButtonLocation`; remove one with
`&botapi.ReplyKeyboardRemove{RemoveKeyboard: true}`.

## Sending media

A file to send is an `InputFile`: `FileID` (already on Telegram), `FileURL`
(Telegram fetches it), or a local upload (`FileFromPath`, `FileFromBytes`,
`FileFromReader`).

```go
bot.SendPhoto(ctx, chat, botapi.FileURL("https://.../cat.jpg"), "caption")
bot.SendDocument(ctx, chat, botapi.FileFromPath("/tmp/report.pdf"), "")
bot.SendVideo(ctx, chat, botapi.FileID(fileID), "")
```

Typed sends: `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`, `SendVoice`,
`SendAnimation`, `SendVideoNote`, `SendSticker`. Albums:
`SendMediaGroup(ctx, chat, []InputMedia{...})` (uploaded items).

## Other sends

`SendLocation`, `SendVenue`, `SendContact`, `SendPoll`, `SendDice`,
`SendChatAction`:

```go
bot.SendChatAction(ctx, chat, botapi.ChatActionTyping)
bot.SendPoll(ctx, chat, "Question?", []string{"A", "B", "C"})
bot.SendDice(ctx, chat, botapi.DiceDie)
```

## Receiving updates

Register handlers with the `On*` methods. A `Handler` is
`func(*Context) error`; the `Context` carries the `*Bot`, the `Update`, and is
itself a `context.Context`.

```go
bot.OnMessage(func(c *botapi.Context) error {
	return c.Reply("you said: " + c.Message().Text)
})

bot.OnEditedMessage(handler)
bot.OnChannelPost(handler)
bot.OnCallbackQuery(handler)
bot.OnInlineQuery(handler)
```

`Context` helpers: `Message()`, `Sender()`, `Chat()`, `Send`, `Reply`,
`AnswerCallback`, `AnswerInline`.

> Updates for the bot's own outgoing messages are filtered out (the HTTP Bot API
> never delivers them), so reply handlers won't answer themselves.

## Predicates

Every `On*` method accepts trailing `Predicate`s (`func(*Update) bool`); the
handler runs only when all match. First match wins across handlers.

```go
bot.OnMessage(handler, botapi.HasText(), botapi.Not(botapi.HasPrefix("/")))
```

Built-ins: `Command`, `HasPrefix`, `HasText`, `TextEquals`, `Regex`,
`ChatTypeIs`, `CallbackData`, `CallbackPrefix`, and the combinators `Not`/`Or`.
Write your own — it's just a function:

```go
func hasPhoto(u *botapi.Update) bool {
	m := u.EffectiveMessage()
	return m != nil && len(m.Photo) > 0
}
```

## Middleware

A `Middleware` is `func(Handler) Handler`. Register global middleware with
`Use`; it wraps every handler:

```go
bot.Use(botapi.Recover(), botapi.Timeout(30*time.Second), botapi.Logging())
```

Built-ins: `Recover` (turns panics into errors), `Timeout`, `Logging`.

## Groups

`Group` scopes shared predicates and middleware to a subset of handlers:

```go
admin := bot.Group(botapi.ChatTypeIs(botapi.ChatTypeSupergroup))
admin.Use(requireAdmin)
admin.OnCommand("ban", "Ban a user", banHandler)
```

## Commands

`OnCommand(name, description, handler, predicates...)` registers a command
handler. On start, the bot publishes all registered commands to Telegram via
`SetMyCommands`, so the client command menu stays in sync. Opt out with
`Options.DisableCommandRegistration`. You can still call
`SetMyCommands`/`GetMyCommands`/`DeleteMyCommands` directly with scopes
(`BotCommandScopeChat`, …).

## Callback & inline queries

Answer a callback query (acknowledge a button tap), optionally with a toast or
alert:

```go
bot.OnCallbackQuery(func(c *botapi.Context) error {
	if err := c.AnswerCallback(botapi.WithCallbackText("Thanks!")); err != nil {
		return err
	}
	m := c.Update.CallbackQuery.Message
	_, err := c.Bot.EditMessageText(c, botapi.ID(m.Chat.ID), m.MessageID, "done")
	return err
}, botapi.CallbackPrefix("vote:"))
```

Answer an inline query with results (enable inline mode in @BotFather first):

```go
bot.OnInlineQuery(func(c *botapi.Context) error {
	return c.AnswerInline([]botapi.InlineQueryResult{
		&botapi.InlineQueryResultArticle{
			ID:                  "1",
			Title:               "Echo",
			InputMessageContent: &botapi.InputTextMessageContent{MessageText: c.Update.InlineQuery.Query},
		},
	})
})
```

`InlineQueryResult` and `InputMessageContent` are sealed unions covering
articles, cached/URL media, and contact/location/venue results.

## Editing, forwarding, deleting

```go
bot.EditMessageText(ctx, chat, messageID, "new text")
bot.EditMessageCaption(ctx, chat, messageID, "new caption")
bot.EditMessageReplyMarkup(ctx, chat, messageID, markup)
bot.ForwardMessage(ctx, toChat, fromChat, messageID)
bot.CopyMessage(ctx, toChat, fromChat, messageID)
bot.DeleteMessage(ctx, chat, messageID)
bot.DeleteMessages(ctx, chat, []int{id1, id2})
```

Live locations: `EditMessageLiveLocation`, `StopMessageLiveLocation`.

## Files

There is no HTTP file server in the MTProto model. `GetFile` decodes a `file_id`
locally (no network) and derives `file_unique_id`; download with `DownloadFile`
or `DownloadFileToPath`, which follow DC migration:

```go
f, err := bot.GetFile(ctx, fileID)
n, err := bot.DownloadFile(ctx, fileID, w) // streams into an io.Writer
```

Incoming media populates the typed fields on `Message` (`Photo`, `Document`,
`Video`, `Sticker`, …), each carrying a usable `file_id`.

## Chat management

Members (supergroups/channels): `BanChatMember`, `UnbanChatMember`,
`RestrictChatMember`, `PromoteChatMember`, `GetChatMember`,
`GetChatAdministrators`, `GetChatMemberCount`. Admin: `PinChatMessage`,
`UnpinChatMessage`, `UnpinAllChatMessages`, `SetChatTitle`,
`SetChatDescription`, `SetChatPermissions`, `SetChatPhoto`, `DeleteChatPhoto`,
`LeaveChat`. Invite links: `ExportChatInviteLink`, `CreateChatInviteLink`,
`EditChatInviteLink`, `RevokeChatInviteLink`. Stickers: `UploadStickerFile`,
`CreateNewStickerSet`, `AddStickerToSet`, `DeleteStickerFromSet`,
`SetStickerPositionInSet`.

## Errors

Methods return errors shaped like the HTTP Bot API: an `*Error` with `Code` and
`Description`. Branch on it with `errors.As` or the helpers:

```go
if _, err := bot.SendMessage(ctx, chat, text); err != nil {
	if wait, ok := botapi.AsFloodWait(err); ok {
		time.Sleep(wait)
	} else if newID, ok := botapi.AsChatMigrated(err); ok {
		_ = newID // retry against newID (group upgraded to supergroup)
	} else if botapi.Code(err) == 403 {
		// blocked, or the bot is not a member of the chat
	}
}
```

Context cancellation passes through unchanged, so `errors.Is(err,
context.Canceled)` works.

## Resilience: flood-wait & rate limiting

Opt in via `Options`:

```go
botapi.Options{
	AppID: appID, AppHash: appHash,
	FloodWait:         true, // retry FLOOD_WAIT-limited requests transparently
	RequestsPerSecond: 25,   // proactive global token-bucket limit
}
```

`FloodWait` waits out limits instead of returning 429; `RequestsPerSecond`
(+ `RequestBurst`) caps outgoing MTProto requests.

## Persistence

By default everything is in memory (nothing survives a restart). Provide a
`Storage` to persist the session, peer access hashes and update state:

```go
db, _ := bbolt.Open("bot.bbolt", 0o666, nil)
opts := botapi.Options{AppID: appID, AppHash: appHash, Storage: storage.NewBBoltStorage(db)}
```

## Running many bots

`pool.Pool` lazily starts and multiplexes bots by token over one process — the
multi-bot front end (e.g. for a service serving many bots):

```go
p, _ := pool.New(pool.Options{AppID: appID, AppHash: appHash, StateDir: "state", IdleTimeout: time.Hour})
go p.RunGC(ctx)

err := p.Do(ctx, token, func(b *botapi.Bot) error {
	_, err := b.SendMessage(ctx, botapi.ID(chatID), "hi")
	return err
})
```

`Do` starts and authorizes the bot on first use (concurrent callers share one
startup), with per-token storage; `RunGC` reaps idle bots.

## The escape hatch

Anything the Bot API surface does not cover is one call away: `bot.Raw()`
returns the underlying `*tg.Client` for direct MTProto, and `bot.Dispatcher()`
exposes the raw update dispatcher.
