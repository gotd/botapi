package botapi

import (
	"context"
	"errors"
	"testing"
)

func TestAnswerCallbackQueryInvalidID(t *testing.T) {
	b := &Bot{}

	err := b.AnswerCallbackQuery(context.Background(), "not-a-number")
	if err == nil {
		t.Fatal("expected error for non-numeric callback query id")
	}

	var apiErr *Error

	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("want 400 Error, got %v", err)
	}
}
