package pool

import (
	"context"
	"errors"
	"testing"
	"time"
)

// fakeManaged returns a managed entry that is already "ready", with a cancel
// hook recording whether it was killed.
func fakeManaged(killed *bool) *managed {
	m := &managed{ready: make(chan struct{})}
	m.cancel = func() { *killed = true }
	m.markReady(nil)
	return m
}

func TestManagedMarkReadyLatches(t *testing.T) {
	m := &managed{ready: make(chan struct{})}
	first := errors.New("first")
	m.markReady(first)
	m.markReady(errors.New("second")) // ignored
	select {
	case <-m.ready:
	default:
		t.Fatal("ready not closed")
	}
	if m.startErr != first {
		t.Fatalf("startErr = %v, want first", m.startErr)
	}
}

func TestManagedIdleBefore(t *testing.T) {
	m := &managed{ready: make(chan struct{})}
	// Never used: not idle.
	if m.idleBefore(time.Now()) {
		t.Fatal("unused bot should not be idle")
	}
	m.use()
	if m.idleBefore(time.Now().Add(-time.Hour)) {
		t.Fatal("recently used bot should not be idle before an old deadline")
	}
	if !m.idleBefore(time.Now().Add(time.Hour)) {
		t.Fatal("bot should be idle before a future deadline")
	}
}

func TestKillRemovesAndCancels(t *testing.T) {
	p, _ := New(Options{AppID: 1, AppHash: "x"})
	killed := false
	p.bots["123:abc"] = fakeManaged(&killed)

	p.Kill("123:abc")
	if !killed {
		t.Fatal("Kill should cancel the bot")
	}
	if _, ok := p.bots["123:abc"]; ok {
		t.Fatal("Kill should remove the bot")
	}
	// Killing an unknown token is a no-op.
	p.Kill("nope")
}

func TestCloseKillsAll(t *testing.T) {
	p, _ := New(Options{AppID: 1, AppHash: "x"})
	k1, k2 := false, false
	p.bots["1:a"] = fakeManaged(&k1)
	p.bots["2:b"] = fakeManaged(&k2)

	p.Close()
	if !k1 || !k2 {
		t.Fatalf("Close should kill all: %v %v", k1, k2)
	}
	if len(p.bots) != 0 {
		t.Fatalf("pool not emptied: %d", len(p.bots))
	}
}

func TestReapCollectsIdle(t *testing.T) {
	p, _ := New(Options{AppID: 1, AppHash: "x"})
	idleKilled, freshKilled := false, false

	idle := fakeManaged(&idleKilled)
	idle.lastUsed = time.Now().Add(-time.Hour)
	p.bots["idle"] = idle

	fresh := fakeManaged(&freshKilled)
	fresh.use()
	p.bots["fresh"] = fresh

	p.reap(time.Now().Add(-time.Minute))

	if !idleKilled {
		t.Fatal("idle bot should be reaped")
	}
	if freshKilled {
		t.Fatal("fresh bot should survive")
	}
	if _, ok := p.bots["idle"]; ok {
		t.Fatal("idle bot not removed")
	}
	if _, ok := p.bots["fresh"]; !ok {
		t.Fatal("fresh bot removed")
	}
}

func TestRunGCReapsThenStops(t *testing.T) {
	p, _ := New(Options{AppID: 1, AppHash: "x", IdleTimeout: 10 * time.Millisecond})
	killed := false
	m := fakeManaged(&killed)
	m.lastUsed = time.Now().Add(-time.Hour)
	p.bots["idle"] = m

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() { p.RunGC(ctx); close(done) }()

	// Wait for the idle bot to be reaped.
	deadline := time.After(2 * time.Second)
	for {
		p.mu.Lock()
		n := len(p.bots)
		p.mu.Unlock()
		if n == 0 {
			break
		}
		select {
		case <-deadline:
			t.Fatal("idle bot was not reaped")
		case <-time.After(5 * time.Millisecond):
		}
	}
	if !killed {
		t.Fatal("reaped bot should be killed")
	}

	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("RunGC did not stop")
	}
}

func TestAcquireDedupes(t *testing.T) {
	p, _ := New(Options{AppID: 1, AppHash: "x"})
	killed := false
	existing := fakeManaged(&killed)
	p.bots["123:abc"] = existing

	tok, err := ParseToken("123:abc")
	if err != nil {
		t.Fatal(err)
	}
	got, err := p.acquire(tok)
	if err != nil {
		t.Fatalf("acquire: %v", err)
	}
	if got != existing {
		t.Fatal("acquire should return the existing managed bot, not start a new one")
	}
}
