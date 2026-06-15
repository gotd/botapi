package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/constant"
	"github.com/gotd/td/tg"
)

func storyUpdates() *tg.Updates {
	return &tg.Updates{
		Updates: []tg.UpdateClass{
			&tg.UpdateStory{
				Peer: &tg.PeerUser{UserID: 10},
				Story: &tg.StoryItem{
					ID:         42,
					Date:       1,
					ExpireDate: 2,
					Media:      &tg.MessageMediaEmpty{},
				},
			},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, FirstName: "Biz", Username: "biz"}},
	}
}

func TestPostStory(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, storyUpdates())

	story, err := newMockBot(inv).PostStory(context.Background(), "bc1",
		InputStoryContentPhoto{Photo: FileURL("https://example.com/p.jpg")}, 86400,
		WithStoryCaption("hi"), WithStoryPostToChatPage(), WithStoryProtectContent())
	if err != nil {
		t.Fatalf("PostStory: %v", err)
	}

	if story.ID != 42 {
		t.Fatalf("story id = %d, want 42", story.ID)
	}

	var wantID constant.TDLibPeerID

	wantID.User(10)

	if story.Chat.ID != int64(wantID) || story.Chat.Type != ChatTypePrivate || story.Chat.Username != "biz" {
		t.Fatalf("story chat = %#v", story.Chat)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.StoriesSendStoryRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	req, isSend := wrapper.Query.(*tg.StoriesSendStoryRequest)
	if !isSend {
		t.Fatalf("query = %#v, want send story", wrapper.Query)
	}

	if _, isSelf := req.Peer.(*tg.InputPeerSelf); !isSelf {
		t.Fatalf("peer = %#v, want self", req.Peer)
	}

	if req.Period != 86400 || req.Caption != "hi" || !req.Pinned || !req.Noforwards {
		t.Fatalf("req = %#v", req)
	}

	ext, isExt := req.Media.(*tg.InputMediaPhotoExternal)
	if !isExt || ext.URL != "https://example.com/p.jpg" {
		t.Fatalf("media = %#v", req.Media)
	}

	if len(req.PrivacyRules) != 1 {
		t.Fatalf("privacy rules = %#v", req.PrivacyRules)
	}

	if _, isAll := req.PrivacyRules[0].(*tg.InputPrivacyValueAllowAll); !isAll {
		t.Fatalf("privacy rule = %#v, want allow all", req.PrivacyRules[0])
	}
}

func TestEditStory(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, storyUpdates())

	story, err := newMockBot(inv).EditStory(context.Background(), "bc1", 42,
		InputStoryContentVideo{Video: FileURL("https://example.com/v.mp4"), Duration: 5, IsAnimation: true})
	if err != nil {
		t.Fatalf("EditStory: %v", err)
	}

	if story.ID != 42 {
		t.Fatalf("story id = %d, want 42", story.ID)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.StoriesEditStoryRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.StoriesEditStoryRequest)
	if !ok {
		t.Fatalf("query = %#v, want edit story", wrapper.Query)
	}

	if req.ID != 42 {
		t.Fatalf("story id = %d", req.ID)
	}

	if _, ok := req.Media.(*tg.InputMediaDocumentExternal); !ok {
		t.Fatalf("media = %#v, want external document", req.Media)
	}
}

func TestRepostStory(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, storyUpdates())

	story, err := newMockBot(inv).RepostStory(context.Background(), "bc1", userRef(20, 3), 7, 43200)
	if err != nil {
		t.Fatalf("RepostStory: %v", err)
	}

	if story.ID != 42 {
		t.Fatalf("story id = %d", story.ID)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.StoriesSendStoryRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.StoriesSendStoryRequest)
	if !ok {
		t.Fatalf("query = %#v, want send story", wrapper.Query)
	}

	if req.FwdFromStory != 7 || req.Period != 43200 {
		t.Fatalf("req = %#v", req)
	}

	if _, ok := req.FwdFromID.(*tg.InputPeerUser); !ok {
		t.Fatalf("fwd from = %#v, want user", req.FwdFromID)
	}
}

func TestDeleteStory(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.IntVector{Elems: []int{42}})

	if err := newMockBot(inv).DeleteStory(context.Background(), "bc1", 42); err != nil {
		t.Fatalf("DeleteStory: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.StoriesDeleteStoriesRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	req, ok := wrapper.Query.(*tg.StoriesDeleteStoriesRequest)
	if !ok {
		t.Fatalf("query = %#v, want delete stories", wrapper.Query)
	}

	if len(req.ID) != 1 || req.ID[0] != 42 {
		t.Fatalf("ids = %#v", req.ID)
	}

	if _, ok := req.Peer.(*tg.InputPeerSelf); !ok {
		t.Fatalf("peer = %#v, want self", req.Peer)
	}
}
