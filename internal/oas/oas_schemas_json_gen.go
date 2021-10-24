// Code generated by ogen, DO NOT EDIT.

package oas

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ogen-go/ogen/conv"
	ht "github.com/ogen-go/ogen/http"
	"github.com/ogen-go/ogen/json"
	"github.com/ogen-go/ogen/uri"
	"github.com/ogen-go/ogen/validate"
)

// No-op definition for keeping imports.
var (
	_ = context.Background()
	_ = fmt.Stringer(nil)
	_ = strings.Builder{}
	_ = errors.Is
	_ = sort.Ints
	_ = chi.Context{}
	_ = http.MethodGet
	_ = io.Copy
	_ = json.Marshal
	_ = bytes.NewReader
	_ = strconv.ParseInt
	_ = time.Time{}
	_ = conv.ToInt32
	_ = uuid.UUID{}
	_ = uri.PathEncoder{}
	_ = url.URL{}
	_ = math.Mod
	_ = validate.Int{}
	_ = ht.NewRequest
	_ = net.IP{}
)

func (AddStickerToSet) WriteJSON(j *json.Stream)        {}
func (AddStickerToSet) ReadJSON(i *json.Iterator) error { return nil }
func (AddStickerToSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (AddStickerToSet) WriteJSONTo(w io.Writer) error   { return nil }

func (AnswerCallbackQuery) WriteJSON(j *json.Stream)        {}
func (AnswerCallbackQuery) ReadJSON(i *json.Iterator) error { return nil }
func (AnswerCallbackQuery) ReadJSONFrom(r io.Reader) error  { return nil }
func (AnswerCallbackQuery) WriteJSONTo(w io.Writer) error   { return nil }

func (AnswerInlineQuery) WriteJSON(j *json.Stream)        {}
func (AnswerInlineQuery) ReadJSON(i *json.Iterator) error { return nil }
func (AnswerInlineQuery) ReadJSONFrom(r io.Reader) error  { return nil }
func (AnswerInlineQuery) WriteJSONTo(w io.Writer) error   { return nil }

func (AnswerPreCheckoutQuery) WriteJSON(j *json.Stream)        {}
func (AnswerPreCheckoutQuery) ReadJSON(i *json.Iterator) error { return nil }
func (AnswerPreCheckoutQuery) ReadJSONFrom(r io.Reader) error  { return nil }
func (AnswerPreCheckoutQuery) WriteJSONTo(w io.Writer) error   { return nil }

func (AnswerShippingQuery) WriteJSON(j *json.Stream)        {}
func (AnswerShippingQuery) ReadJSON(i *json.Iterator) error { return nil }
func (AnswerShippingQuery) ReadJSONFrom(r io.Reader) error  { return nil }
func (AnswerShippingQuery) WriteJSONTo(w io.Writer) error   { return nil }

func (BanChatMember) WriteJSON(j *json.Stream)        {}
func (BanChatMember) ReadJSON(i *json.Iterator) error { return nil }
func (BanChatMember) ReadJSONFrom(r io.Reader) error  { return nil }
func (BanChatMember) WriteJSONTo(w io.Writer) error   { return nil }

// WriteJSON implements json.Marshaler.
func (s ChatPermissions) WriteJSON(j *json.Stream) {
	j.WriteObjectStart()
	more := json.NewMore(j)
	defer more.Reset()
	if s.CanAddWebPagePreviews.Set {
		more.More()
		j.WriteObjectField("can_add_web_page_previews")
		s.CanAddWebPagePreviews.WriteJSON(j)
	}
	if s.CanChangeInfo.Set {
		more.More()
		j.WriteObjectField("can_change_info")
		s.CanChangeInfo.WriteJSON(j)
	}
	if s.CanInviteUsers.Set {
		more.More()
		j.WriteObjectField("can_invite_users")
		s.CanInviteUsers.WriteJSON(j)
	}
	if s.CanPinMessages.Set {
		more.More()
		j.WriteObjectField("can_pin_messages")
		s.CanPinMessages.WriteJSON(j)
	}
	if s.CanSendMediaMessages.Set {
		more.More()
		j.WriteObjectField("can_send_media_messages")
		s.CanSendMediaMessages.WriteJSON(j)
	}
	if s.CanSendMessages.Set {
		more.More()
		j.WriteObjectField("can_send_messages")
		s.CanSendMessages.WriteJSON(j)
	}
	if s.CanSendOtherMessages.Set {
		more.More()
		j.WriteObjectField("can_send_other_messages")
		s.CanSendOtherMessages.WriteJSON(j)
	}
	if s.CanSendPolls.Set {
		more.More()
		j.WriteObjectField("can_send_polls")
		s.CanSendPolls.WriteJSON(j)
	}
	j.WriteObjectEnd()
}

// WriteJSONTo writes ChatPermissions json value to io.Writer.
func (s ChatPermissions) WriteJSONTo(w io.Writer) error {
	j := json.GetStream(w)
	defer json.PutStream(j)
	s.WriteJSON(j)
	return j.Flush()
}

// ReadJSONFrom reads ChatPermissions json value from io.Reader.
func (s *ChatPermissions) ReadJSONFrom(r io.Reader) error {
	buf := json.GetBuffer()
	defer json.PutBuffer(buf)

	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}
	i := json.GetIterator()
	i.ResetBytes(buf.Bytes())
	defer json.PutIterator(i)

	return s.ReadJSON(i)
}

