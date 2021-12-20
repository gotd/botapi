package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// GetGameHighScores implements oas.Handler.
func (b *BotAPI) GetGameHighScores(ctx context.Context, req oas.GetGameHighScores) (oas.ResultArrayOfGameHighScore, error) {
	return oas.ResultArrayOfGameHighScore{}, &NotImplementedError{}
}

// SetGameScore implements oas.Handler.
func (b *BotAPI) SetGameScore(ctx context.Context, req oas.SetGameScore) (oas.Result, error) {
	return oas.Result{}, &NotImplementedError{}
}
