# Roadmap

> Rebuild `github.com/gotd/botapi` from a codegen-first OpenAPI/ogen project into
> a hand-written, MTProto-backed Bot API **library** that beats existing HTTP
> Bot API clients on ergonomics and performance.
> Companion docs: `architecture.md`, `building-blocks.md`.

Legend: ☐ todo · ◐ in progress · ☑ done

## Phase 0 — Planning & docs (current)

- ☑ Map `gotd/td` building blocks (`building-blocks.md`)
- ☑ Inventory existing `botapi` reusable translation logic
- ☑ Target architecture (`architecture.md`)
- ☑ This roadmap
- ☐ Confirm scope decisions with maintainer (see "Decisions needed")

## Phase 1 — Demolition & skeleton

Strip codegen, keep the engine, stand up an empty-but-compiling library.
**In progress** on branch `phase1-rebuild` — see `phase1-status.md` for the
detailed handoff.

- ☑ Remove `internal/oas`, `botdoc` OAS emission (`oas.go`),
  `cmd/gotd-bot-oas`, `_oas/`, ogen tooling (`tools.go`/`generate.go`)
- ☑ Keep `botdoc` **fetch/extract** only; moved under `internal/botdoc`
- ☑ Preserve `internal/botapi` + `internal/pool` + `cmd/botapi` as non-compiled
  seed under `_seed/` (re-point to `tdbot`/`pool` in later phases);
  `internal/botstorage` → `storage` (public)
- ☑ New root `Bot` type with construction, `Run(ctx)`, `Raw() *tg.Client`,
  wiring the verified update chain (peers hook → gaps → dispatcher)
- ☑ Repo builds with the translation layer detached — **`go mod tidy`,
  `go build ./...`, `go vet ./...` and `go test ./...` all green**
- ☑ Bump `gotd/td` to latest (**v0.117.0 → v0.156.0**); `tidy` dropped
  `go-chi/chi`, moved `ogen` to indirect, added `gotd/log` + `gotd/log/logzap`
- ◐ Update `Makefile` (done: dropped `generate`, added `lint`); rewrite
  `README.md` (☑ rewritten for the MTProto-backed library)

## Phase 2 — Core type system

The hand-written Bot API surface. This is the bulk of the work. **Done** on
`main` (commits on top of Phase 1): see `enums.go`, `errors.go`, `markup.go`,
`unions.go`, `types_*.go`, `message_origin.go`, `chat_member.go`,
`input_media.go`, `update.go`.

- ☑ Primitive types: `User`, `Chat`, `Message`, `MessageEntity`, `PhotoSize`,
  `Document`, `Update`, `ResponseParameters` (+ media/query supporting types)
- ☑ Typed enums: `ParseMode`, `ChatType`, `ChatAction`, `MessageEntityType`, … +
  the union discriminators
- ☑ Sealed-interface unions: `ChatID`, `InputFile`, `ReplyMarkup`,
  `MessageOrigin`, `ChatMember`, `ReactionType`, `MenuButton`, `InputMedia`
- ◐ Constructors / fluent setters (`ID`, `Username`, keyboard builders) — done.
  More fluent setters land alongside the methods in Phase 3.
- ☑ Typed error hierarchy (`Error`, `AsFloodWait`, `ErrNotImplemented`)
- ☑ Exhaustiveness lint config for unions (`gochecksumtype` + `exhaustive`;
  golangci config migrated to v2)

## Phase 3 — Outgoing methods (translation)

Hand-written over the gotd sender on our types. **Done** on `main`: methods on
`*Bot` with shared functional `SendOption`s; pure translation
(`markup_to_tg.go`, `entities.go`, `errors_map.go`, `convert.go`) is unit
tested. Live-Telegram paths are compile-/lint-verified.

- ☑ `SendMessage` (text + HTML/MarkdownV2/Markdown) — `send.go`
- ☑ Media sends: `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`,
  `SendVoice`, `SendAnimation`, `SendVideoNote`, `SendSticker` (file_id via
  `fileid`, URL, or upload via the uploader); `SendMediaGroup` (uploaded albums)
- ☑ `SendContact`, `SendDice`, `SendVenue`, `SendLocation`, `SendPoll`
- ☑ `SendChatAction`
- ☑ Keyboards both directions (`markup_to_tg.go` out, `convert.go` in)
- ◐ Edits: `EditMessageText/Caption/ReplyMarkup` done; `EditMessageMedia` TODO
- ☑ `ForwardMessage`, `CopyMessage`, `DeleteMessage(s)`
- ◐ Peer/chat-id resolution (`resolve.go`: TDLib id + @username); access-hash
  miss hardening continues as real traffic exposes cases

