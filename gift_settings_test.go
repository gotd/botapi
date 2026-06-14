package botapi

import (
	"context"
	"testing"

	"github.com/gotd/td/tg"
)

func TestSetBusinessAccountGiftSettings(t *testing.T) {
	inv := newMockInvoker()
	// Both the get and the set go through the business wrapper; the reply
	// decodes as GlobalPrivacySettings for either call. HideReadMarks stands in
	// for an unrelated setting that must survive the read-modify-write.
	inv.reply(tg.InvokeWithBusinessConnectionRequestTypeID, &tg.GlobalPrivacySettings{HideReadMarks: true})

	err := newMockBot(inv).SetBusinessAccountGiftSettings(context.Background(), "bc1", true, AcceptedGiftTypes{
		UnlimitedGifts:    true,
		UniqueGifts:       true,
		GiftsFromChannels: true,
	})
	if err != nil {
		t.Fatalf("SetBusinessAccountGiftSettings: %v", err)
	}

	wrapper := tg.InvokeWithBusinessConnectionRequest{Query: &tg.AccountSetGlobalPrivacySettingsRequest{}}

	inv.decode(t, tg.InvokeWithBusinessConnectionRequestTypeID, &wrapper)

	if wrapper.ConnectionID != "bc1" {
		t.Fatalf("connection id = %q", wrapper.ConnectionID)
	}

	req, ok := wrapper.Query.(*tg.AccountSetGlobalPrivacySettingsRequest)
	if !ok {
		t.Fatalf("query = %#v, want set global privacy", wrapper.Query)
	}

	settings := req.Settings
	if !settings.HideReadMarks {
		t.Fatalf("unrelated setting not preserved: %#v", settings)
	}

	if !settings.DisplayGiftsButton {
		t.Fatalf("display gifts button not set")
	}

	disallowed, ok := settings.GetDisallowedGifts()
	if !ok {
		t.Fatalf("disallowed gifts not set")
	}

	if disallowed.DisallowUnlimitedStargifts {
		t.Fatalf("unlimited gifts should be allowed")
	}

	if !disallowed.DisallowLimitedStargifts {
		t.Fatalf("limited gifts should be disallowed")
	}

	if !disallowed.DisallowPremiumGifts {
		t.Fatalf("premium gifts should be disallowed")
	}

	if disallowed.DisallowUniqueStargifts || disallowed.DisallowStargiftsFromChannels {
		t.Fatalf("unique/channel gifts should be allowed: %#v", disallowed)
	}
}
