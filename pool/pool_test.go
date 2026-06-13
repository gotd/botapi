package pool

import (
	"context"
	"testing"
	"time"

	"github.com/gotd/botapi"
)

func TestParseToken(t *testing.T) {
	tok, err := ParseToken("123456:ABC-DEF")
	if err != nil {
		t.Fatal(err)
	}
	if tok.ID != 123456 || tok.Secret != "ABC-DEF" {
		t.Fatalf("unexpected: %#v", tok)
	}
	if tok.String() != "123456:ABC-DEF" {
		t.Fatalf("round-trip: %q", tok.String())
	}

	for _, bad := range []string{"", "nocolon", "abc:def", "123456"} {
		if _, err := ParseToken(bad); err == nil {
			t.Fatalf("ParseToken(%q) should fail", bad)
		}
	}
}

func TestNewRequiresAppIdentity(t *testing.T) {
	if _, err := New(Options{}); err == nil {
		t.Fatal("New should require AppID/AppHash")
	}
	if _, err := New(Options{AppID: 1, AppHash: "x"}); err != nil {
		t.Fatalf("New with identity: %v", err)
	}
}

func TestDoRejectsInvalidToken(t *testing.T) {
	p, err := New(Options{AppID: 1, AppHash: "x"})
	if err != nil {
		t.Fatal(err)
	}
	err = p.Do(context.Background(), "not-a-token", func(*botapi.Bot) error { return nil })
	if err == nil {
		t.Fatal("Do should reject an invalid token before starting anything")
	}
}

func TestRunGCNoTimeoutStopsOnContext(t *testing.T) {
	p, err := New(Options{AppID: 1, AppHash: "x"})
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { p.RunGC(ctx); close(done) }()
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("RunGC did not return after context cancellation")
	}
}
