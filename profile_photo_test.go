package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetMyProfilePhoto(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PhotosUploadProfilePhotoRequestTypeID, &tg.PhotosPhoto{Photo: &tg.PhotoEmpty{}})

	photo := InputProfilePhotoStatic{Photo: &InputFileUpload{Name: "p.jpg", Bytes: []byte("img")}}
	if err := newMockBot(inv).SetMyProfilePhoto(context.Background(), photo); err != nil {
		t.Fatalf("SetMyProfilePhoto: %v", err)
	}

	var req tg.PhotosUploadProfilePhotoRequest

	inv.decode(t, tg.PhotosUploadProfilePhotoRequestTypeID, &req)

	if _, ok := req.GetFile(); !ok {
		t.Fatal("static photo should set File")
	}
}

func TestRemoveMyProfilePhoto(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PhotosUpdateProfilePhotoRequestTypeID, &tg.PhotosPhoto{Photo: &tg.PhotoEmpty{}})

	if err := newMockBot(inv).RemoveMyProfilePhoto(context.Background()); err != nil {
		t.Fatalf("RemoveMyProfilePhoto: %v", err)
	}

	var req tg.PhotosUpdateProfilePhotoRequest

	inv.decode(t, tg.PhotosUpdateProfilePhotoRequestTypeID, &req)

	if _, ok := req.ID.(*tg.InputPhotoEmpty); !ok {
		t.Fatalf("id = %#v, want empty", req.ID)
	}
}