Deferred within Phase 3: explicit-entity sends (parse modes cover formatting);
`SendMediaGroup` with file_id/URL items (only uploads compose through the
high-level album API); `EditMessageMedia`.

## Phase 4 — Incoming: updates & handler framework

**Done** on `main`. The framework lives in the root package (on `*Bot`),
consistent with the Phase 3 methods. `installHandlers` (called from New) binds
the raw `tg.UpdateDispatcher` to a concurrency-safe router. See
`examples/echo` for the end-to-end pipe.

- ☑ `tg.Update*` → `botapi.Update` mapping: new/edited messages, channel posts
  (broadcast → channel_post, supergroup → message), callback & inline queries
  (`updates_map.go`, `on.go`); senders resolved from harvested `Entities`
- ☑ `convert.go` re-point + reply-to + forward-origin resolution
  (user/hidden/chat/channel)
- ☑ Routing: `Group`, `Use`/middleware (global + group-scoped), predicates,
  first-match dispatch, `Context` (`handler.go`, `group.go`, `context.go`)
- ☑ Built-in middleware: `Recover`, `Timeout`, `Logging` (rate-limit → Phase 6)
- ☑ Built-in predicates: `Command`, `HasPrefix`, `HasText`, `TextEquals`,
  `Regex`, `ChatTypeIs`, `CallbackData`/`CallbackPrefix`, `Not`/`Or`
  (media predicate → Phase 5 alongside incoming media)
- ☑ `Bot.On*` convenience: `OnMessage`, `OnEditedMessage`, `OnChannelPost`,
  `OnCallbackQuery`, `OnInlineQuery`, `OnCommand`
