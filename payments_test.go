package botapi

import (
	"testing"

	"github.com/gotd/td/tg"
)

func TestShippingQueryFromTg(t *testing.T) {
	e := tg.Entities{Users: map[int64]*tg.User{7: {ID: 7, FirstName: "Ann"}}}
	u := &tg.UpdateBotShippingQuery{
		QueryID: 100,
		UserID:  7,
		Payload: []byte("order-1"),
		ShippingAddress: tg.PostAddress{
			StreetLine1: "1 Main St",
			City:        "Town",
			CountryISO2: "US",
			PostCode:    "12345",
		},
	}
	q := shippingQueryFromTg(e, u)

	if q.ID != "100" || q.From.FirstName != "Ann" || q.InvoicePayload != "order-1" {
		t.Fatalf("query: %#v", q)
	}

	if q.ShippingAddress.CountryCode != "US" || q.ShippingAddress.City != "Town" {
		t.Fatalf("address: %#v", q.ShippingAddress)
	}
}

func TestPreCheckoutQueryFromTg(t *testing.T) {
	u := &tg.UpdateBotPrecheckoutQuery{
		QueryID:          200,
		UserID:           9,
		Payload:          []byte("p"),
		Currency:         "USD",
		TotalAmount:      1500,
		ShippingOptionID: "express",
	}
	u.SetInfo(tg.PaymentRequestedInfo{Name: "Bob", Email: "b@x.io"})

	q := preCheckoutQueryFromTg(tg.Entities{}, u)
	if q.ID != "200" || q.From.ID != 9 || q.Currency != "USD" || q.TotalAmount != 1500 {
		t.Fatalf("query: %#v", q)
	}

	if q.ShippingOptionID != "express" {
		t.Fatalf("shipping option: %q", q.ShippingOptionID)
	}

	if q.OrderInfo == nil || q.OrderInfo.Name != "Bob" || q.OrderInfo.Email != "b@x.io" {
		t.Fatalf("order info: %#v", q.OrderInfo)
	}
}

func TestShippingOptionsToTg(t *testing.T) {
	opts := []ShippingOption{
		{ID: "std", Title: "Standard", Prices: []LabeledPrice{{Label: "Base", Amount: 500}}},
	}
	got := shippingOptionsToTg(opts)

	if len(got) != 1 || got[0].ID != "std" || len(got[0].Prices) != 1 || got[0].Prices[0].Amount != 500 {
		t.Fatalf("converted: %#v", got)
	}
}
