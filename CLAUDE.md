# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`github.com/gotd/botapi` is a Telegram Bot API library that implements the familiar Bot API
surface (types, methods, updates) **directly over MTProto** via `github.com/gotd/td` — it does
**not** talk HTTP to `api.telegram.org`. A bot connects once and serves updates over a
persistent connection. The raw `*tg.Client` is always reachable through `Bot.Raw()` for anything
the typed surface does not cover.

The project was rebuilt from an old codegen-first (OpenAPI/ogen) design into a hand-written
library; that reconstruction is complete and ongoing work is incremental Bot API feature parity.
**Note:** `docs/architecture.md` and `docs/roadmap.md` describe the original target/plan —
`architecture.md` in particular references a *target* layout (`tdbot/`, `handler/` subpackages)
that does **not** exist; the real code is a flat root package. Trust the code over those docs.

## Commands

- Build: `go build ./...`
- Test (race, as CI runs it): `make test` → `go test --timeout 5m -race ./...`
- Single test: `go test -run TestName -count=1 ./`
- Coverage: `make coverage` (writes `profile.out`, prints per-func coverage)
- Lint: `make lint` → `golangci-lint run ./...`
- Tidy modules: `make tidy`
- Inspect the published Bot API (reference oracle): `go run ./cmd/botdoc` (see its `-methods`,
  `-name`, `-json`, `-file` flags)

`TestRunEndToEnd` (`run_e2e_test.go`) boots a real in-process **teled** server backed by a
throwaway PostgreSQL **container** (`teledtest`). It self-skips on hosts without Docker/container
support, so a passing local run without Docker does **not** mean it ran.

## Architecture

### Lifecycle (`bot.go`, `options.go`)

`New` builds an unconnected `Bot` from a BotFather token plus an MTProto app identity
(`AppID`/`AppHash` — these are *not* the token). It does no I/O. `Run` connects, authorizes as a
bot, initializes peers, and blocks in gap recovery serving updates until the context is canceled;
`OnStart` fires once after recovery is live.

`New` resolves a construction cycle: `peers.Manager` needs a `*tg.Client`, but the client's update
handler is built *from* the manager. It builds the manager against a placeholder `*tg.Client`, then
backfills it after `telegram.NewClient` (the same trick gotd examples use). Don't "simplify" this away.

The engine sits on these `gotd/td` primitives: `telegram.Client`, `peers.Manager` (access-hash
storage + harvesting via `UpdateHook`), `updates.Manager` (gap recovery), `message.Sender`, and
`tg.UpdateDispatcher`.

### File-group conventions (flat root package `botapi`)

- `types_*.go` — hand-written Bot API data types (Message, Chat, User, payments, media, queries).
- `unions.go`, `enums.go` — **sealed-interface unions** (`ChatID`, `InputFile`, `ReplyMarkup`,
  `ChatMember`, `InlineQueryResult`, …) and typed enums. Unions use an unexported marker method so
  illegal states are unrepresentable. Switches over them are linter-enforced exhaustive — see below.
- `codec*.go` — JSON via `github.com/go-faster/jx`: receivable entities implement `Encode(*jx.Encoder)`
  / `Decode(*jx.Decoder)`, exposed to `encoding/json` through generated `MarshalJSON`/`UnmarshalJSON`.
  The `json:"..."` struct tags are documentation of wire names, not what drives (de)serialization.
- `convert*.go` — **inbound** translation: `tg.*` → `botapi.*` (`convert.go`, `convert_media.go`,
  `convert_member.go`). Pure where possible.
- Outbound **methods hang off `*Bot`**, context-first, taking a `ChatID` target and functional
  `SendOption`s (`send.go`, `send_*.go`, `edit*.go`, `chat_*.go`, `commands.go`, `forward_delete.go`, …).
- `markup.go` / `markup_to_tg.go` — keyboards, both directions.
- `errors.go` / `errors_map.go` — methods return an `*Error{Code, Description}` shaped like the HTTP
  Bot API; `errors_map.go` maps `tgerr` RPC errors to those codes. Context cancellation passes
  through unchanged (`errors.Is(err, ctx.Err())` still works).

### Handler framework

`on.go` installs handlers on the raw `tg.UpdateDispatcher` and converts each raw update into a
`botapi.Update`, then `route`s it. `handler.go` holds the `router`, `Context` (embeds the request
context + `Bot` + `Update`), and the `Handler`/`Predicate`/`Middleware` types. Register with the
`On*` helpers, narrow with `predicates.go`, layer with `middleware.go` / `Group` / `Use`.
Update-conversion failures are logged and swallowed so one bad update never tears down the stream.

### Addressing peers / access hashes

Sending requires a peer's MTProto access hash. The bot harvests and persists hashes for peers it
sees, but addressing a peer it hasn't seen (e.g. after a restart) needs stored peer data.
`resolve.go` turns a `ChatID` (numeric via TDLib id convention, or `@username` by domain) into a
peer. `PeerRef` (`peerref.go`) is a self-contained, JSON-serializable `{kind, id, access_hash}` you
can persist and later send to via `Peer(ref)` without re-resolution. A `PeerRef` is send-only — it
can't back the `peers.Peer` that chat-management methods need.

### Subpackages

- `pool/` — runs/multiplexes many bots by token in one process (lazy start, idle GC). The multi-bot
  front end.
- `storage/` — `BBoltStorage`, one bbolt file holding session + peer cache + update state. The
  composite `Storage` interface (in `options.go`) requires all of those; `Options.Storage` is
  optional and everything stays in memory when nil — but is effectively mandatory in production:
  without a persisted session every `Run` re-authorizes the bot and Telegram answers repeated
  logins with a growing `FLOOD_WAIT`.
- `internal/botdoc/` + `cmd/botdoc/` — fetch/extract the published Bot API docs. Kept as a reference
  oracle and as a conformance check (`conformance_test.go`) that catches *unacknowledged* drift
  between the hand-written surface and the published API.

## Testing patterns

- **Hermetic per-method tests** use `mockInvoker` (`harness_test.go`): an in-memory `tg.Invoker`
  that dispatches by request `TypeID`, records calls, and returns canned responses the real client
  decodes back. `newMockBot` wires a `Bot` to it; `userRef`/`channelRef`/`tdlibChannel` build
  addressable targets that skip resolution.
- **End-to-end** tests drive a real teled server through `teledtest`. `Options` has unexported test
  seams (`resolver`, `publicKeys`, `dcList`) used to point a `Bot` at the in-process server; these
  are settable only from within the `botapi` package's own tests.

## Conventions

- **Exhaustiveness is enforced**, not optional: `exhaustive` (typed enums) and `gochecksumtype`
  (sealed unions) run in CI. Both treat a `default:` case as covering the remaining variants. When
  you add a union variant or enum value, fix every switch or add a `default:`.
- Lint config (`.golangci.yml`) is strict (gosec, gocritic diagnostics, lll @140, govet shadow,
  unparam, …). Run `make lint` before finishing.
- Commit messages follow Conventional Commits (`commitlint` runs in CI): `feat:`, `fix:`, `test:`,
  `chore:`, `docs:`, etc.
- No `reflect` in the hot path and allocation-consciousness are explicit design goals — prefer the
  typed jx codecs and fully-typed request building over reflection-based shortcuts.

## Documentation

When updating the docs, also update `../docs` (gotd.dev).
