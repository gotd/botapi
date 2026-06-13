package botapi

import (
	"context"

	"github.com/gotd/botapi/internal/oas"
)

// AnswerCallbackQuery implements oas.Handler.
func (b *BotAPI) AnswerCallbackQuery(ctx context.Context, req *oas.AnswerCallbackQuery) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}

// AnswerInlineQuery implements oas.Handler.
func (b *BotAPI) AnswerInlineQuery(ctx context.Context, req *oas.AnswerInlineQuery) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}

// AnswerPreCheckoutQuery implements oas.Handler.
func (b *BotAPI) AnswerPreCheckoutQuery(ctx context.Context, req *oas.AnswerPreCheckoutQuery) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}

// AnswerShippingQuery implements oas.Handler.
func (b *BotAPI) AnswerShippingQuery(ctx context.Context, req *oas.AnswerShippingQuery) (*oas.Result, error) {
	return nil, &NotImplementedError{}
}
