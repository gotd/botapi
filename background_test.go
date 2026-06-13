package botapi

import (
	"context"
	"testing"
)

func TestBackgroundContext(t *testing.T) {
	b := newTestBot(t)

	// Before Run, Background is already canceled so background sends fail fast.
	if err := b.Background().Err(); err == nil {
		t.Fatal("Background before Run should be canceled")
	}

	// While running, it returns the live run context.
	ctx, cancel := context.WithCancel(context.Background())
	b.setRunCtx(ctx)
	if got := b.Background(); got != ctx {
		t.Fatalf("Background should return the run context, got %v", got)
	}
	if b.Background().Err() != nil {
		t.Fatal("live run context should not be canceled")
	}

	// Canceling the run context (bot stopping) cancels background work.
	cancel()
	if b.Background().Err() == nil {
		t.Fatal("canceled run context should propagate")
	}

	// After Run returns, runCtx is cleared back to a canceled context.
	b.setRunCtx(nil)
	if b.Background().Err() == nil {
		t.Fatal("Background after stop should be canceled")
	}
}

func TestContextBackgroundDelegates(t *testing.T) {
	b := newTestBot(t)
	c := &Context{Bot: b, Update: &Update{}}
	if c.Background() != b.Background() {
		t.Fatal("Context.Background should delegate to Bot.Background")
	}
}
