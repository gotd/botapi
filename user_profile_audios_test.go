package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetUserProfileAudios(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.UsersGetSavedMusicRequestTypeID, &tg.UsersSavedMusic{
		Count: 3,
		Documents: []tg.DocumentClass{
			&tg.Document{
				ID: 1, AccessHash: 2, FileReference: []byte{0}, MimeType: "audio/mpeg", Size: 1000, DCID: 2,
				Attributes: []tg.DocumentAttributeClass{
					&tg.DocumentAttributeAudio{Duration: 180, Title: "Song", Performer: "Artist"},
				},
			},
		},
	})

	got, err := newMockBot(inv).GetUserProfileAudios(context.Background(), 99, WithProfileAudiosLimit(50), WithProfileAudiosOffset(5))
	if err != nil {
		t.Fatalf("GetUserProfileAudios: %v", err)
	}

	if got.TotalCount != 3 || len(got.Audios) != 1 {
		t.Fatalf("got = %#v", got)
	}

	if got.Audios[0].Duration != 180 || got.Audios[0].Title != "Song" || got.Audios[0].Performer != "Artist" {
		t.Fatalf("audio = %#v", got.Audios[0])
	}

	var req tg.UsersGetSavedMusicRequest

	inv.decode(t, tg.UsersGetSavedMusicRequestTypeID, &req)

	if req.Limit != 50 || req.Offset != 5 {
		t.Fatalf("req = %#v", req)
	}
}
