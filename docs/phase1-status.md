# Phase 1 status & handoff

Working branch: **`phase1-rebuild`** (off `main` @ `bd3da21`). All changes below
are **staged/unstaged but NOT committed**. Build is now **green** (`go mod tidy`,
`go build`, `go vet`, `go test` all pass) and `gotd/td` is bumped to the latest
release — see "Done".

Goal of Phase 1 (from `roadmap.md`): strip codegen, preserve the MTProto
translation logic as non-compiled seed, relocate clean reusable code, and stand
up a compiling root-package `Bot` skeleton.

## Done

- **Seeded (NOT compiled)** the oas-dependent translation layer. Go ignores
  `_`-prefixed dirs, so these are preserved as reference for re-pointing in
  Phase 3+ but excluded from the build:
  - `internal/botapi` → `_seed/botapi`
  - `internal/pool` → `_seed/pool`  (authoritative known-good v0.117.0 wiring is `_seed/pool/pool.go`)
  - `cmd/botapi` → `_seed/cmd-botapi`
  - `internal/guard_test.go` → `_seed/guard_test.go`
- **Relocated reusable, clean code:**
  - `internal/botstorage` → `storage/` (public); package renamed `botstorage` → `storage`; doc moved to `storage/doc.go`. `storage.BBoltStorage` implements `peers.Storage`, `peers.Cache`, `updates.StateStorage`, `session.Storage`.
  - `botdoc` → `internal/botdoc` (kept as reference oracle / future conformance test).
- **Deleted codegen:** `internal/oas`, `_oas/`, `cmd/gotd-bot-oas`, `botdoc/oas.go` (OAS emission), `tools.go`, `generate.go`.
- **New root package `botapi`** (the public library):
  - `doc.go` — rewritten package doc (MTProto-backed Bot API).
  - `options.go` — `Options{AppID,AppHash,Logger,Device,Storage,OnStart}` + combined `Storage` interface (session + peers + cache + updates state). Nil Storage ⇒ in-memory.
  - `bot.go` — `Bot` skeleton: `New`, `Run(ctx)`, `Raw`, `Dispatcher`, `Sender`, `Peers`, `Self`. Wiring mirrors `_seed/pool/pool.go` `createClient` (peers→gaps→UpdateHook→telegram.NewClient, with the `*rawPlaceholder = *client.API()` cycle-break; bot auth in `Run` via `Auth().Bot`, then `gaps.Run` with `IsBot:true`). The terminal update handler is a `tg.UpdateDispatcher` (the seed pool used a no-op there).
- **Makefile** — removed `generate` target, added `lint`.
- **Build green + gotd bumped.** Ran `go mod tidy` → `go build ./...` →
  `go vet ./...` → `go test ./...`, all pass. `tidy` dropped `go-chi/chi` and
  moved `ogen` to indirect (still pulled transitively by `gotd/td`); `goquery`
  stays (used by `internal/botdoc`).
  - **`gotd/td` v0.117.0 → v0.156.0.** The one breaking change that touched our
    code: logging moved off `*zap.Logger` onto the `github.com/gotd/log.Logger`
    port. `Options.Logger` stays `*zap.Logger` (caller-friendly); `bot.go`
    bridges it once via `logzap.New(opt.Logger)` and names sub-loggers with the
    package func `log.Named(lg, "peers"|"gaps"|"client")` (no longer a method).
    New deps: `gotd/log v0.1.0`, `gotd/log/logzap v0.1.1`.

## Next steps (resume here)

1. **README.md** — rewrite (still describes the old OpenAPI/ogen project).
2. Optional but nice: a minimal `examples/echo` and a smoke/compile test.
3. **Commit** Phase 1 (only when the maintainer asks). Suggested message:
   `refactor: strip codegen, seed translation layer, add MTProto-backed Bot skeleton`.

### Notes confirmed during the build pass

- `Options.Storage` (combined interface) flows nil straight through to
  `peers.Options{Storage,Cache}`, `updates.Config.Storage` and
  `telegram.Options.SessionStorage`; gotd defaults nil→in-memory, so no
  explicit branch needed. Builds and vets clean.
- `updates.Config.Handler` accepts the `tg.UpdateDispatcher` value; the same
  value is stored in `Bot.disp` and shares the handler map.

## Notes / decisions still open (from roadmap)

- `file_unique_id` derivation (stubbed in seed).
- Whether to ship `getUpdates`/webhook compatibility shims.
- Unions: sealed interfaces (recommended) vs generics.
- `appID`/`appHash`: currently **required** in `Options` (no bundled default).
- `storage/` is public now; `tdbot`/`pool` re-point happens in Phase 3/7 by
  copying from `_seed/` and swapping `oas.*` types for the hand-written ones.
