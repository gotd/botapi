package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func giftStickerDoc() *tg.Document {
	return &tg.Document{
		ID:            5,
		AccessHash:    6,
		FileReference: []byte{1},
		DCID:          2,
		Size:          1024,
		MimeType:      mimeStickerAnimated,
		Attributes: []tg.DocumentAttributeClass{
			&tg.DocumentAttributeSticker{Alt: "🎁", Stickerset: &tg.InputStickerSetEmpty{}},
		},
	}
}

func regularSavedGift(msgID int) tg.SavedStarGift {
	sg := tg.SavedStarGift{
		Date: 100,
		Gift: &tg.StarGift{ID: 100, Sticker: giftStickerDoc(), Stars: 50},
	}
	sg.SetMsgID(msgID)
	sg.SetConvertStars(40)
	sg.SetMessage(tg.TextWithEntities{Text: "thanks"})

	return sg
}

func TestGetUserGifts(t *testing.T) {
	resp := &tg.PaymentsSavedStarGifts{
		Count: 1,
		Gifts: []tg.SavedStarGift{regularSavedGift(7)},
	}
	resp.SetNextOffset("next")

	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetSavedStarGiftsRequestTypeID, resp)

	got, err := newMockBot(inv).GetUserGifts(context.Background(), 42, WithGiftListLimit(10), WithGiftListExcludeUnique())
	if err != nil {
		t.Fatalf("GetUserGifts: %v", err)
	}

	if got.TotalCount != 1 || got.NextOffset != "next" || len(got.Gifts) != 1 {
		t.Fatalf("owned gifts = %#v", got)
	}

	reg, ok := got.Gifts[0].(OwnedGiftRegular)
	if !ok {
		t.Fatalf("gift = %#v, want regular", got.Gifts[0])
	}

	if reg.OwnedGiftID != OwnedGiftFromMessage(7) {
		t.Fatalf("owned_gift_id = %q", reg.OwnedGiftID)
	}

	if reg.Gift.ID != "100" || reg.ConvertStarCount != 40 || reg.Text != "thanks" {
		t.Fatalf("regular gift = %#v", reg)
	}

	var req tg.PaymentsGetSavedStarGiftsRequest

	inv.decode(t, tg.PaymentsGetSavedStarGiftsRequestTypeID, &req)

	if _, ok := req.Peer.(*tg.InputPeerUser); !ok {
		t.Fatalf("peer = %#v, want user", req.Peer)
	}

	if req.Limit != 10 || !req.ExcludeUnique {
		t.Fatalf("req = %#v", req)
	}
}

func TestGetChatGifts(t *testing.T) {
	sg := tg.SavedStarGift{
		Date: 200,
		Gift: &tg.StarGift{ID: 100, Sticker: giftStickerDoc(), Stars: 50},
	}
	sg.SetSavedID(9)

	inv := newMockInvoker()
	inv.reply(tg.PaymentsGetSavedStarGiftsRequestTypeID, &tg.PaymentsSavedStarGifts{
		Count: 1,
		Gifts: []tg.SavedStarGift{sg},
	})

	got, err := newMockBot(inv).GetChatGifts(context.Background(), channelRef(777, 888))
	if err != nil {
		t.Fatalf("GetChatGifts: %v", err)
	}

	reg, ok := got.Gifts[0].(OwnedGiftRegular)
	if !ok {
		t.Fatalf("gift = %#v, want regular", got.Gifts[0])
	}

	if reg.OwnedGiftID != "chat:777:888:9" {
		t.Fatalf("owned_gift_id = %q", reg.OwnedGiftID)
	}

	var req tg.PaymentsGetSavedStarGiftsRequest

	inv.decode(t, tg.PaymentsGetSavedStarGiftsRequestTypeID, &req)

	if _, ok := req.Peer.(*tg.InputPeerChannel); !ok {
		t.Fatalf("peer = %#v, want channel", req.Peer)
	}
}

