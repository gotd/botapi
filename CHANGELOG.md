# Changelog

All notable changes to this project are documented here. The format is based on
[Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and the project aims to
follow [Semantic Versioning](https://semver.org/spec/v2.0.0.html). Commit
messages follow [Conventional Commits](https://www.conventionalcommits.org/).

## [Unreleased]

### Changed

- **Logging port** — the library now logs through `github.com/gotd/log` instead
  of zap directly. `Options.Logger` and `Bot.Logger()` (and `pool.Options.Logger`)
  are now `log.Logger`. **Breaking:** wrap a `*zap.Logger` with
  `github.com/gotd/log/logzap.New` (or a `*slog.Logger` with `logslog.New`). The
  library no longer depends on zap.

### Added

- **Rich messages** (Bot API 10.1) — `SendRichMessage`/`SendRichHTML`/
  `SendRichMarkdown` send structured page-block content built with
  `github.com/gotd/td/telegram/message/rich`. The `examples/rich` bot showcases
  the page-block and rich-text constructors valid in a bot-sent message.
  (Instant-View-page-only blocks — Title, Subtitle, Header, Subheader, Kicker,
  AuthorDate, Cover, RelatedArticles — and the auto-link inline styles are
  rejected by the server with `RICH_VALIDATE_CTOR_NOT_ALLOWED`.)
- **Background sends** — `Bot.Background()` / `Context.Background()` expose a
  run-lifetime context for proactive sends to any chat from timers, queues or
  goroutines, instead of the per-update handler context. Plus `Bot.Logger()`.
- **Serializable peer references** — `Bot.PeerRef` captures a chat's id and
  access hash into a JSON-serializable `PeerRef`; `Peer(ref)` addresses it
  directly (no stored peer data, no re-resolution), so background/scheduled
  sends survive a restart.

## [0.1.0] - 2026-06-14

`botapi` was rebuilt from a codegen-first OpenAPI/ogen project into a
hand-written, MTProto-backed Bot API **library** built on
[`gotd/td`](https://github.com/gotd/td). It exposes the familiar Bot API surface
but speaks MTProto directly — no `api.telegram.org`, no `getUpdates`/webhooks.

### Added

- **Client** — `Bot` with `New`/`Run`/`Raw`, a persistent MTProto connection and
  gap-aware update stream. Optional bbolt `Storage` for session/peers/state.
- **Types** — hand-written Bot API types and sealed-interface unions (`ChatID`,
  `InputFile`, `ReplyMarkup`, `InputMedia`, `ChatMember`, `MessageOrigin`,
  `InlineQueryResult`, `InputMessageContent`, `BotCommandScope`,
  `PassportElementError`, …) with compile-time exhaustiveness checks.
- **Sending** — `SendMessage` (HTML/MarkdownV2/Markdown), photo/document/video/
  audio/voice/animation/video-note/sticker, media groups, location/venue/
  contact/poll/dice, chat actions; edits (text/caption/markup/media, live
  location), forward/copy/delete, `StopPoll`.
- **Receiving** — a handler framework (`On*`, predicates, middleware, `Group`,
  `Context`) over the native update stream; messages, edits, channel posts,
  callback/inline queries, shipping/pre-checkout queries; incoming media, polls,
  contacts and forward origins mapped to the typed `Message`.
- **Files** — `GetFile`, `DownloadFile`/`DownloadFileToPath`, local
  `file_unique_id` derivation.
- **Queries** — `AnswerCallbackQuery`, `AnswerInlineQuery`,
  `AnswerShippingQuery`, `AnswerPreCheckoutQuery`.
- **Chat management** — members (ban/unban/restrict/promote, get member(s)/
  admins/count, custom title), admin (pin/unpin, title/description/photo/
  permissions, sticker set, leave), invite links, `GetChat`,
  `GetUserProfilePhotos`.
- **Commands** — `Set`/`Get`/`DeleteMyCommands` with scopes; `OnCommand`
  auto-publishes the command menu.
- **Stickers** — `UploadStickerFile`, sticker-set create/add/delete/reorder,
  `GetStickerSet`, `SetStickerSetThumb`.
- **Payments & games** — `SendInvoice`, `SetPassportDataErrors`, `SendGame`,
  `SetGameScore`, `GetGameHighScores`.
- **Resilience** — Bot-API-shaped errors with a comprehensive `tgerr` mapping
  and `AsFloodWait`/`AsChatMigrated`/`Code` helpers; opt-in flood-wait retry and
  a proactive rate limiter.
- **Multi-bot** — `pool.Pool` runs and multiplexes many bots by token with idle
  GC.
- **Tooling & docs** — `cmd/botdoc` (fetch/inspect the published API), a
  method-drift conformance test, hot-path benchmarks, package docs, a usage
  guide and runnable examples (`echo`, `buttons`, `inline`, `media`,
  `advanced`).

[Unreleased]: https://github.com/gotd/botapi/compare/v0.1.0...HEAD
[0.1.0]: https://github.com/gotd/botapi/releases/tag/v0.1.0
