package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

// noopHandler is a stub telegram.UpdateHandler for middleware-assembly tests.
type noopHandler struct{}

func (noopHandler) Handle(context.Context, tg.UpdatesClass) error { return nil }

func TestBuildMiddlewares(t *testing.T) {
	cases := []struct {
		name string
		opt  Options
		want int // expected chain length (the update hook is always present)
	}{
		{"Default", Options{}, 1},
		{"FloodWait", Options{FloodWait: true}, 2},
		{"RateLimit", Options{RequestsPerSecond: 30}, 2},
		{"Both", Options{FloodWait: true, MaxFloodWaitRetries: 3, RequestsPerSecond: 30, RequestBurst: 5}, 3},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := buildMiddlewares(c.opt, noopHandler{})
			if len(got) != c.want {
				t.Fatalf("chain length: got %d, want %d", len(got), c.want)
			}
		})
	}
}
