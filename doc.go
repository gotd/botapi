// Package botapi is a Telegram Bot API library implemented over MTProto using
// github.com/gotd/td.
//
// Unlike HTTP Bot API clients, botapi does not talk to api.telegram.org. It
// exposes the Bot API surface (types, methods, updates) but speaks MTProto
// directly, which avoids Bot API server rate limits and gives access to the
// raw client when needed (Bot.Raw).
//
// This package is under active reconstruction; see docs/roadmap.md.
package botapi