// ReadJSON reads ChatPermissions from json stream.
func (s *ChatPermissions) ReadJSON(i *json.Iterator) error {
	i.ReadObjectCB(func(i *json.Iterator, k string) bool {
		switch k {
		case "can_add_web_page_previews":
			s.CanAddWebPagePreviews.Reset()
			if err := s.CanAddWebPagePreviews.ReadJSON(i); err != nil {
				i.ReportError("Field CanAddWebPagePreviews", err.Error())
				return false
			}
			return true
		case "can_change_info":
			s.CanChangeInfo.Reset()
			if err := s.CanChangeInfo.ReadJSON(i); err != nil {
				i.ReportError("Field CanChangeInfo", err.Error())
				return false
			}
			return true
		case "can_invite_users":
			s.CanInviteUsers.Reset()
			if err := s.CanInviteUsers.ReadJSON(i); err != nil {
				i.ReportError("Field CanInviteUsers", err.Error())
				return false
			}
			return true
		case "can_pin_messages":
			s.CanPinMessages.Reset()
			if err := s.CanPinMessages.ReadJSON(i); err != nil {
				i.ReportError("Field CanPinMessages", err.Error())
				return false
			}
			return true
		case "can_send_media_messages":
			s.CanSendMediaMessages.Reset()
			if err := s.CanSendMediaMessages.ReadJSON(i); err != nil {
				i.ReportError("Field CanSendMediaMessages", err.Error())
				return false
			}
			return true
		case "can_send_messages":
			s.CanSendMessages.Reset()
			if err := s.CanSendMessages.ReadJSON(i); err != nil {
				i.ReportError("Field CanSendMessages", err.Error())
				return false
			}
			return true
		case "can_send_other_messages":
			s.CanSendOtherMessages.Reset()
			if err := s.CanSendOtherMessages.ReadJSON(i); err != nil {
				i.ReportError("Field CanSendOtherMessages", err.Error())
				return false
			}
			return true
		case "can_send_polls":
			s.CanSendPolls.Reset()
			if err := s.CanSendPolls.ReadJSON(i); err != nil {
				i.ReportError("Field CanSendPolls", err.Error())
				return false
			}
			return true
		default:
			i.Skip()
			return true
		}
	})
	return i.Error
}

func (CopyMessage) WriteJSON(j *json.Stream)        {}
func (CopyMessage) ReadJSON(i *json.Iterator) error { return nil }
func (CopyMessage) ReadJSONFrom(r io.Reader) error  { return nil }
func (CopyMessage) WriteJSONTo(w io.Writer) error   { return nil }

func (CreateChatInviteLink) WriteJSON(j *json.Stream)        {}
func (CreateChatInviteLink) ReadJSON(i *json.Iterator) error { return nil }
func (CreateChatInviteLink) ReadJSONFrom(r io.Reader) error  { return nil }
func (CreateChatInviteLink) WriteJSONTo(w io.Writer) error   { return nil }