- ☑ **Decision: no `getUpdates`/webhook shim.** Updates arrive on the
  persistent MTProto connection; the native handler framework + `Bot.Run` is
  the only model. Decided against an HTTP-poll/webhook compatibility surface
  (resolves Decisions-needed #2).

Deferred within Phase 4: `chat_member`/`my_chat_member`, `poll`/`poll_answer`
and `chosen_inline_result` update routing (types exist; dispatcher wiring lands
with the chat-management and query work in Phase 5).

## Phase 5 — Files, queries, chat management

**Done** on `main` (payment answers and a couple of sticker reads aside). The
methods are on `*Bot` with the same functional-option style as Phase 3. See
`file.go`, `answer.go`, `commands.go`, `chat_member_methods.go`,
`chat_admin.go`, `chat_photo.go`, `invite_links.go`, `live_location.go`,
`inline_query_result.go`, `input_message_content.go`, `sticker.go`.

- ☑ `GetFile` + download (`DownloadFile`/`DownloadFileToPath`);
  `file_unique_id` resolved — derived locally from the decoded `file_id` with
  the TDLib scheme (web/document exact; legacy photos via volume/local id, newer
  photo sources fall back to media id). Resolves Decisions-needed #1. No HTTP
  file server in the MTProto-native model, so `GetFile` is decode-only.
- ◐ `AnswerCallbackQuery` (+ `Context.AnswerCallback`) and `AnswerInlineQuery`
  (+ `Context.AnswerInline`) done — the `InlineQueryResult` union (article;
  photo/gif/mpeg4 gif by URL; cached photo/gif/sticker/document/video/voice/
  audio by file_id; contact/location/venue) and the `InputMessageContent` union
  (text/location/venue/contact). Payment answers
  `AnswerPreCheckoutQuery`/`AnswerShippingQuery` deferred (need payment-update
  plumbing, which has no incoming updates wired yet).
- ☑ Chat members: `Ban`/`Unban`/`Restrict`/`PromoteChatMember`,
  `GetChatMember`, `GetChatAdministrators`, `GetChatMemberCount`
  (supergroups/channels via `channels.*`); `ChatPermissions`/`ChatAdminRights`
  with MTProto rights mapping; participant → `ChatMember` converter.
- ☑ Chat admin: pin/unpin (`PinChatMessage`/`UnpinChatMessage`/
  `UnpinAllChatMessages`), `SetChatTitle`/`SetChatDescription`,
  `SetChatPermissions`, `LeaveChat`, `SetChatPhoto`/`DeleteChatPhoto`; invite
  links (`Export`/`Create`/`Edit`/`RevokeChatInviteLink`).
- ☑ Commands: `Set`/`Get`/`DeleteMyCommands` with the `BotCommandScope` union.
- ☑ Live location: `EditMessageLiveLocation`, `StopMessageLiveLocation`.
- ◐ Stickers: `UploadStickerFile`, `CreateNewStickerSet`, `AddStickerToSet`,
  `DeleteStickerFromSet`, `SetStickerPositionInSet` done (`InputSticker` +
  `StickerFormat`). `GetStickerSet`/`SetStickerSetThumb` deferred (need full
  `Sticker[]` conversion).

Deferred within Phase 5: payment answers
(`AnswerPreCheckoutQuery`/`AnswerShippingQuery`) until payment updates land;
`GetStickerSet`/`SetStickerSetThumb`.

## Phase 6 — Errors, rate limiting, resilience

**Done** on `main`. See `errors_map.go`, `errors.go`, `ratelimit.go`.

- ☑ `tgerr` → Bot API mapping completed: a ~50-entry table of verbatim official
  descriptions (`errors_map.go`, mirroring telegram-bot-api `Client.cpp`), plus
  the server's code-normalization (sub-400/404 → 400, SCREAMING_CASE 403 → 400)
  and prefix/casing fallback for unmapped errors. Helpers `AsFloodWait`,
  `AsChatMigrated`, `Code`.
- ☑ Flood-wait retry + proactive limiter: opt-in `Options.FloodWait`
  (+ `MaxFloodWaitRetries`) and `Options.RequestsPerSecond` (+ `RequestBurst`),
  wired as client invoker middlewares via `gotd/contrib`
  (`floodwait`/`ratelimit`). Off by default.
- ☑ Context-cancellation semantics: `asAPIError` passes `context.Canceled`/
  `DeadlineExceeded` through unchanged (even when RPC-wrapped) so callers can
  `errors.Is` on `ctx.Err()`.
- ◐ Reconnect handled transparently by the gotd client + gaps manager; group →
  supergroup migration is surfaced via `AsChatMigrated`. A connection-state
  callback is deferred to Phase 7 polish.

## Phase 7 — Multi-bot, server, polish

- ☐ `pool.Pool` re-pointed at public `Bot`; GC, keepalive
- ☐ `cmd/botapi` HTTP server as an optional example (local Bot-API server)
- ◐ Examples: `examples/echo` (handler + middleware), `examples/buttons`
  (inline keyboards + callback queries), `examples/inline` (inline mode). Media
  bot still to add.
- ☑ Allocation tests on hot paths (`bench_test.go`): entity/markup/user
  conversion and `file_unique_id`, with `-benchmem`.
- ☐ Conformance test against kept `botdoc` extractor (API-version drift guard)
- ◐ Docs: package docs (`doc.go`) and README done; reference/migration guide
  still to write.

## Phase 8 — Release

- ☐ CI (lint, race tests, conformance), codecov
- ☐ Semantic version, changelog, conventional commits
- ☐ Announce; migration guide

---

## Sequencing notes

- Phases 2–4 are the critical path; 3 unblocks the most user value.
- Re-pointing reused code (`tdbot`) only needs the Phase-2 types to exist, so
  build a **vertical slice first**: `User`/`Chat`/`Message`/`ChatID` →
  `SendMessage` → echo update flow. Prove the whole pipe before going wide.
- Keep each method's translation behind the seed logic where it already exists;
  prefer re-pointing over rewriting.

## Decisions needed (maintainer)

1. ~~**`file_unique_id`**~~ — **Resolved (Phase 5): derived locally** from the
   decoded `file_id` with the TDLib scheme (web/document exact; photos best
   effort, stable per file).
2. ~~**HTTP Bot-API compatibility**~~ — **Resolved (Phase 4): MTProto-native
   handlers only, no `getUpdates`/webhook shim.**
3. ~~**Unions**~~ — **Resolved (Phase 2): sealed interfaces**, guarded by the
   `gochecksumtype` exhaustiveness linter.
4. ~~**Module surface**~~ — **Resolved (Phases 2–4): single root package**;
   `pool` re-points in Phase 7.
5. **`appID`/`appHash`** — bundled default vs. caller-provided (currently
   caller-provided and required).