func TestGetBusinessAccountGifts(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.PaymentsSavedStarGifts{Count: 0})

	got, err := newMockBot(inv).GetBusinessAccountGifts(context.Background(), "bc1", WithGiftListSortByPrice())
	if err != nil {
		t.Fatalf("GetBusinessAccountGifts: %v", err)
	}

	if got.TotalCount != 0 || len(got.Gifts) != 0 {
		t.Fatalf("owned gifts = %#v", got)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.PaymentsGetSavedStarGiftsRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	req, ok := wrapper.Query.(*tg.PaymentsGetSavedStarGiftsRequest)
	if !ok {
		t.Fatalf("query = %#v", wrapper.Query)
	}

	if _, ok := req.Peer.(*tg.InputPeerSelf); !ok {
		t.Fatalf("peer = %#v, want self", req.Peer)
	}

	if !req.SortByValue {
		t.Fatalf("sort by value not set")
	}
}

func TestOwnedGiftUniqueFromTg(t *testing.T) {
	unique := &tg.StarGiftUnique{
		ID:             1,
		GiftID:         100,
		Title:          "Star",
		Slug:           "Star-7",
		Num:            7,
		RequirePremium: true,
		Attributes: []tg.StarGiftAttributeClass{
			&tg.StarGiftAttributeModel{
				Name:     "Model A",
				Document: giftStickerDoc(),
				Rarity:   &tg.StarGiftAttributeRarity{Permille: 15},
			},
			&tg.StarGiftAttributePattern{
				Name:     "Pattern B",
				Document: giftStickerDoc(),
				Rarity:   &tg.StarGiftAttributeRarityRare{},
			},
			&tg.StarGiftAttributeBackdrop{
				Name:         "Sky",
				CenterColor:  0x010203,
				EdgeColor:    0x040506,
				PatternColor: 0x070809,
				TextColor:    0x0a0b0c,
				Rarity:       &tg.StarGiftAttributeRarity{Permille: 5},
			},
		},
	}

	sg := tg.SavedStarGift{Date: 300, Gift: unique}
	sg.SetMsgID(11)
	sg.SetFromID(&tg.PeerUser{UserID: 50})
	sg.SetTransferStars(99)
	sg.SetCanTransferAt(1234)

	users := map[int64]*tg.User{50: {ID: 50, FirstName: "Sender"}}

	got := ownedGiftFromTg(sg, users, nil, 0, 0)

	uniq, ok := got.(OwnedGiftUnique)
	if !ok {
		t.Fatalf("gift = %#v, want unique", got)
	}

	if uniq.OwnedGiftID != OwnedGiftFromMessage(11) {
		t.Fatalf("owned_gift_id = %q", uniq.OwnedGiftID)
	}

	if uniq.SenderUser == nil || uniq.SenderUser.FirstName != "Sender" {
		t.Fatalf("sender = %#v", uniq.SenderUser)
	}

	if uniq.TransferStarCount != 99 || !uniq.CanBeTransferred || uniq.NextTransferDate != 1234 {
		t.Fatalf("transfer = %#v", uniq)
	}

	g := uniq.Gift
	if g.GiftID != "100" || g.BaseName != "Star" || g.Name != "Star-7" || g.Number != 7 || !g.IsPremium {
		t.Fatalf("unique gift = %#v", g)
	}

	if g.Model.Name != "Model A" || g.Model.RarityPerMille != 15 || g.Model.Sticker.FileID == "" {
		t.Fatalf("model = %#v", g.Model)
	}

	if g.Symbol.Name != "Pattern B" || g.Symbol.Sticker.FileID == "" {
		t.Fatalf("symbol = %#v", g.Symbol)
	}

	if g.Backdrop.Name != "Sky" || g.Backdrop.RarityPerMille != 5 || g.Backdrop.Colors.SymbolColor != 0x070809 {
		t.Fatalf("backdrop = %#v", g.Backdrop)
	}
}
