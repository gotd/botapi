package botapi

import (
	"context"

	"github.com/gotd/td/telegram/message"
	"github.com/gotd/td/tg"
)

// InvoiceParams describes an invoice to send with SendInvoice.
type InvoiceParams struct {
	Title         string         `json:"title"`
	Description   string         `json:"description"`
	Payload       string         `json:"payload"` // bot-defined, not shown to the user
	ProviderToken string         `json:"provider_token"`
	Currency      string         `json:"currency"`
	Prices        []LabeledPrice `json:"prices"`

	MaxTipAmount        int    `json:"max_tip_amount,omitempty"`
	SuggestedTipAmounts []int  `json:"suggested_tip_amounts,omitempty"`
	StartParameter      string `json:"start_parameter,omitempty"`
	ProviderData        string `json:"provider_data,omitempty"` // JSON for the payment provider

	PhotoURL    string `json:"photo_url,omitempty"`
	PhotoSize   int    `json:"photo_size,omitempty"`
	PhotoWidth  int    `json:"photo_width,omitempty"`
	PhotoHeight int    `json:"photo_height,omitempty"`

	NeedName            bool `json:"need_name,omitempty"`
	NeedPhoneNumber     bool `json:"need_phone_number,omitempty"`
	NeedEmail           bool `json:"need_email,omitempty"`
	NeedShippingAddress bool `json:"need_shipping_address,omitempty"`

	SendPhoneNumberToProvider bool `json:"send_phone_number_to_provider,omitempty"`
	SendEmailToProvider       bool `json:"send_email_to_provider,omitempty"`
	IsFlexible                bool `json:"is_flexible,omitempty"`
}

// invoiceMedia builds the MTProto invoice media shared by SendInvoice and
// CreateInvoiceLink.
func invoiceMedia(params InvoiceParams) *tg.InputMediaInvoice {
	suggested := make([]int64, len(params.SuggestedTipAmounts))
	for i, a := range params.SuggestedTipAmounts {
		suggested[i] = int64(a)
	}

	invoice := tg.Invoice{
		NameRequested:            params.NeedName,
		PhoneRequested:           params.NeedPhoneNumber,
		EmailRequested:           params.NeedEmail,
		ShippingAddressRequested: params.NeedShippingAddress,
		Flexible:                 params.IsFlexible,
		PhoneToProvider:          params.SendPhoneNumberToProvider,
		EmailToProvider:          params.SendEmailToProvider,
		Currency:                 params.Currency,
		Prices:                   pricesToTg(params.Prices),
		MaxTipAmount:             int64(params.MaxTipAmount),
		SuggestedTipAmounts:      suggested,
	}

	media := &tg.InputMediaInvoice{
		Title:       params.Title,
		Description: params.Description,
		Invoice:     invoice,
		Payload:     []byte(params.Payload),
		Provider:    params.ProviderToken,
		StartParam:  params.StartParameter,
	}
	if params.ProviderData != "" {
		media.ProviderData = tg.DataJSON{Data: params.ProviderData}
	}

	if params.PhotoURL != "" {
		media.SetPhoto(tg.InputWebDocument{
			URL:      params.PhotoURL,
			Size:     params.PhotoSize,
			MimeType: "image/jpeg",
			Attributes: []tg.DocumentAttributeClass{
				&tg.DocumentAttributeImageSize{W: params.PhotoWidth, H: params.PhotoHeight},
			},
		})
	}

	return media
}

// SendInvoice sends an invoice.
func (b *Bot) SendInvoice(ctx context.Context, chat ChatID, params InvoiceParams, opts ...SendOption) (*Message, error) {
	var cfg sendConfig

	for _, o := range opts {
		o(&cfg)
	}

	peer, err := b.resolveInputPeer(ctx, chat)
	if err != nil {
		return nil, err
	}

	media := invoiceMedia(params)

	builder := &b.sender.To(peer).Builder

	builder, err = b.applySendConfig(builder, cfg)
	if err != nil {
		return nil, err
	}

	resp, err := builder.Media(ctx, message.Media(media))

	return b.sentMessage(ctx, peer, resp, err)
}

// pricesToTg converts Bot API labeled prices to MTProto.
func pricesToTg(prices []LabeledPrice) []tg.LabeledPrice {
	out := make([]tg.LabeledPrice, 0, len(prices))
	for _, p := range prices {
		out = append(out, tg.LabeledPrice{Label: p.Label, Amount: int64(p.Amount)})
	}

	return out
}
