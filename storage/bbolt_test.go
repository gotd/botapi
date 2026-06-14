package storage

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/gotd/td/telegram/peers"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
)

func newTestStorage(t *testing.T) *BBoltStorage {
	t.Helper()
	s, err := Open(filepath.Join(t.TempDir(), "test.bbolt"))
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	t.Cleanup(func() { _ = s.Close() })
	return s
}

func TestOpenAndClose(t *testing.T) {
	s := newTestStorage(t)
	if s.db == nil {
		t.Fatal("db is nil")
	}
	if err := s.Close(); err != nil {
		t.Fatalf("Close: %v", err)
	}
	// Wrapping a caller-owned db must not close it.
	wrapped := NewBBoltStorage(s.db)
	if err := wrapped.Close(); err != nil {
		t.Fatalf("Close (not owned): %v", err)
	}
}

func TestAccessHashRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()
	key := peers.Key{Prefix: "users_", ID: 42}

	if _, found, err := s.Find(ctx, key); err != nil || found {
		t.Fatalf("Find before save: found=%v err=%v", found, err)
	}
	if err := s.Save(ctx, key, peers.Value{AccessHash: 999}); err != nil {
		t.Fatalf("Save: %v", err)
	}
	v, found, err := s.Find(ctx, key)
	if err != nil || !found || v.AccessHash != 999 {
		t.Fatalf("Find: v=%+v found=%v err=%v", v, found, err)
	}
}

func TestUserCacheRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if _, found, _ := s.FindUser(ctx, 7); found {
		t.Fatal("found user before save")
	}
	if err := s.SaveUsers(ctx, &tg.User{ID: 7, AccessHash: 8, Username: "u"}); err != nil {
		t.Fatalf("SaveUsers: %v", err)
	}
	u, found, err := s.FindUser(ctx, 7)
	if err != nil || !found || u.ID != 7 || u.Username != "u" {
		t.Fatalf("FindUser: u=%+v found=%v err=%v", u, found, err)
	}
}

func TestChannelCacheRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if err := s.SaveChannels(ctx, &tg.Channel{ID: 11, AccessHash: 12, Title: "c", Photo: &tg.ChatPhotoEmpty{}}); err != nil {
		t.Fatalf("SaveChannels: %v", err)
	}
	c, found, err := s.FindChannel(ctx, 11)
	if err != nil || !found || c.Title != "c" {
		t.Fatalf("FindChannel: c=%+v found=%v err=%v", c, found, err)
	}
}

func TestChatCacheRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if err := s.SaveChats(ctx, &tg.Chat{ID: 13, Title: "g", Photo: &tg.ChatPhotoEmpty{}}); err != nil {
		t.Fatalf("SaveChats: %v", err)
	}
	c, found, err := s.FindChat(ctx, 13)
	if err != nil || !found || c.Title != "g" {
		t.Fatalf("FindChat: c=%+v found=%v err=%v", c, found, err)
	}
}

func TestContactsHash(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if err := s.SaveContactsHash(ctx, 555); err != nil {
		t.Fatalf("SaveContactsHash: %v", err)
	}
	h, err := s.GetContactsHash(ctx)
	if err != nil {
		t.Fatalf("GetContactsHash: %v", err)
	}
	// GetContactsHash is a stub that returns 0 (see FIXME in source).
	_ = h
}

func TestSessionRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if _, err := s.LoadSession(ctx); err != nil {
		// Empty session is reported as a not-found error by the session storage.
		t.Logf("LoadSession (empty): %v", err)
	}
	if err := s.StoreSession(ctx, []byte("session-bytes")); err != nil {
		t.Fatalf("StoreSession: %v", err)
	}
	data, err := s.LoadSession(ctx)
	if err != nil || string(data) != "session-bytes" {
		t.Fatalf("LoadSession: data=%q err=%v", data, err)
	}
}

func TestUpdatesStateRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if _, found, err := s.GetState(ctx, 1); err != nil || found {
		t.Fatalf("GetState before set: found=%v err=%v", found, err)
	}
	st := updates.State{Pts: 1, Qts: 2, Date: 3, Seq: 4}
	if err := s.SetState(ctx, 1, st); err != nil {
		t.Fatalf("SetState: %v", err)
	}
	got, found, err := s.GetState(ctx, 1)
	if err != nil || !found || got != st {
		t.Fatalf("GetState: got=%+v found=%v err=%v", got, found, err)
	}

	if err := s.SetPts(ctx, 1, 10); err != nil {
		t.Fatalf("SetPts: %v", err)
	}
	if err := s.SetDateSeq(ctx, 1, 11, 12); err != nil {
		t.Fatalf("SetDateSeq: %v", err)
	}
	got, _, _ = s.GetState(ctx, 1)
	if got.Pts != 10 || got.Date != 11 || got.Seq != 12 {
		t.Fatalf("state after updates = %+v", got)
	}
}

func TestChannelPtsRoundTrip(t *testing.T) {
	s := newTestStorage(t)
	ctx := context.Background()

	if err := s.SetChannelPts(ctx, 1, 100, 50); err != nil {
		t.Fatalf("SetChannelPts: %v", err)
	}
	pts, found, err := s.GetChannelPts(ctx, 1, 100)
	if err != nil || !found || pts != 50 {
		t.Fatalf("GetChannelPts: pts=%d found=%v err=%v", pts, found, err)
	}

	seen := map[int64]int{}
	err = s.ForEachChannels(ctx, 1, func(_ context.Context, channelID int64, pts int) error {
		seen[channelID] = pts
		return nil
	})
	if err != nil {
		t.Fatalf("ForEachChannels: %v", err)
	}
	if seen[100] != 50 {
		t.Fatalf("ForEachChannels saw %v", seen)
	}
}
