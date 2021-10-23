package botapi

//go:generate go run ./cmd/gotd-bot-oas _oas/openapi.json
//go:generate go run github.com/ogen-go/ogen/cmd/ogen --schema _oas/openapi.json --target internal/oas --package oas --clean
