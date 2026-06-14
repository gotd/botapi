# botapi [![Go Reference](https://img.shields.io/badge/go-pkg-00ADD8)](https://pkg.go.dev/github.com/gotd/botapi#section-documentation) [![codecov](https://img.shields.io/codecov/c/github/gotd/botapi?label=cover)](https://codecov.io/gh/gotd/botapi) [![experimental](https://img.shields.io/badge/-experimental-blueviolet)](https://gotd.org/docs/projects/status#experimental)

A Telegram Bot API library for Go, implemented directly over **MTProto** using
[`gotd/td`](https://github.com/gotd/td) — not over HTTP to `api.telegram.org`.

It exposes the familiar [Bot API](https://core.telegram.org/bots/api) surface
(types, methods, updates) but speaks MTProto on a persistent connection. That
sidesteps the Bot API server's rate limits, removes the `getUpdates`/webhook
round trip, and keeps the raw `gotd/td` client one method call away
(`Bot.Raw()`) for anything not yet covered.

> **Status: under active reconstruction.** The repo is being rebuilt from a
> codegen-first OpenAPI/ogen project into a hand-written library. The `Bot`
> client skeleton compiles and connects today; the typed Bot API surface
> (methods, types, handler framework) is being filled in. See
> [`docs/roadmap.md`](./docs/roadmap.md) for what's done and what's next.

## Why MTProto instead of HTTP

| | HTTP Bot API client | botapi |
| --- | --- | --- |
| Transport | HTTPS to `api.telegram.org` | MTProto via `gotd/td` |
| Updates | `getUpdates` long-poll / webhook | persistent connection, no polling |
| Rate limits | Bot API server limits | MTProto limits only |
| Escape hatch | none | raw `*tg.Client` via `Bot.Raw()` |

## Design goals

In priority order (see [`docs/architecture.md`](./docs/architecture.md)):

1. **Zero-reflection performance** — fully typed request/response building, no
   `reflect` in the hot path; allocation-tested like `gotd/td`.
2. **Type-safe unions & enums** — `ChatID`, `InputFile`, `ChatMember`,
   `ReplyMarkup`, parse modes, etc. as sealed interfaces and typed constants,
   not stringly-typed structs.
3. **First-class context & structured errors** — context-first API; typed
   errors (flood-wait, retry-after, network vs API vs not-implemented);
   proactive rate limiting.
4. **A great handler framework** — composable middleware, router and predicates
   over a native MTProto update stream.

## Usage

```go
package main

import (
	"context"

	"github.com/gotd/botapi"
)

func main() {
	// Bots still need an MTProto app identity (https://my.telegram.org).
	// This is NOT the bot token.
	bot, err := botapi.New("<bot-token>", botapi.Options{
		AppID:   123456,
		AppHash: "<app-hash>",
	})
	if err != nil {
		panic(err)
	}

	bot.OnCommand("start", "Start the bot", func(c *botapi.Context) error {
		_, err := c.Reply("Hello!")
		return err
	})

	// Connects, authorizes as a bot, and serves updates until ctx is cancelled.
	if err := bot.Run(context.Background()); err != nil {
		panic(err)
	}
}
```

See the [**guide**](./docs/guide.md) for the full surface (sending, media,
keyboards, handlers, predicates, middleware, commands, files, chat management,
errors, pooling) and [`examples/`](./examples) for runnable bots
(`echo`, `buttons`, `inline`, `media`, `rich`, `background`, `advanced`,
`business`).

`Options.Storage` is optional — leave it nil to keep session, peers and update
state in memory (nothing survives a restart). `storage.Open("bot.bbolt")`
persists all of it to a single bbolt file; close it on shutdown. Every example
under [`examples/`](./examples) persists its session this way by default, so
they reconnect without re-authorizing.

## Package layout

- `botapi` (root) — the public library: the `Bot` client, options, and the
  hand-written Bot API surface as it lands.
- `pool` — runs and multiplexes many bots by token over one process.
- `storage` — bbolt-backed session/peer/update-state storage.
- `internal/botdoc` — fetches and extracts the published Bot API docs; kept as a
  reference oracle and a conformance check against API-version drift.
- `cmd/botdoc` — CLI to fetch and inspect the published Bot API docs.

## Acknowledgements

- [Bot API reference](https://core.telegram.org/bots/api) — the spec.
- [`gotd/td`](https://github.com/gotd/td) — the MTProto engine.
- [reference Bot API server](https://github.com/tdlib/telegram-bot-api).
