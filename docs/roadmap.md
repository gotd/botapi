# Roadmap

> Rebuild `github.com/gotd/botapi` from a codegen-first OpenAPI/ogen project into
> a hand-written, MTProto-backed Bot API **library** that beats `telego`.
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
- ◐ Constructors / fluent setters (`ID`, `Username`, keyboard builders) — done;
  the `telegoutil` equivalent. More fluent setters land alongside the methods
  in Phase 3.
- ☑ Typed error hierarchy (`Error`, `AsFloodWait`, `ErrNotImplemented`)
- ☑ Exhaustiveness lint config for unions (`gochecksumtype` + `exhaustive`;
  golangci config migrated to v2)

## Phase 3 — Outgoing methods (translation)

Re-point the kept translation logic from `oas.*` to our types, then fill stubs.

- ☐ `SendMessage` (text + entities/HTML) — re-point `send.go`
- ☐ Media sends: `SendPhoto`, `SendDocument`, `SendVideo`, `SendAudio`,
  `SendVoice`, `SendAnimation`, `SendVideoNote`, `SendMediaGroup`,
  `SendSticker` (currently stubbed — uploader + `fileid`)
- ☐ `SendContact`, `SendDice`, `SendVenue`, `SendLocation`, `SendPoll`
- ☐ `SendChatAction` (already mapped)
- ☐ Keyboards both directions (`markup.go`)
- ☐ Edits: `EditMessageText/Caption/Media/ReplyMarkup`
- ☐ `ForwardMessage(s)`, `CopyMessage(s)`, `DeleteMessage(s)`
- ☐ Peer/chat-id resolution hardening (`ResolveTDLibID`, access-hash misses)

## Phase 4 — Incoming: updates & handler framework

- ☐ `tg.Update*` → `botapi.Update` mapping (message, edited, callback query,
  inline query, chat member, etc.) via the dispatcher
- ☐ `convert_message.go` re-point + finish reply-to / forward resolution
- ☐ `handler` package: `Group`, `Use`/middleware, predicates, routing, context
- ☐ Built-in middleware: `Recover`, `Timeout`, rate-limit, logging
- ☐ Built-in predicates: command, prefix, chat type, media, regex
- ☐ `Bot.On*` convenience registration
- ☐ Decide & implement any `getUpdates`/webhook compatibility shim

## Phase 5 — Files, queries, chat management

- ☐ `GetFile` + download; `file_id`/`file_unique_id` (resolve the stub)
- ☐ `UploadStickerFile`, sticker set methods
- ☐ `AnswerCallbackQuery`, `AnswerInlineQuery`, `AnswerPreCheckoutQuery`,
  `AnswerShippingQuery`
- ☐ Chat members: ban/unban/restrict/promote, `GetChatMember(s)`,
  `GetChatAdministrators`, `GetChatMemberCount`
- ☐ Chat admin: pin/unpin, photo, title, permissions, invite links (partly done)
- ☐ Commands: `Set/Get/DeleteMyCommands` (done in seed — re-point)
- ☐ Live location: `EditMessageLiveLocation`, `StopMessageLiveLocation`

## Phase 6 — Errors, rate limiting, resilience

- ☐ Complete `tgerr` → Bot API code mapping (`errors_map.go`)
- ☐ Flood-wait retry middleware (`tgerr.AsFloodWait`) + proactive limiter
- ☐ Context-cancellation semantics (return wrapped `ctx.Err()`)
- ☐ Reconnect/migration behavior surfaced sanely to callers

## Phase 7 — Multi-bot, server, polish

- ☐ `pool.Pool` re-pointed at public `Bot`; GC, keepalive
- ☐ `cmd/botapi` HTTP server as an optional example (local Bot-API server)
- ☐ Examples: echo bot, media bot, inline bot, handler/middleware
- ☐ Allocation tests on hot paths; benchmarks vs telego
- ☐ Conformance test against kept `botdoc` extractor (API-version drift guard)
- ☐ Docs: package docs, migration-from-telego guide, README

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

1. **`file_unique_id`** — derive properly now (Phase 5) or keep stubbed initially?
2. **HTTP Bot-API compatibility** — ship `getUpdates`/webhook shims for drop-in
   telego migration, or MTProto-native handlers only?
3. **Unions** — sealed interfaces (recommended) vs. type-param generics?
4. **Module surface** — single root package, or split `handler`/`pool` into
   sub-packages from day one?
5. **`appID`/`appHash`** — bundled default vs. caller-provided (bots still need
   an app identity for MTProto).
