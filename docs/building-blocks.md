# gotd/td building blocks

This document maps the high-level `github.com/gotd/td` primitives the Bot API
library is built on. Unlike HTTP Bot API clients, this library does **not** talk
HTTP to `api.telegram.org`. It implements the Bot API *surface* directly over
MTProto using `gotd/td`. Everything below is a building block we wrap or
translate.

All signatures verified against the local checkout at `/src/gotd/td`
(`github.com/gotd/td`, module path confirmed). File:line references are to that
tree.

---

## 1. `telegram` — the MTProto client

`telegram.Client` is the connection. It is created unstarted and only connects
inside `Run`.

```go
// telegram/client.go:169
func NewClient(appID int, appHash string, opt Options) *Client

// connect.go — blocks; callback runs while connected, reconnects on transient errors
func (c *Client) Run(ctx context.Context, f func(ctx context.Context) error) error

// invoke.go:22 — the raw generated API, wrapped with middlewares
func (c *Client) API() *tg.Client

// the auth helper (bot login lives here)
func (c *Client) Auth() *auth.Client

// the current bot/user
func (c *Client) Self(ctx context.Context) (*tg.User, error)
```

### `telegram.Options` (the wiring points we care about)

`telegram/options.go:28`. The fields the library sets:

| Field | Purpose |
| --- | --- |
| `UpdateHandler UpdateHandler` | Where raw `tg.UpdatesClass` are delivered. We chain `peers.Manager.UpdateHook` → `updates.Manager` → our dispatcher here. |
| `SessionStorage SessionStorage` | Persisted MTProto session (auth key). Backed by bbolt per bot. |
| `Middlewares []Middleware` | Invoker middlewares — flood-wait retry, rate limiting, tracing, the `updhook` that routes updates. |
| `Logger log.Logger` | `go.uber.org/zap` based logging. |
| `Device DeviceConfig` | Device info reported at session init. |
| `NoUpdates bool` | Disable update stream (not used for bots that receive updates). |
| `OnSelfError` / `OnSelfSuccess` | Hooks around self-fetch. |

### Bot authentication

```go
// telegram/auth/bot.go:12 — imports bot authorization with a BotFather token
func (c *auth.Client) Bot(ctx context.Context, token string) (*tg.AuthAuthorization, error)

// telegram/auth/status.go — {Authorized bool, User *tg.User}
func (c *auth.Client) Status(ctx context.Context) (*auth.Status, error)

// telegram/bot.go:22 — convenience: reads BOT_TOKEN, logs in, runs callback
func BotFromEnvironment(ctx context.Context, opts Options,
    setup func(*Client) error, cb func(ctx context.Context, client *Client) error) error
```

Canonical bot startup:

```go
client.Run(ctx, func(ctx context.Context) error {
    st, err := client.Auth().Status(ctx)
    if err != nil { return err }
    if !st.Authorized {
        if _, err := client.Auth().Bot(ctx, token); err != nil { return err }
    }
    // peers init, gaps.Run, then serve.
    return nil
})
```

---

## 2. `telegram/message` — the Sender (outgoing)

`message.Sender` is the fluent builder for everything we send. It owns an
uploader and a peer resolver.

```go
// telegram/message/sender.go:25
func NewSender(raw *tg.Client) *Sender

// telegram/message/peer.go — pick the destination, returns *RequestBuilder
func (s *Sender) To(p tg.InputPeerClass) *RequestBuilder      // :55  already-resolved peer
func (s *Sender) Self() *RequestBuilder                        // :63  saved messages
func (s *Sender) Resolve(from string, ...) *RequestBuilder     // :80  @username, t.me/…, phone
func (s *Sender) ResolveDomain(domain string, ...) *RequestBuilder // :105
```

`RequestBuilder` embeds `Builder`. Builder configuration is chainable and
returns `tg.UpdatesClass` on send:

