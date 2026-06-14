package botapi

import (
	"context"
	"strconv"

	"github.com/gotd/td/tg"
)

// RevenueWithdrawalState is a sealed union describing the state of a revenue
// withdrawal operation.
//
// Concrete variants: RevenueWithdrawalStatePending,
// RevenueWithdrawalStateSucceeded, RevenueWithdrawalStateFailed.
type RevenueWithdrawalState interface {
	isRevenueWithdrawalState()
}

// RevenueWithdrawalStatePending reports that the withdrawal is in progress.
type RevenueWithdrawalStatePending struct {
	Type string `json:"type"`
}

// RevenueWithdrawalStateSucceeded reports that the withdrawal succeeded.
type RevenueWithdrawalStateSucceeded struct {
	Type string `json:"type"`
	// Date is the Unix time when the withdrawal was completed.
	Date int `json:"date"`
	// URL is an HTTPS URL that can be used to see transaction details.
	URL string `json:"url"`
}

// RevenueWithdrawalStateFailed reports that the withdrawal failed and the
// transaction was refunded.
type RevenueWithdrawalStateFailed struct {
	Type string `json:"type"`
}

func (RevenueWithdrawalStatePending) isRevenueWithdrawalState()   {}
func (RevenueWithdrawalStateSucceeded) isRevenueWithdrawalState() {}
func (RevenueWithdrawalStateFailed) isRevenueWithdrawalState()    {}

// Revenue withdrawal state discriminators.
const (
	revenueWithdrawalStatePending   = "pending"
	revenueWithdrawalStateSucceeded = "succeeded"
	revenueWithdrawalStateFailed    = "failed"
)

// TransactionPartner is a sealed union describing the source of a transaction,
// or its recipient for outgoing transactions.
//
// Concrete variants: TransactionPartnerUser, TransactionPartnerChat,
// TransactionPartnerAffiliateProgram, TransactionPartnerFragment,
// TransactionPartnerTelegramAds, TransactionPartnerTelegramApi,
// TransactionPartnerOther.
type TransactionPartner interface {
	isTransactionPartner()
}

// TransactionPartnerUser describes a transaction with a user.
type TransactionPartnerUser struct {
	Type string `json:"type"`
	// TransactionType is the type of the transaction; one of "invoice_payment",
	// "paid_media_payment", "gift_purchase", "premium_purchase" or
	// "business_account_transfer".
	TransactionType string `json:"transaction_type"`
	// User is the information about the user.
	User User `json:"user"`
	// InvoicePayload is the bot-specified invoice payload, for "invoice_payment".
	InvoicePayload string `json:"invoice_payload,omitempty"`
	// SubscriptionPeriod is the duration, in seconds, of the paid subscription.
	SubscriptionPeriod int `json:"subscription_period,omitempty"`
	// PremiumSubscriptionDuration is the number of months the gifted Telegram
	// Premium subscription will be active for, for "premium_purchase".
	PremiumSubscriptionDuration int `json:"premium_subscription_duration,omitempty"`
}

// TransactionPartnerChat describes a transaction with a chat.
type TransactionPartnerChat struct {
	Type string `json:"type"`
	// Chat is the information about the chat.
	Chat Chat `json:"chat"`
}

// TransactionPartnerAffiliateProgram describes the affiliate program that
// issued the affiliate commission received via this transaction.
type TransactionPartnerAffiliateProgram struct {
	Type string `json:"type"`
	// SponsorUser is the information about the bot that sponsored the affiliate
	// program, if any.
	SponsorUser *User `json:"sponsor_user,omitempty"`
	// CommissionPerMille is the number of Telegram Stars received by the bot for
	// each 1000 Telegram Stars received by the affiliate program sponsor.
	CommissionPerMille int `json:"commission_per_mille"`
}

// TransactionPartnerFragment describes a withdrawal transaction with Fragment.
type TransactionPartnerFragment struct {
	Type string `json:"type"`
	// WithdrawalState is the state of the transaction if the transaction is
	// outgoing.
	WithdrawalState RevenueWithdrawalState `json:"withdrawal_state,omitempty"`
}

// TransactionPartnerTelegramAds describes a withdrawal transaction to the
// Telegram Ads platform.
type TransactionPartnerTelegramAds struct {
	Type string `json:"type"`
}

