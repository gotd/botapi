# Architecture

> Status: **planning**. This describes the target architecture for the rebuilt
> `github.com/gotd/botapi` — a hand-written, MTProto-backed Telegram Bot API
> library. See `roadmap.md` for sequencing and `building-blocks.md` for the
> `gotd/td` primitives we sit on.

## What we are building

A Go library that exposes the **Telegram Bot API** surface (types, methods,
updates) but implements it **directly over MTProto** via `gotd/td` — not over
HTTP to `api.telegram.org`. Existing HTTP Bot API clients set the ergonomic bar
to beat; the Bot API docs (<https://core.telegram.org/bots/api>) are the spec.

Design goals, in priority order (per project decision):

1. **Zero-reflection performance** — request/response building is fully typed,
   no `reflect` in the hot path; allocation-tested like `gotd/td`.
2. **Type-safe unions & enums** — `ChatID`, `InputFile`, `ChatMember`,
   `ReplyMarkup`, `MessageEntity`, parse modes, etc. are sealed interfaces or
   generics, not stringly-typed structs.
3. **First-class context & structured errors** — context-first API; typed
   errors (`FloodWait`, retry-after, network vs API vs not-implemented);
   proactive rate limiting.
4. **A great handler framework** — composable middleware/router/predicates over
   a native MTProto update stream.

## What changes from today's `botapi`

Today's repo is **codegen-first**: `botdoc` scrapes the docs → OpenAPI →
`ogen` generates `internal/oas` (client + server) → `internal/botapi` implements
the *server* `oas.Handler` on top of `gotd/td`. We are inverting this.

| Component | Fate |
| --- | --- |
| `internal/oas` (ogen output) | **delete** — replaced by hand-written types |
| `botdoc` OpenAPI generation (`oas.go`, OAS emit) | **delete** |
| `cmd/gotd-bot-oas`, `_oas/openapi.json`, ogen tooling | **delete** |
| `botdoc` doc **fetch/extract** (HTML → structured API) | **keep** — used as a reference oracle for hand-writing & for doc strings / a verification test that our hand-written surface matches the published API |
| `internal/botapi` translation logic (convert_message, markup, peers, errors_map, send*, fileid) | **keep & re-point** from `oas.*` types to our hand-written types |
| `internal/pool`, `internal/botstorage` | **keep** — client lifecycle, bbolt storage |
| `cmd/botapi` (HTTP server) | **drop** — the library is the product; no HTTP Bot-API server is planned |

Net: the project becomes a **library** whose public API is hand-written Bot API
types + a `Bot` client, with the MTProto translation as its engine.

## Package layout (target)

Module: `github.com/gotd/botapi` (unchanged path; the *shape* changes).

```
botapi/                      root package — the public library
  bot.go                     Bot client: construction, Run, raw API escape hatch
  options.go                 Options, functional opts
  types_*.go                 hand-written Bot API types (Message, Chat, User, …)
  enums.go                   typed enums (ParseMode, ChatType, ChatAction, …)
  unions.go                  sealed-interface/generic unions (ChatID, InputFile, …)
  methods_*.go               hand-written methods (SendMessage, GetChat, …)
  errors.go                  typed error hierarchy + RPC mapping
  updates.go                 Update type + the update source
  doc.go

  tdbot/        (internal)   MTProto engine: translation tg.* ⇆ botapi.*
     convert_message.go      tg.Message  -> botapi.Message
     convert_media.go        media/entities/file_id
     markup.go               keyboards both directions
     peers.go                ChatID -> InputPeer (access hashes)
     send.go / send_media.go outgoing translation
     errors_map.go           tgerr -> botapi error codes
     engine.go               owns sender/peers/uploader/downloader/dispatcher

  handler/                   the dispatcher framework
     handler.go              Group, Use (middleware), predicates, routing
     predicates.go
     context.go

  pool/                      multi-bot pooling (from internal/pool)
  storage/                   bbolt session/peer/state storage (from botstorage)

  internal/botdoc/           kept doc fetch/extract (reference oracle + tests)
  cmd/botdoc/                optional: fetch & inspect the published API
  examples/
  docs/
```

> `tdbot` is internal so the MTProto coupling never leaks into the public type
> surface — a user holds `botapi.Message`, never a `tg.Message`. The raw
> `*tg.Client` is still reachable via an explicit escape hatch
> (`bot.Raw()`), mirroring `gotd/td`'s philosophy.

## Core types & the four goals

### Type-safe unions

A common HTTP-client approach models `ChatID` as a two-field struct
(`{ID int64; Username string}`) and `InputFile` similarly — illegal states are
representable. We use sealed
interfaces (a private method makes them unforgeable from outside the package),
with ergonomic constructors:

```go
type ChatID interface{ isChatID() }
type ChatIDInt int64
type ChatIDUsername string      // both implement isChatID()

func ID(id int64) ChatID          { return ChatIDInt(id) }
func Username(u string) ChatID    { return ChatIDUsername(u) }

type InputFile interface{ isInputFile() }
// InputFileID | InputFileURL | InputFileUpload(reader/path/bytes)
```

Discriminated incoming unions (`ChatMember`, `MessageOrigin`, `ReactionType`,
`MenuButton`, `InputMedia`, inline-query results) are sealed interfaces with one
concrete type per variant and a type switch — compile-time exhaustiveness via a
linter, no `interface{}` + reflection, no runtime "try each type" unmarshal.

### Typed enums

```go
type ParseMode string
const ( ParseModeHTML ParseMode = "HTML"; ParseModeMarkdownV2 ParseMode = "MarkdownV2" )

type ChatAction string  // ChatActionTyping, ChatActionUploadPhoto, …
type ChatType string     // private/group/supergroup/channel
```

These map to MTProto concepts internally; the public value is a typed constant,
not a bare string the caller can mistype.

### Zero-reflection request building

There is no JSON marshaling step at all on the wire — methods translate their
typed params straight into `gotd/td` `message.Builder`/`tg.*` calls. Where we do
serialize (e.g. callback data, `file_id`), it is explicit byte handling.
Hot-path methods get `testutil.ZeroAlloc`/`MaxAlloc` coverage as in `gotd/td`.

### First-class context & errors

Every method is `func (b *Bot) X(ctx, params) (Result, error)`. Errors form a
typed hierarchy translated from `tgerr`:

```go
type Error struct {            // implements error
    Code        int            // Bot-API-compatible (400/403/429/…)
    Description string
    Parameters  *ResponseParameters  // retry_after, migrate_to_chat_id
    raw         error          // wrapped tgerr.Error
}

func AsFloodWait(err error) (retryAfter time.Duration, ok bool)
var ErrNotImplemented = …
```

Proactive rate limiting and flood-wait retry live as **invoker middlewares**
(`gotd/td` `telegram.Middleware`), so they apply uniformly and are testable.

### Handler framework

A native update source built on `tg.UpdateDispatcher` → mapped to Bot API
`Update`, then a composable router:

```go
b.OnMessage(handler, FilterCommand("start"))
b.OnCallback(handler, FilterPrefix("vote:"))
grp := b.Group(FilterChatType(ChatTypePrivate))
grp.Use(Recover(), Timeout(30*time.Second))
```

Predicates are `func(ctx, Update) bool`; middleware is
`func(next Handler) Handler`. Context carries the `*Bot`, the update, and
per-update values — designed in from the start, with no per-request HTTP
`context.WithoutCancel` foot-gun because updates arrive over a persistent
MTProto stream rather than per-request HTTP.

## Update flow (no long-poll, no webhook)

Because we are on MTProto, there is **no `getUpdates` and no webhook**. Updates
arrive on the persistent connection. The chain (verified wiring) is:

```
tg updates → peers.Manager.UpdateHook (harvest access hashes)
           → updates.Manager.Handle    (gap recovery)
           → tg.UpdateDispatcher        (typed fan-out)
           → tdbot mapping              (tg.Update* -> botapi.Update)
           → handler router             (predicates, middleware, handlers)
```

We still **offer** a `SetWebhook`/`GetUpdates`-shaped compatibility surface only
where it makes sense for drop-in migration, but the primary model is "register
handlers, call `bot.Run(ctx)`".

## Multi-bot & single-bot

- **Single bot**: `botapi.New(appID, appHash, token, opts…)` → `bot.Run(ctx)`.
- **Many bots** (a server): `pool.Pool` keyed by token, lazy-constructs and
  GCs idle `Bot`s, each with its own `*.bbolt`. This is the existing pool,
  re-pointed at the public `Bot`.

## Testing strategy

- **Translation unit tests** with `tgmock` (mock `tg.Invoker`): feed known
  `tg.*` objects, assert `botapi.*` output, and vice-versa. No live Telegram.
- **`tgtest`** in-process server for end-to-end auth + send + receive.
- **Allocation tests** on hot paths.
- **Conformance test**: the kept `botdoc` extractor parses the live docs; a test
  asserts our hand-written method/type set matches the published surface (names,
  required fields), catching drift when Telegram ships a new Bot API version.
- Coverage filters generated-free since nothing is generated now.

## Open questions (tracked in roadmap)

- `file_unique_id` derivation (currently stubbed in the seed code).
- How much HTTP Bot-API compatibility (`getUpdates`/webhook shims) to ship.
- Whether `ChatID`/`InputFile` use sealed interfaces vs. type-param generics —
  leaning sealed interfaces for ergonomics + exhaustiveness linting.
