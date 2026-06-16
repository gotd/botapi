# Examples

Runnable bots showing the library in practice. Each is its own `main` package:

```bash
APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/echo
```

Bots need an MTProto app identity (`APP_ID`/`APP_HASH`, from
<https://my.telegram.org>) plus a BotFather token (`BOT_TOKEN`).

## demo — the full tour

[`demo`](./demo) is the comprehensive one: a single bot that exercises every
major feature, split across focused files so each subsystem reads on its own —
`main.go` (wiring & lifecycle), `commands.go` (formatting, content, editing),
`keyboards.go` (inline + reply keyboards and callbacks), `media.go` (sending and
receiving media), `inline.go` (inline mode), `admin.go` (a group-scoped command
set: chat info, reactions, pinning, the raw escape hatch), `text.go` (free-text
predicates and edits) and shared `middleware.go`/`helpers.go`. Enable inline mode
in @BotFather to try the inline queries.

```bash
APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/demo
```

The other directories are smaller, single-feature bots.

## Logging

The bots log structured **JSONL** via zap. The shared
[`examples.NewLogger`](./logger.go) uses `zap.NewProductionConfig` but lowers the
level to **Debug**, so MTProto RPC traces and the business peer diagnostics show
up — handy when verifying behavior against the live API.

Raw JSON is hard to read in a terminal. Pipe it through
[`github.com/go-faster/pl`](https://github.com/go-faster/pl), which tails and
pretty-prints exactly this `zap.NewProductionConfig` JSONL.

### Install pl

```bash
go install github.com/go-faster/pl/cmd/pl@latest
```

### Use it

zap writes to **stderr**, so redirect it into `pl` with `2>&1`:

```bash
APP_ID=12345 APP_HASH=abcdef BOT_TOKEN=123:abc go run ./examples/business 2>&1 | pl
```

Useful flags:

- `--level info` — hide debug lines (keep info and above)
- `--no-time` — drop timestamps
- `--no-color` — disable colors (or set `NO_COLOR`)
- `-f service.log` — follow a file like `tail -f`

Non-JSON lines pass through untouched, so mixed output (e.g. a panic stack
trace) stays readable.

To capture a session and read it back later:

```bash
go run ./examples/business 2>session.log
pl session.log          # read once
pl -f session.log       # follow while the bot runs
```