func (CreateNewStickerSet) WriteJSON(j *json.Stream)        {}
func (CreateNewStickerSet) ReadJSON(i *json.Iterator) error { return nil }
func (CreateNewStickerSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (CreateNewStickerSet) WriteJSONTo(w io.Writer) error   { return nil }

func (DeleteChatPhoto) WriteJSON(j *json.Stream)        {}
func (DeleteChatPhoto) ReadJSON(i *json.Iterator) error { return nil }
func (DeleteChatPhoto) ReadJSONFrom(r io.Reader) error  { return nil }
func (DeleteChatPhoto) WriteJSONTo(w io.Writer) error   { return nil }

func (DeleteChatStickerSet) WriteJSON(j *json.Stream)        {}
func (DeleteChatStickerSet) ReadJSON(i *json.Iterator) error { return nil }
func (DeleteChatStickerSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (DeleteChatStickerSet) WriteJSONTo(w io.Writer) error   { return nil }

func (DeleteMessage) WriteJSON(j *json.Stream)        {}
func (DeleteMessage) ReadJSON(i *json.Iterator) error { return nil }
func (DeleteMessage) ReadJSONFrom(r io.Reader) error  { return nil }
func (DeleteMessage) WriteJSONTo(w io.Writer) error   { return nil }

func (DeleteMyCommands) WriteJSON(j *json.Stream)        {}
func (DeleteMyCommands) ReadJSON(i *json.Iterator) error { return nil }
func (DeleteMyCommands) ReadJSONFrom(r io.Reader) error  { return nil }
func (DeleteMyCommands) WriteJSONTo(w io.Writer) error   { return nil }

func (DeleteStickerFromSet) WriteJSON(j *json.Stream)        {}
func (DeleteStickerFromSet) ReadJSON(i *json.Iterator) error { return nil }
func (DeleteStickerFromSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (DeleteStickerFromSet) WriteJSONTo(w io.Writer) error   { return nil }

func (DeleteWebhook) WriteJSON(j *json.Stream)        {}
func (DeleteWebhook) ReadJSON(i *json.Iterator) error { return nil }
func (DeleteWebhook) ReadJSONFrom(r io.Reader) error  { return nil }
func (DeleteWebhook) WriteJSONTo(w io.Writer) error   { return nil }

func (EditChatInviteLink) WriteJSON(j *json.Stream)        {}
func (EditChatInviteLink) ReadJSON(i *json.Iterator) error { return nil }
func (EditChatInviteLink) ReadJSONFrom(r io.Reader) error  { return nil }
func (EditChatInviteLink) WriteJSONTo(w io.Writer) error   { return nil }

func (EditMessageCaption) WriteJSON(j *json.Stream)        {}
func (EditMessageCaption) ReadJSON(i *json.Iterator) error { return nil }
func (EditMessageCaption) ReadJSONFrom(r io.Reader) error  { return nil }
func (EditMessageCaption) WriteJSONTo(w io.Writer) error   { return nil }

func (EditMessageLiveLocation) WriteJSON(j *json.Stream)        {}
func (EditMessageLiveLocation) ReadJSON(i *json.Iterator) error { return nil }
func (EditMessageLiveLocation) ReadJSONFrom(r io.Reader) error  { return nil }
func (EditMessageLiveLocation) WriteJSONTo(w io.Writer) error   { return nil }

func (EditMessageMedia) WriteJSON(j *json.Stream)        {}
func (EditMessageMedia) ReadJSON(i *json.Iterator) error { return nil }
func (EditMessageMedia) ReadJSONFrom(r io.Reader) error  { return nil }
func (EditMessageMedia) WriteJSONTo(w io.Writer) error   { return nil }

func (EditMessageReplyMarkup) WriteJSON(j *json.Stream)        {}
func (EditMessageReplyMarkup) ReadJSON(i *json.Iterator) error { return nil }
func (EditMessageReplyMarkup) ReadJSONFrom(r io.Reader) error  { return nil }
func (EditMessageReplyMarkup) WriteJSONTo(w io.Writer) error   { return nil }

func (EditMessageText) WriteJSON(j *json.Stream)        {}
func (EditMessageText) ReadJSON(i *json.Iterator) error { return nil }
func (EditMessageText) ReadJSONFrom(r io.Reader) error  { return nil }
func (EditMessageText) WriteJSONTo(w io.Writer) error   { return nil }

func (ExportChatInviteLink) WriteJSON(j *json.Stream)        {}
func (ExportChatInviteLink) ReadJSON(i *json.Iterator) error { return nil }
func (ExportChatInviteLink) ReadJSONFrom(r io.Reader) error  { return nil }
func (ExportChatInviteLink) WriteJSONTo(w io.Writer) error   { return nil }

func (ForwardMessage) WriteJSON(j *json.Stream)        {}
func (ForwardMessage) ReadJSON(i *json.Iterator) error { return nil }
func (ForwardMessage) ReadJSONFrom(r io.Reader) error  { return nil }
func (ForwardMessage) WriteJSONTo(w io.Writer) error   { return nil }

func (GetChat) WriteJSON(j *json.Stream)        {}
func (GetChat) ReadJSON(i *json.Iterator) error { return nil }
func (GetChat) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetChat) WriteJSONTo(w io.Writer) error   { return nil }

func (GetChatAdministrators) WriteJSON(j *json.Stream)        {}
func (GetChatAdministrators) ReadJSON(i *json.Iterator) error { return nil }
func (GetChatAdministrators) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetChatAdministrators) WriteJSONTo(w io.Writer) error   { return nil }

func (GetChatMember) WriteJSON(j *json.Stream)        {}
func (GetChatMember) ReadJSON(i *json.Iterator) error { return nil }
func (GetChatMember) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetChatMember) WriteJSONTo(w io.Writer) error   { return nil }