// TransactionPartnerTelegramApi describes a transaction with payment for paid
// broadcasting.
type TransactionPartnerTelegramApi struct {
	Type string `json:"type"`
	// RequestCount is the number of successful requests that exceeded regular
	// limits and were therefore billed.
	RequestCount int `json:"request_count"`
}

// TransactionPartnerOther describes a transaction with an unknown source or
// recipient.
type TransactionPartnerOther struct {
	Type string `json:"type"`
}

func (TransactionPartnerUser) isTransactionPartner()             {}
func (TransactionPartnerChat) isTransactionPartner()             {}
func (TransactionPartnerAffiliateProgram) isTransactionPartner() {}
func (TransactionPartnerFragment) isTransactionPartner()         {}
func (TransactionPartnerTelegramAds) isTransactionPartner()      {}
func (TransactionPartnerTelegramApi) isTransactionPartner()      {}
func (TransactionPartnerOther) isTransactionPartner()            {}

// Transaction partner discriminators.
const (
	transactionPartnerUser             = "user"
	transactionPartnerChat             = "chat"
	transactionPartnerAffiliateProgram = "affiliate_program"
	transactionPartnerFragment         = "fragment"
	transactionPartnerTelegramAds      = "telegram_ads"
	transactionPartnerTelegramApi      = "telegram_api"
	transactionPartnerOther            = "other"
)

// Transaction type discriminators for TransactionPartnerUser.
const (
	transactionTypeInvoicePayment         = "invoice_payment"
	transactionTypePaidMediaPayment       = "paid_media_payment"
	transactionTypeGiftPurchase           = "gift_purchase"
	transactionTypePremiumPurchase        = "premium_purchase"
	transactionTypeBusinessAccountTranfer = "business_account_transfer"
)

// StarTransaction describes a Telegram Star transaction.
type StarTransaction struct {
	// ID is the unique identifier of the transaction. Coincides with the
	// identifier of the original transaction for refund transactions.
	ID string `json:"id"`
	// Amount is the integer amount of Telegram Stars transferred by the
	// transaction.
	Amount int `json:"amount"`
	// NanostarAmount is the number of 1/1000000000 shares of Telegram Stars
	// transferred by the transaction; from 0 to 999999999.
	NanostarAmount int `json:"nanostar_amount,omitempty"`
	// Date is the Unix time when the transaction was created.
	Date int `json:"date"`
	// Source is the source of an incoming transaction (for transactions with a
	// positive amount); only for incoming transactions.
	Source TransactionPartner `json:"source,omitempty"`
	// Receiver is the receiver of an outgoing transaction (for transactions
	// with a negative amount); only for outgoing transactions.
	Receiver TransactionPartner `json:"receiver,omitempty"`
}

// GetStarTransactions returns the bot's Telegram Star transactions in
// chronological order. offset is the number of transactions to skip; limit is
// the maximum number of transactions to return (1-100, defaults to 100).
func (b *Bot) GetStarTransactions(ctx context.Context, offset, limit int) ([]StarTransaction, error) {
	if limit <= 0 {
		limit = 100
	}

	req := &tg.PaymentsGetStarsTransactionsRequest{
		Peer:  &tg.InputPeerSelf{},
		Limit: limit,
	}
	if offset > 0 {
		req.Offset = strconv.Itoa(offset)
	}

	res, err := b.raw.PaymentsGetStarsTransactions(ctx, req)
	if err != nil {
		return nil, asAPIError(err)
	}

	users := usersByID(res.Users)
	chats := chatsByID(res.Chats)

	out := make([]StarTransaction, 0, len(res.History))
	for i := range res.History {
		out = append(out, starTransactionFromTg(res.History[i], users, chats))
	}

	return out, nil
}

// starTransactionFromTg converts a raw MTProto stars transaction into the Bot
// API StarTransaction. The sign of the amount selects whether the partner is
// the source (incoming) or the receiver (outgoing).
func starTransactionFromTg(tx tg.StarsTransaction, users map[int64]*tg.User, chats map[int64]tg.ChatClass) StarTransaction {
	amount, nanos := starsAmountValues(tx.Amount)

	out := StarTransaction{
		ID:   tx.ID,
		Date: tx.Date,
	}

	partner := transactionPartnerFromTg(tx, users, chats)

	// A positive amount is money coming in (source); a negative amount is money
	// going out (receiver). The Bot API always reports a non-negative amount.
	if amount > 0 || nanos > 0 {
		out.Amount = amount
		out.NanostarAmount = nanos
		out.Source = partner
	} else {
		out.Amount = -amount
		out.NanostarAmount = -nanos
		out.Receiver = partner
	}

	return out
}

