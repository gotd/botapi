package botapi

import (
	"context"
	"errors"
	"testing"
)

func TestResolvePeerEmptyUsername(t *testing.T) {
	b, err := New("123:abc", Options{AppID: 1, AppHash: "hash"})
	if err != nil {
		t.Fatal(err)
	}

	_, err = b.resolvePeer(context.Background(), Username("@"))
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("empty username should be a 400, got %v", err)
	}
}