func (GetChatMemberCount) WriteJSON(j *json.Stream)        {}
func (GetChatMemberCount) ReadJSON(i *json.Iterator) error { return nil }
func (GetChatMemberCount) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetChatMemberCount) WriteJSONTo(w io.Writer) error   { return nil }

func (GetFile) WriteJSON(j *json.Stream)        {}
func (GetFile) ReadJSON(i *json.Iterator) error { return nil }
func (GetFile) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetFile) WriteJSONTo(w io.Writer) error   { return nil }

func (GetGameHighScores) WriteJSON(j *json.Stream)        {}
func (GetGameHighScores) ReadJSON(i *json.Iterator) error { return nil }
func (GetGameHighScores) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetGameHighScores) WriteJSONTo(w io.Writer) error   { return nil }

func (GetMyCommands) WriteJSON(j *json.Stream)        {}
func (GetMyCommands) ReadJSON(i *json.Iterator) error { return nil }
func (GetMyCommands) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetMyCommands) WriteJSONTo(w io.Writer) error   { return nil }

func (GetStickerSet) WriteJSON(j *json.Stream)        {}
func (GetStickerSet) ReadJSON(i *json.Iterator) error { return nil }
func (GetStickerSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetStickerSet) WriteJSONTo(w io.Writer) error   { return nil }

func (GetUpdates) WriteJSON(j *json.Stream)        {}
func (GetUpdates) ReadJSON(i *json.Iterator) error { return nil }
func (GetUpdates) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetUpdates) WriteJSONTo(w io.Writer) error   { return nil }

func (GetUserProfilePhotos) WriteJSON(j *json.Stream)        {}
func (GetUserProfilePhotos) ReadJSON(i *json.Iterator) error { return nil }
func (GetUserProfilePhotos) ReadJSONFrom(r io.Reader) error  { return nil }
func (GetUserProfilePhotos) WriteJSONTo(w io.Writer) error   { return nil }

func (InlineKeyboardMarkup) WriteJSON(j *json.Stream)        {}
func (InlineKeyboardMarkup) ReadJSON(i *json.Iterator) error { return nil }
func (InlineKeyboardMarkup) ReadJSONFrom(r io.Reader) error  { return nil }
func (InlineKeyboardMarkup) WriteJSONTo(w io.Writer) error   { return nil }

func (LeaveChat) WriteJSON(j *json.Stream)        {}
func (LeaveChat) ReadJSON(i *json.Iterator) error { return nil }
func (LeaveChat) ReadJSONFrom(r io.Reader) error  { return nil }
func (LeaveChat) WriteJSONTo(w io.Writer) error   { return nil }

// WriteJSON implements json.Marshaler.
func (s MaskPosition) WriteJSON(j *json.Stream) {
	j.WriteObjectStart()
	more := json.NewMore(j)
	defer more.Reset()
	more.More()
	j.WriteObjectField("point")
	j.WriteString(s.Point)
	more.More()
	j.WriteObjectField("scale")
	j.WriteFloat64(s.Scale)
	more.More()
	j.WriteObjectField("x_shift")
	j.WriteFloat64(s.XShift)
	more.More()
	j.WriteObjectField("y_shift")
	j.WriteFloat64(s.YShift)
	j.WriteObjectEnd()
}

// WriteJSONTo writes MaskPosition json value to io.Writer.
func (s MaskPosition) WriteJSONTo(w io.Writer) error {
	j := json.GetStream(w)
	defer json.PutStream(j)
	s.WriteJSON(j)
	return j.Flush()
}