```go
b := sender.To(peer).
    NoWebpage().            // disable link preview
    Silent().               // no notification
    NoForwards().           // protect content
    Reply(replyToMsgID).    // reply
    Markup(replyMarkup)     // tg.ReplyMarkupClass

upd, err := b.StyledText(ctx, html.String(nil, text))  // entities/HTML
upd, err := b.Text(ctx, "plain")                        // plain text
upd, err := b.Media(ctx, message.Contact(...))          // any media
```

Key send methods (all return `tg.UpdatesClass`):

| Method | Use |
| --- | --- |
| `Text(ctx, msg)` / `Textf` | plain text |
| `StyledText(ctx, ...StyledTextOption)` | formatted text/entities |
| `Photo` / `Video` / `Audio` / `Document` / `Voice` / `GIF` | typed media from a `FileLocation` |
| `Media(ctx, MediaOption)` | generic media (contact, dice, poll, venue…) |
| `Upload(UploadOption) *UploadBuilder` | upload then send local files |

Upload sources (`message.From*`) produce `UploadOption`:
`FromPath`, `FromBytes`, `FromReader`, `FromFile`, `FromURL`, `FromFS`.

Typing/chat actions:

```go
// telegram/message/typing.go
rb.TypingAction().Typing(ctx)        // "typing"
rb.TypingAction().UploadDocument(ctx, progress)
rb.TypingAction().ChooseSticker(ctx)
```

`FileLocation` interface (`message/file.go`) is satisfied by `tg.Photo`,
`tg.Document`, etc. — the bridge from a decoded `file_id` back into a send.

---

## 3. `telegram/peers` — peer resolution & access-hash storage

This is the single most important building block, because the Bot API speaks in
bare `int64` chat IDs while MTProto needs an `InputPeer` **with an access hash**.

```go
// telegram/peers/options.go:31
func (o Options) Build(api *tg.Client) *Manager   // Options{Storage, Cache, Logger}

// telegram/peers/id.go:14 — resolve a Bot-API-style 64-bit ID to a Peer
func (m *Manager) ResolveTDLibID(ctx, peerID constant.TDLibPeerID) (Peer, error)

// telegram/peers/resolve.go:160 — resolve @username
func (m *Manager) ResolveDomain(ctx, domain string) (Peer, error)

// typed getters that cache access hashes
func (m *Manager) GetUser(ctx, *tg.InputUser) (User, error)
func (m *Manager) GetChannel(ctx, tg.InputChannelClass) (Channel, error)

func (m *Manager) Self(ctx) (User, error)
func (m *Manager) Init(ctx) error
```

`Peer` exposes `InputPeer() tg.InputPeerClass`, which is what `Sender.To` wants.

### The TDLib ID convention

Bot API chat IDs are `constant.TDLibPeerID` values: users are positive, chats
are offset-negative, channels/supergroups are `-100…`-prefixed. `ResolveTDLibID`
decodes that into the right peer kind and pulls the access hash from storage.
**Bots can resolve any peer they've "seen"** (received an update from, or that
appeared in `Entities`) because the access hash was persisted at that moment.

### Storage / Cache interfaces

```go
// telegram/peers/storage.go
type Storage interface {
    Save(ctx, key Key, value Value) error            // Key{Prefix,ID} -> {AccessHash}
    Find(ctx, key Key) (Value, found bool, err error)
    SavePhone / FindPhone / GetContactsHash / SaveContactsHash
}
type Cache interface {  // SaveUsers/FindUser/SaveChannels/FindChannel/…
}
```

Built-ins: `peers.NewInmemoryStorage()`, `peers.NoopCache`. We back both with
bbolt for persistence (see `botstorage`).

### Update integration hook

```go
// telegram/peers/integration.go:41 — wrap the next handler so every update's
// users/chats are harvested into storage before our code sees it
func (m *Manager) UpdateHook(next telegram.UpdateHandler) telegram.UpdateHandler
```

This is what keeps access hashes fresh without explicit calls.

---

