package botapi

import (
	"context"
	"errors"
	"testing"
)

func newTestBot(t *testing.T) *Bot {
	t.Helper()

	b, err := New("123:abc", Options{AppID: 1, AppHash: "hash"})
	if err != nil {
		t.Fatal(err)
	}

	return b
}

func TestPhotoMediaURL(t *testing.T) {
	b := newTestBot(t)
	// URL media needs no network to construct.
	media, err := b.photoMedia(context.Background(), FileURL("https://example.com/p.jpg"), nil)
	if err != nil || media == nil {
		t.Fatalf("url photo media: got %v, %v", media, err)
	}
}

func TestMediaEmptyUploadRejected(t *testing.T) {
	b := newTestBot(t)
	_, err := b.documentMedia(context.Background(), &InputFileUpload{}, nil)

	var apiErr *Error

	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("empty upload should be a 400, got %v", err)
	}
}

func TestPhotoMediaBadFileID(t *testing.T) {
	b := newTestBot(t)
	_, err := b.photoMedia(context.Background(), FileID("not-a-valid-file-id"), nil)

	var apiErr *Error

	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("bad file_id should be a 400, got %v", err)
	}
}