// ReadJSONFrom reads MaskPosition json value from io.Reader.
func (s *MaskPosition) ReadJSONFrom(r io.Reader) error {
	buf := json.GetBuffer()
	defer json.PutBuffer(buf)

	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}
	i := json.GetIterator()
	i.ResetBytes(buf.Bytes())
	defer json.PutIterator(i)

	return s.ReadJSON(i)
}

// ReadJSON reads MaskPosition from json stream.
func (s *MaskPosition) ReadJSON(i *json.Iterator) error {
	i.ReadObjectCB(func(i *json.Iterator, k string) bool {
		switch k {
		case "point":
			s.Point = i.ReadString()
			return i.Error == nil
		case "scale":
			s.Scale = i.ReadFloat64()
			return i.Error == nil
		case "x_shift":
			s.XShift = i.ReadFloat64()
			return i.Error == nil
		case "y_shift":
			s.YShift = i.ReadFloat64()
			return i.Error == nil
		default:
			i.Skip()
			return true
		}
	})
	return i.Error
}

// WriteJSON writes json value of bool to json stream.
func (o OptBool) WriteJSON(j *json.Stream) {
	j.WriteBool(bool(o.Value))
}

// ReadJSON reads json value of bool from json iterator.
func (o *OptBool) ReadJSON(i *json.Iterator) error {
	switch i.WhatIsNext() {
	case json.BoolValue:
		o.Set = true
		o.Value = bool(i.ReadBool())
		return i.Error
	default:
		return fmt.Errorf("unexpected type %d while reading OptBool", i.WhatIsNext())
	}
	return nil
}

// WriteJSON writes json value of string to json stream.
func (o OptString) WriteJSON(j *json.Stream) {
	j.WriteString(string(o.Value))
}

// ReadJSON reads json value of string from json iterator.
func (o *OptString) ReadJSON(i *json.Iterator) error {
	switch i.WhatIsNext() {
	case json.StringValue:
		o.Set = true
		o.Value = string(i.ReadString())
		return i.Error
	default:
		return fmt.Errorf("unexpected type %d while reading OptString", i.WhatIsNext())
	}
	return nil
}

func (PinChatMessage) WriteJSON(j *json.Stream)        {}
func (PinChatMessage) ReadJSON(i *json.Iterator) error { return nil }
func (PinChatMessage) ReadJSONFrom(r io.Reader) error  { return nil }
func (PinChatMessage) WriteJSONTo(w io.Writer) error   { return nil }

func (PromoteChatMember) WriteJSON(j *json.Stream)        {}
func (PromoteChatMember) ReadJSON(i *json.Iterator) error { return nil }
func (PromoteChatMember) ReadJSONFrom(r io.Reader) error  { return nil }
func (PromoteChatMember) WriteJSONTo(w io.Writer) error   { return nil }

func (RestrictChatMember) WriteJSON(j *json.Stream)        {}
func (RestrictChatMember) ReadJSON(i *json.Iterator) error { return nil }
func (RestrictChatMember) ReadJSONFrom(r io.Reader) error  { return nil }
func (RestrictChatMember) WriteJSONTo(w io.Writer) error   { return nil }

func (RevokeChatInviteLink) WriteJSON(j *json.Stream)        {}
func (RevokeChatInviteLink) ReadJSON(i *json.Iterator) error { return nil }
func (RevokeChatInviteLink) ReadJSONFrom(r io.Reader) error  { return nil }
func (RevokeChatInviteLink) WriteJSONTo(w io.Writer) error   { return nil }

func (SendAnimation) WriteJSON(j *json.Stream)        {}
func (SendAnimation) ReadJSON(i *json.Iterator) error { return nil }
func (SendAnimation) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendAnimation) WriteJSONTo(w io.Writer) error   { return nil }

func (SendAudio) WriteJSON(j *json.Stream)        {}
func (SendAudio) ReadJSON(i *json.Iterator) error { return nil }
func (SendAudio) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendAudio) WriteJSONTo(w io.Writer) error   { return nil }

func (SendChatAction) WriteJSON(j *json.Stream)        {}
func (SendChatAction) ReadJSON(i *json.Iterator) error { return nil }
func (SendChatAction) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendChatAction) WriteJSONTo(w io.Writer) error   { return nil }