## 4. `telegram/updates` — gap-aware update manager

MTProto updates can have gaps; the manager recovers them via
`updates.getDifference`.

```go
// telegram/updates/manager.go:40
func New(cfg Config) *Manager

// :55 — feed raw updates in (called via the UpdateHandler chain)
func (m *Manager) Handle(ctx, u tg.UpdatesClass) error

// :84 — start gap recovery loop for a given account
func (m *Manager) Run(ctx, api API, userID int64, opt AuthOptions) error
```

`updates.Config` (`config.go:22`) fields we set: `Handler` (terminal update
sink — our dispatcher), `Storage StateStorage`, `AccessHasher`,
`UserAccessHasher`, `OnChannelTooLong`, `Logger`.

`AuthOptions` (`manager.go:74`):

```go
type AuthOptions struct {
    IsBot   bool                      // true for us
    Forget  bool                      // ignore local state, resync from server
    OnStart func(ctx context.Context) // fired once recovery is live
}
```

Wiring order (from the existing pool, the pattern we keep):

```
telegram.Options.UpdateHandler =
    peers.Manager.UpdateHook(          // harvest access hashes
        updates.Manager (Handle)       // gap recovery
    )
// updates.Manager.Config.Handler = our tg.UpdateDispatcher (terminal)
gaps.Run(ctx, api, me.ID, updates.AuthOptions{IsBot: true, OnStart: …})
```

---

## 5. `tg.UpdateDispatcher` — typed update fan-out

The terminal handler. Register typed callbacks; it routes by update type and
hands you pre-resolved `Entities` (users/chats/channels in the update) so you
rarely need an extra RPC.

```go
// tg/tl_handlers_gen.go:41
func NewUpdateDispatcher() UpdateDispatcher

// :117
func (u UpdateDispatcher) OnNewMessage(NewMessageHandler)
// :497
func (u UpdateDispatcher) OnBotCallbackQuery(BotCallbackQueryHandler)
// plus OnNewChannelMessage, OnEditMessage, OnBotInlineQuery, OnDeleteMessages, …
```

```go
type NewMessageHandler func(ctx, e Entities, u *UpdateNewMessage) error

type Entities struct {  // tl_handlers_gen.go:47
    Short bool
    Users    map[int64]*User
    Chats    map[int64]*Chat
    Channels map[int64]*Channel
}
```

`UpdateDispatcher` implements `telegram.UpdateHandler`, so it slots in as the
`updates.Config.Handler`. This is the raw layer our Bot-API update layer maps
into Bot API `Update` objects.

---

## 6. `telegram/uploader` & `telegram/downloader` — files

```go
// telegram/uploader/uploader.go:16
func NewUploader(rpc Client) *Uploader
func (u *Uploader) FromPath(ctx, path) (tg.InputFileClass, error)
func (u *Uploader) FromReader(ctx, name, io.Reader) (tg.InputFileClass, error)
func (u *Uploader) FromBytes(ctx, name, []byte) (tg.InputFileClass, error)
func (u *Uploader) FromURL(ctx, url) (tg.InputFileClass, error)
func (u *Uploader) WithThreads(n).WithPartSize(n).WithProgress(p)
```

```go
// telegram/downloader/downloader.go:10
func NewDownloader() *Downloader            // also client.Downloader()
func (d *Downloader) Download(rpc, tg.InputFileLocationClass) *Builder
b.Stream(ctx, io.Writer) / b.Parallel(ctx, io.WriterAt) / b.ToPath(ctx, path)
```

The Sender already owns an uploader; we use the standalone uploader/downloader
for `getFile`/`uploadStickerFile`-style operations.

---

## 7. `fileid` — Bot API `file_id` ⇆ MTProto locations

The Bot API exposes opaque `file_id` strings; MTProto uses
`{id, access_hash, file_reference, dc}`. `fileid` is the codec.

