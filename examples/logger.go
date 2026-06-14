// Package examples holds helpers shared by the runnable example bots.
package examples

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger builds the logger shared by the example bots: the standard zap
// production configuration (structured JSONL on stderr) with the level lowered
// to Debug, so the library's debug output — MTProto RPC traces and the business
// peer diagnostics — is visible while debugging.
//
// The output is JSONL in the shape github.com/go-faster/pl expects; pipe it
// through pl for readable, colorized logs (see the examples README).
func NewLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	cfg.Level = zap.NewAtomicLevelAt(zapcore.DebugLevel)

	return cfg.Build()
}