// transactionPartnerFromTg maps a raw stars-transaction peer (plus the
// transaction flags) onto the Bot API TransactionPartner union.
func transactionPartnerFromTg(tx tg.StarsTransaction, users map[int64]*tg.User, chats map[int64]tg.ChatClass) TransactionPartner {
	switch peer := tx.Peer.(type) {
	case *tg.StarsTransactionPeer:
		return partnerForPeer(tx, peer.Peer, users, chats)
	case *tg.StarsTransactionPeerFragment:
		return TransactionPartnerFragment{
			Type:            transactionPartnerFragment,
			WithdrawalState: withdrawalStateFromTg(tx),
		}
	case *tg.StarsTransactionPeerAds:
		return TransactionPartnerTelegramAds{Type: transactionPartnerTelegramAds}
	case *tg.StarsTransactionPeerAPI:
		return TransactionPartnerTelegramApi{
			Type:         transactionPartnerTelegramApi,
			RequestCount: tx.FloodskipNumber,
		}
	default:
		// AppStore, PlayMarket, PremiumBot and Unsupported peers have no Bot API
		// equivalent.
		return TransactionPartnerOther{Type: transactionPartnerOther}
	}
}

// partnerForPeer builds the user/chat partner for a starsTransactionPeer.
func partnerForPeer(tx tg.StarsTransaction, p tg.PeerClass, users map[int64]*tg.User, chats map[int64]tg.ChatClass) TransactionPartner {
	switch p := p.(type) {
	case *tg.PeerUser:
		var user User

		if u, ok := users[p.UserID]; ok {
			user = userFromTgUser(u)
		} else {
			user = User{ID: p.UserID}
		}

		return TransactionPartnerUser{
			Type:                        transactionPartnerUser,
			TransactionType:             userTransactionType(tx),
			User:                        user,
			InvoicePayload:              string(tx.BotPayload),
			SubscriptionPeriod:          tx.SubscriptionPeriod,
			PremiumSubscriptionDuration: tx.PremiumGiftMonths,
		}
	case *tg.PeerChat:
		return TransactionPartnerChat{
			Type: transactionPartnerChat,
			Chat: chatFromRaw(p, chats),
		}
	case *tg.PeerChannel:
		return TransactionPartnerChat{
			Type: transactionPartnerChat,
			Chat: chatFromRaw(p, chats),
		}
	default:
		return TransactionPartnerOther{Type: transactionPartnerOther}
	}
}

// userTransactionType derives the Bot API transaction_type for a user partner
// from the raw transaction flags.
func userTransactionType(tx tg.StarsTransaction) string {
	switch {
	case tx.BusinessTransfer:
		return transactionTypeBusinessAccountTranfer
	case tx.PremiumGiftMonths > 0:
		return transactionTypePremiumPurchase
	case len(tx.ExtendedMedia) > 0:
		return transactionTypePaidMediaPayment
	case isGiftTransaction(tx):
		return transactionTypeGiftPurchase
	default:
		return transactionTypeInvoicePayment
	}
}

// isGiftTransaction reports whether the transaction is a gift purchase.
func isGiftTransaction(tx tg.StarsTransaction) bool {
	_, ok := tx.GetStargift()

	return ok || tx.Gift
}

// withdrawalStateFromTg derives the Fragment withdrawal state from the raw
// transaction flags and fields.
func withdrawalStateFromTg(tx tg.StarsTransaction) RevenueWithdrawalState {
	switch {
	case tx.Pending:
		return RevenueWithdrawalStatePending{Type: revenueWithdrawalStatePending}
	case tx.Failed:
		return RevenueWithdrawalStateFailed{Type: revenueWithdrawalStateFailed}
	default:
		return RevenueWithdrawalStateSucceeded{
			Type: revenueWithdrawalStateSucceeded,
			Date: tx.TransactionDate,
			URL:  tx.TransactionURL,
		}
	}
}

// starsAmountValues extracts the integer and nanostar parts from a stars
// amount. Non-star (TON) amounts are reported as zero.
func starsAmountValues(a tg.StarsAmountClass) (amount, nanos int) {
	if v, ok := a.(*tg.StarsAmount); ok {
		return int(v.Amount), v.Nanos
	}

	return 0, 0
}
