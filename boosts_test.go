package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestGetUserChatBoosts(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.PremiumGetUserBoostsRequestTypeID, &tg.PremiumBoostsList{
		Count: 2,
		Boosts: []tg.Boost{
			{ID: "b1", UserID: 10, Date: 100, Expires: 200},
			{ID: "b2", UserID: 10, Date: 150, Expires: 250, Giveaway: true, GiveawayMsgID: 7, Unclaimed: true},
		},
		Users: []tg.UserClass{&tg.User{ID: 10, AccessHash: 1, FirstName: "Boo"}},
	})

	boosts, err := newMockBot(inv).GetUserChatBoosts(context.Background(), tdlibChannel(50), 10)
	if err != nil {
		t.Fatalf("GetUserChatBoosts: %v", err)
	}

	if len(boosts) != 2 {
		t.Fatalf("got %d boosts", len(boosts))
	}

	premium, ok := boosts[0].Source.(ChatBoostSourcePremium)
	if !ok || premium.User.FirstName != "Boo" || boosts[0].BoostID != "b1" {
		t.Fatalf("boost[0] = %#v", boosts[0])
	}

	giveaway, ok := boosts[1].Source.(ChatBoostSourceGiveaway)
	if !ok || giveaway.GiveawayMessageID != 7 || !giveaway.IsUnclaimed {
		t.Fatalf("boost[1] = %#v", boosts[1])
	}

	if giveaway.User == nil || giveaway.User.ID != 10 {
		t.Fatalf("giveaway user = %#v", giveaway.User)
	}
}