```go
// fileid/decode.go:18
func DecodeFileID(s string) (FileID, error)
// fileid/encode.go:10
func EncodeFileID(id FileID) (string, error)

// fileid/from.go — build a FileID from MTProto objects
func FromDocument(doc *tg.Document) FileID       // :9
func FromPhoto(photo *tg.Photo, thumbType rune) FileID  // :39
// + FromChatPhoto, FromUserPhoto, …

type FileID struct {
    Type            Type   // Photo, Document, Video, ProfilePhoto, …
    DC              int
    ID, AccessHash  int64
    FileReference   []byte
    URL             string
    PhotoSizeSource PhotoSizeSource
}
```

Outgoing: decode the user's `file_id` → build a `tg.InputDocument`/`InputPhoto`
→ send without re-upload. Incoming: `FromDocument`/`FromPhoto` → `EncodeFileID`
→ hand the user a `file_id`. (Note: `file_unique_id` is a separate, currently
stubbed, derivation — see roadmap.)

---

## 8. `tgerr` — RPC error matching

MTProto errors carry a type and optional numeric argument
(`FLOOD_WAIT_3` → type `FLOOD_WAIT`, arg `3`). This is the source we translate
into Bot API `{error_code, description}`.

```go
// tgerr/error.go:14
type Error struct { Code int; Message, Type string; Argument int }
func (e *Error) IsType(t string) bool
func (e *Error) IsOneOf(...string) bool

// tgerr/error.go:139
func Is(err error, tt ...string) bool
// tgerr/flood_wait.go:24
func AsFloodWait(err error) (d time.Duration, ok bool)
```

Common types we map: `PEER_ID_INVALID`, `CHAT_NOT_FOUND`, `CHANNEL_PRIVATE`,
`USER_IS_BLOCKED`, `INPUT_USER_DEACTIVATED`, `MESSAGE_NOT_MODIFIED`,
`QUERY_ID_INVALID`, `FLOOD_WAIT_*`.

---

## 9. Session storage (`botstorage`)

A single bbolt-backed type implements three gotd interfaces at once:

- `session.Storage` — MTProto auth key (telegram.SessionStorage)
- `peers.Storage` + `peers.Cache` — access hashes & peer cache
- `updates.StateStorage` — pts/qts/seq gap state

One `*.bbolt` file per bot ID. This is reused largely as-is.

---

## 10. How the blocks compose (one bot)

```
                      ┌─────────────────────────────────────────┐
   BotFather token →  │ telegram.Client (MTProto, one per bot)   │
                      │   Auth().Bot(token)                      │
                      └───────────────┬──────────────────────────┘
                                      │ API() *tg.Client
        ┌─────────────────────────────┼───────────────────────────────┐
        ▼ outgoing                     ▼ peers                          ▼ incoming
  message.Sender            peers.Manager (Storage/Cache)      UpdateHandler chain:
  uploader/downloader       ResolveTDLibID / ResolveDomain      peers.UpdateHook
  fileid codec              InputPeer + access hashes            → updates.Manager (gaps)
                                                                 → tg.UpdateDispatcher
                                                                    (OnNewMessage, …)
        └──────────────── translation layer (ours) ───────────────────┘
                                      │
                        Bot API types & methods (hand-written)
```

The **translation layer** — MTProto `tg.*` ⇆ Bot API types — is the code we
own. The blocks above do the protocol heavy lifting; we never re-implement
MTProto, peer resolution, gap recovery, or file transfer.

---

## 11. What this buys us over HTTP Bot API clients

| HTTP Bot API client | this library (MTProto) |
| --- | --- |
| Subject to Bot API server rate limits | Direct MTProto; flood-wait only |
| `file_id` opaque, must round-trip server | `fileid` codec is local; can construct locations |
| Update delivery via long-poll/webhook HTTP | Native gap-aware update stream |
| No access to raw MTProto | `client.API()` escape hatch always available |
| Reflection-based request building | Typed builders, zero reflection |

See `architecture.md` for how we expose these as a clean, type-safe Bot API.
