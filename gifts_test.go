package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetAvailableGifts(t *testing.T) {
	doc := &tg.Document{
		ID:            5,
		AccessHash:    6,
		FileReference: []byte{1},
		DCID:          2,
		Size:          2048,
		MimeType:      mimeStickerAnimated,
		Attributes: []tg.DocumentAttributeClass{
			&tg.DocumentAttributeImageSize{W: 512, H: 512},
			&tg.DocumentAttributeSticker{Alt: "🎁", Stickerset: &tg.InputStickerSetEmpty{}},
		},
	}

	limited := tg.StarGift{
		ID:             100,
		Sticker:        doc,
		Stars:          50,
		Limited:        true,
		RequirePremium: true,
	}
	limited.SetUpgradeStars(25)
	limited.SetAvailabilityTotal(1000)
	limited.SetAvailabilityRemains(900)
	limited.SetPerUserTotal(5)
	limited.SetPerUserRemains(3)
	limited.SetUpgradeVariants(7)
	limited.SetBackground(tg.StarGiftBackground{CenterColor: 0x112233, EdgeColor: 0x445566, TextColor: 0xffffff})
	limited.SetReleasedBy(&tg.PeerChannel{ChannelID: 777})

	plain := tg.StarGift{
		ID:      200,
		Sticker: doc,
		Stars:   10,
	}

	resp := &tg.PaymentsStarGifts{
		Gifts: []tg.StarGiftClass{
			&limited,
			&tg.StarGiftUnique{ID: 300}, // collectible, must be skipped
			&plain,
		},
		Chats: []tg.ChatClass{
			&tg.Channel{ID: 777, Title: "Publisher", Username: "pub", Broadcast: true, Photo: &tg.ChatPhotoEmpty{}},
		},
	}

	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetStarGiftsRequestTypeID, resp)

	got, err := newMockBot(inv).GetAvailableGifts(context.Background())
	if err != nil {
		t.Fatalf("GetAvailableGifts: %v", err)
	}

	if len(got) != 2 {
		t.Fatalf("len = %d, want 2 (unique gift skipped)", len(got))
	}

	g := got[0]
	if g.ID != "100" || g.StarCount != 50 || g.UpgradeStarCount != 25 {
		t.Fatalf("counts: %#v", g)
	}

	if !g.IsPremium {
		t.Fatalf("is_premium not set: %#v", g)
	}

	if g.TotalCount != 1000 || g.RemainingCount != 900 {
		t.Fatalf("availability: %#v", g)
	}

	if g.PersonalTotalCount != 5 || g.PersonalRemainingCount != 3 {
		t.Fatalf("per-user: %#v", g)
	}

	if g.UniqueGiftVariantCount != 7 {
		t.Fatalf("variants: %#v", g)
	}

	if g.Sticker.FileID == "" || g.Sticker.Emoji != "🎁" || !g.Sticker.IsAnimated {
		t.Fatalf("sticker: %#v", g.Sticker)
	}

	if g.Background == nil || g.Background.CenterColor != 0x112233 || g.Background.TextColor != 0xffffff {
		t.Fatalf("background: %#v", g.Background)
	}

	if g.PublisherChat == nil || g.PublisherChat.Title != "Publisher" || g.PublisherChat.Type != ChatTypeChannel {
		t.Fatalf("publisher: %#v", g.PublisherChat)
	}

	plainGift := got[1]
	if plainGift.ID != "200" || plainGift.UpgradeStarCount != 0 || plainGift.Background != nil || plainGift.PublisherChat != nil {
		t.Fatalf("plain gift: %#v", plainGift)
	}

	var req tg.PaymentsGetStarGiftsRequest

	inv.decode(t, tg.PaymentsGetStarGiftsRequestTypeID, &req)

	if req.Hash != 0 {
		t.Fatalf("hash = %d, want 0", req.Hash)
	}
}