func (SendContact) WriteJSON(j *json.Stream)        {}
func (SendContact) ReadJSON(i *json.Iterator) error { return nil }
func (SendContact) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendContact) WriteJSONTo(w io.Writer) error   { return nil }

func (SendDice) WriteJSON(j *json.Stream)        {}
func (SendDice) ReadJSON(i *json.Iterator) error { return nil }
func (SendDice) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendDice) WriteJSONTo(w io.Writer) error   { return nil }

func (SendDocument) WriteJSON(j *json.Stream)        {}
func (SendDocument) ReadJSON(i *json.Iterator) error { return nil }
func (SendDocument) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendDocument) WriteJSONTo(w io.Writer) error   { return nil }

func (SendGame) WriteJSON(j *json.Stream)        {}
func (SendGame) ReadJSON(i *json.Iterator) error { return nil }
func (SendGame) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendGame) WriteJSONTo(w io.Writer) error   { return nil }

func (SendInvoice) WriteJSON(j *json.Stream)        {}
func (SendInvoice) ReadJSON(i *json.Iterator) error { return nil }
func (SendInvoice) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendInvoice) WriteJSONTo(w io.Writer) error   { return nil }

func (SendLocation) WriteJSON(j *json.Stream)        {}
func (SendLocation) ReadJSON(i *json.Iterator) error { return nil }
func (SendLocation) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendLocation) WriteJSONTo(w io.Writer) error   { return nil }

func (SendMediaGroup) WriteJSON(j *json.Stream)        {}
func (SendMediaGroup) ReadJSON(i *json.Iterator) error { return nil }
func (SendMediaGroup) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendMediaGroup) WriteJSONTo(w io.Writer) error   { return nil }

func (SendMessage) WriteJSON(j *json.Stream)        {}
func (SendMessage) ReadJSON(i *json.Iterator) error { return nil }
func (SendMessage) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendMessage) WriteJSONTo(w io.Writer) error   { return nil }

func (SendPhoto) WriteJSON(j *json.Stream)        {}
func (SendPhoto) ReadJSON(i *json.Iterator) error { return nil }
func (SendPhoto) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendPhoto) WriteJSONTo(w io.Writer) error   { return nil }

func (SendPoll) WriteJSON(j *json.Stream)        {}
func (SendPoll) ReadJSON(i *json.Iterator) error { return nil }
func (SendPoll) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendPoll) WriteJSONTo(w io.Writer) error   { return nil }

func (SendSticker) WriteJSON(j *json.Stream)        {}
func (SendSticker) ReadJSON(i *json.Iterator) error { return nil }
func (SendSticker) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendSticker) WriteJSONTo(w io.Writer) error   { return nil }

func (SendVenue) WriteJSON(j *json.Stream)        {}
func (SendVenue) ReadJSON(i *json.Iterator) error { return nil }
func (SendVenue) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendVenue) WriteJSONTo(w io.Writer) error   { return nil }

func (SendVideo) WriteJSON(j *json.Stream)        {}
func (SendVideo) ReadJSON(i *json.Iterator) error { return nil }
func (SendVideo) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendVideo) WriteJSONTo(w io.Writer) error   { return nil }

func (SendVideoNote) WriteJSON(j *json.Stream)        {}
func (SendVideoNote) ReadJSON(i *json.Iterator) error { return nil }
func (SendVideoNote) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendVideoNote) WriteJSONTo(w io.Writer) error   { return nil }

func (SendVoice) WriteJSON(j *json.Stream)        {}
func (SendVoice) ReadJSON(i *json.Iterator) error { return nil }
func (SendVoice) ReadJSONFrom(r io.Reader) error  { return nil }
func (SendVoice) WriteJSONTo(w io.Writer) error   { return nil }

func (SetChatAdministratorCustomTitle) WriteJSON(j *json.Stream)        {}
func (SetChatAdministratorCustomTitle) ReadJSON(i *json.Iterator) error { return nil }
func (SetChatAdministratorCustomTitle) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetChatAdministratorCustomTitle) WriteJSONTo(w io.Writer) error   { return nil }

func (SetChatDescription) WriteJSON(j *json.Stream)        {}
func (SetChatDescription) ReadJSON(i *json.Iterator) error { return nil }
func (SetChatDescription) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetChatDescription) WriteJSONTo(w io.Writer) error   { return nil }

