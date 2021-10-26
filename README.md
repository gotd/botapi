# botapi

The telegram-bot-api, but in go. WIP.
* [API Reference](https://core.telegram.org/bots/api)
* [Reference implementation](https://github.com/tdlib/telegram-bot-api)
* [Generated OpenAPI v3 Schema](./_oas/openapi.json)

## Features
* Parsing of API documentation with defaults, format, enum and constraints inference
* OpenAPI v3 specification generation
* Server and Client generation based on OpenAPI v3 specification

## Roadmap
- [x] Parse definition
- [x] Generate OpenAPI v3 Specification
- [x] Generate client and server from OpenAPi v3 using [ogen](https://github.com/ogen-go/ogen)
- [ ] Infer enums
- [ ] Infer defaults
- [ ] Use rich text for documentation
- [ ] More links to documentation
- [ ] Support Emoji
