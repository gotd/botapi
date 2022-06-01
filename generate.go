package botapi

//go:generate go run ./cmd/gotd-bot-oas _oas/openapi.json
//go:generate go run github.com/ogen-go/ogen/cmd/ogen --target internal/oas --package oas --clean _oas/openapi.json