func (SetChatPermissions) WriteJSON(j *json.Stream)        {}
func (SetChatPermissions) ReadJSON(i *json.Iterator) error { return nil }
func (SetChatPermissions) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetChatPermissions) WriteJSONTo(w io.Writer) error   { return nil }

func (SetChatPhoto) WriteJSON(j *json.Stream)        {}
func (SetChatPhoto) ReadJSON(i *json.Iterator) error { return nil }
func (SetChatPhoto) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetChatPhoto) WriteJSONTo(w io.Writer) error   { return nil }

func (SetChatStickerSet) WriteJSON(j *json.Stream)        {}
func (SetChatStickerSet) ReadJSON(i *json.Iterator) error { return nil }
func (SetChatStickerSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetChatStickerSet) WriteJSONTo(w io.Writer) error   { return nil }

func (SetChatTitle) WriteJSON(j *json.Stream)        {}
func (SetChatTitle) ReadJSON(i *json.Iterator) error { return nil }
func (SetChatTitle) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetChatTitle) WriteJSONTo(w io.Writer) error   { return nil }

func (SetGameScore) WriteJSON(j *json.Stream)        {}
func (SetGameScore) ReadJSON(i *json.Iterator) error { return nil }
func (SetGameScore) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetGameScore) WriteJSONTo(w io.Writer) error   { return nil }

func (SetMyCommands) WriteJSON(j *json.Stream)        {}
func (SetMyCommands) ReadJSON(i *json.Iterator) error { return nil }
func (SetMyCommands) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetMyCommands) WriteJSONTo(w io.Writer) error   { return nil }

func (SetPassportDataErrors) WriteJSON(j *json.Stream)        {}
func (SetPassportDataErrors) ReadJSON(i *json.Iterator) error { return nil }
func (SetPassportDataErrors) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetPassportDataErrors) WriteJSONTo(w io.Writer) error   { return nil }

func (SetStickerPositionInSet) WriteJSON(j *json.Stream)        {}
func (SetStickerPositionInSet) ReadJSON(i *json.Iterator) error { return nil }
func (SetStickerPositionInSet) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetStickerPositionInSet) WriteJSONTo(w io.Writer) error   { return nil }

func (SetStickerSetThumb) WriteJSON(j *json.Stream)        {}
func (SetStickerSetThumb) ReadJSON(i *json.Iterator) error { return nil }
func (SetStickerSetThumb) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetStickerSetThumb) WriteJSONTo(w io.Writer) error   { return nil }

func (SetWebhook) WriteJSON(j *json.Stream)        {}
func (SetWebhook) ReadJSON(i *json.Iterator) error { return nil }
func (SetWebhook) ReadJSONFrom(r io.Reader) error  { return nil }
func (SetWebhook) WriteJSONTo(w io.Writer) error   { return nil }

func (StopMessageLiveLocation) WriteJSON(j *json.Stream)        {}
func (StopMessageLiveLocation) ReadJSON(i *json.Iterator) error { return nil }
func (StopMessageLiveLocation) ReadJSONFrom(r io.Reader) error  { return nil }
func (StopMessageLiveLocation) WriteJSONTo(w io.Writer) error   { return nil }

func (StopPoll) WriteJSON(j *json.Stream)        {}
func (StopPoll) ReadJSON(i *json.Iterator) error { return nil }
func (StopPoll) ReadJSONFrom(r io.Reader) error  { return nil }
func (StopPoll) WriteJSONTo(w io.Writer) error   { return nil }

func (UnbanChatMember) WriteJSON(j *json.Stream)        {}
func (UnbanChatMember) ReadJSON(i *json.Iterator) error { return nil }
func (UnbanChatMember) ReadJSONFrom(r io.Reader) error  { return nil }
func (UnbanChatMember) WriteJSONTo(w io.Writer) error   { return nil }

func (UnpinAllChatMessages) WriteJSON(j *json.Stream)        {}
func (UnpinAllChatMessages) ReadJSON(i *json.Iterator) error { return nil }
func (UnpinAllChatMessages) ReadJSONFrom(r io.Reader) error  { return nil }
func (UnpinAllChatMessages) WriteJSONTo(w io.Writer) error   { return nil }

func (UnpinChatMessage) WriteJSON(j *json.Stream)        {}
func (UnpinChatMessage) ReadJSON(i *json.Iterator) error { return nil }
func (UnpinChatMessage) ReadJSONFrom(r io.Reader) error  { return nil }
func (UnpinChatMessage) WriteJSONTo(w io.Writer) error   { return nil }

