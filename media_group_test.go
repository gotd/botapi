package botapi

import (
	"context"
	"errors"
	"testing"
)

func TestSendMediaGroupCountValidation(t *testing.T) {
	b := newTestBot(t)
	_, err := b.SendMediaGroup(context.Background(), ID(1), []InputMedia{
		&InputMediaPhoto{Type: InputMediaPhotoType, Media: FileFromBytes("a", []byte("x"))},
	})
	var apiErr *Error
	if !errors.As(err, &apiErr) || apiErr.Code != 400 {
		t.Fatalf("single-item group should be a 400, got %v", err)
	}
}

func TestMediaGroupRejectsNonUpload(t *testing.T) {
	b := newTestBot(t)
	_, err := b.inputMediaToMulti(context.Background(), &InputMediaPhoto{
		Type:  InputMediaPhotoType,
		Media: FileURL("https://example.com/a.jpg"),
	})
	if !errors.Is(err, error(errNonUploadInAlbum)) {
		t.Fatalf("URL item in album should be rejected, got %v", err)
	}
}
