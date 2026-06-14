package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetBusinessAccountProfilePhotoStatic(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PhotosPhoto{Photo: &tg.PhotoEmpty{}})

	photo := InputProfilePhotoStatic{Photo: &InputFileUpload{Name: "p.jpg", Bytes: []byte("img")}}
	if err := newMockBot(inv).SetBusinessAccountProfilePhoto(context.Background(), "bc1", photo, true); err != nil {
		t.Fatalf("SetBusinessAccountProfilePhoto: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PhotosUploadProfilePhotoRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	req, ok := wrapper.Query.(*tg.PhotosUploadProfilePhotoRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := req.GetFile(); !ok {
		t.Fatal("static photo should set File")
	}

	if !req.GetFallback() {
		t.Fatal("is_public should set Fallback")
	}
}

func TestSetBusinessAccountProfilePhotoAnimated(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PhotosPhoto{Photo: &tg.PhotoEmpty{}})

	photo := InputProfilePhotoAnimated{
		Animation:          &InputFileUpload{Name: "a.mp4", Bytes: []byte("vid")},
		MainFrameTimestamp: 1.5,
	}
	if err := newMockBot(inv).SetBusinessAccountProfilePhoto(context.Background(), "bc1", photo, false); err != nil {
		t.Fatalf("SetBusinessAccountProfilePhoto: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PhotosUploadProfilePhotoRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.PhotosUploadProfilePhotoRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := req.GetVideo(); !ok {
		t.Fatal("animated photo should set Video")
	}

	if ts, ok := req.GetVideoStartTs(); !ok || ts != 1.5 {
		t.Fatalf("video start ts = %v, ok=%v", ts, ok)
	}

	if req.GetFallback() {
		t.Fatal("Fallback should be unset when is_public is false")
	}
}

func TestSetBusinessAccountProfilePhotoRejectsFileID(t *testing.T) {
	inv := newMockInvoker()

	photo := InputProfilePhotoStatic{Photo: FileID("abc")}

	err := newMockBot(inv).SetBusinessAccountProfilePhoto(context.Background(), "bc1", photo, false)
	if err == nil {
		t.Fatal("expected error for non-upload profile photo")
	}

	if inv.count() != 0 {
		t.Fatal("should not make an RPC for an invalid file")
	}
}

func TestRemoveBusinessAccountProfilePhoto(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PhotosPhoto{Photo: &tg.PhotoEmpty{}})

	if err := newMockBot(inv).RemoveBusinessAccountProfilePhoto(context.Background(), "bc1", true); err != nil {
		t.Fatalf("RemoveBusinessAccountProfilePhoto: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PhotosUpdateProfilePhotoRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.PhotosUpdateProfilePhotoRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := req.ID.(*tg.InputPhotoEmpty); !ok {
		t.Fatalf("id = %#v, want empty", req.ID)
	}

	if !req.GetFallback() {
		t.Fatal("is_public should set Fallback")
	}
}