func (UploadStickerFile) WriteJSON(j *json.Stream)        {}
func (UploadStickerFile) ReadJSON(i *json.Iterator) error { return nil }
func (UploadStickerFile) ReadJSONFrom(r io.Reader) error  { return nil }
func (UploadStickerFile) WriteJSONTo(w io.Writer) error   { return nil }

// WriteJSON implements json.Marshaler.
func (s User) WriteJSON(j *json.Stream) {
	j.WriteObjectStart()
	more := json.NewMore(j)
	defer more.Reset()
	if s.CanJoinGroups.Set {
		more.More()
		j.WriteObjectField("can_join_groups")
		s.CanJoinGroups.WriteJSON(j)
	}
	if s.CanReadAllGroupMessages.Set {
		more.More()
		j.WriteObjectField("can_read_all_group_messages")
		s.CanReadAllGroupMessages.WriteJSON(j)
	}
	more.More()
	j.WriteObjectField("first_name")
	j.WriteString(s.FirstName)
	more.More()
	j.WriteObjectField("id")
	j.WriteInt(s.ID)
	more.More()
	j.WriteObjectField("is_bot")
	j.WriteBool(s.IsBot)
	if s.LanguageCode.Set {
		more.More()
		j.WriteObjectField("language_code")
		s.LanguageCode.WriteJSON(j)
	}
	if s.LastName.Set {
		more.More()
		j.WriteObjectField("last_name")
		s.LastName.WriteJSON(j)
	}
	if s.SupportsInlineQueries.Set {
		more.More()
		j.WriteObjectField("supports_inline_queries")
		s.SupportsInlineQueries.WriteJSON(j)
	}
	if s.Username.Set {
		more.More()
		j.WriteObjectField("username")
		s.Username.WriteJSON(j)
	}
	j.WriteObjectEnd()
}

// WriteJSONTo writes User json value to io.Writer.
func (s User) WriteJSONTo(w io.Writer) error {
	j := json.GetStream(w)
	defer json.PutStream(j)
	s.WriteJSON(j)
	return j.Flush()
}

// ReadJSONFrom reads User json value from io.Reader.
func (s *User) ReadJSONFrom(r io.Reader) error {
	buf := json.GetBuffer()
	defer json.PutBuffer(buf)

	if _, err := buf.ReadFrom(r); err != nil {
		return err
	}
	i := json.GetIterator()
	i.ResetBytes(buf.Bytes())
	defer json.PutIterator(i)

	return s.ReadJSON(i)
}

// ReadJSON reads User from json stream.
func (s *User) ReadJSON(i *json.Iterator) error {
	i.ReadObjectCB(func(i *json.Iterator, k string) bool {
		switch k {
		case "can_join_groups":
			s.CanJoinGroups.Reset()
			if err := s.CanJoinGroups.ReadJSON(i); err != nil {
				i.ReportError("Field CanJoinGroups", err.Error())
				return false
			}
			return true
		case "can_read_all_group_messages":
			s.CanReadAllGroupMessages.Reset()
			if err := s.CanReadAllGroupMessages.ReadJSON(i); err != nil {
				i.ReportError("Field CanReadAllGroupMessages", err.Error())
				return false
			}
			return true
		case "first_name":
			s.FirstName = i.ReadString()
			return i.Error == nil
		case "id":
			s.ID = i.ReadInt()
			return i.Error == nil
		case "is_bot":
			s.IsBot = i.ReadBool()
			return i.Error == nil
		case "language_code":
			s.LanguageCode.Reset()
			if err := s.LanguageCode.ReadJSON(i); err != nil {
				i.ReportError("Field LanguageCode", err.Error())
				return false
			}
			return true
		case "last_name":
			s.LastName.Reset()
			if err := s.LastName.ReadJSON(i); err != nil {
				i.ReportError("Field LastName", err.Error())
				return false
			}
			return true
		case "supports_inline_queries":
			s.SupportsInlineQueries.Reset()
			if err := s.SupportsInlineQueries.ReadJSON(i); err != nil {
				i.ReportError("Field SupportsInlineQueries", err.Error())
				return false
			}
			return true
		case "username":
			s.Username.Reset()
			if err := s.Username.ReadJSON(i); err != nil {
				i.ReportError("Field Username", err.Error())
				return false
			}
			return true
		default:
			i.Skip()
			return true
		}
	})
	return i.Error
}
