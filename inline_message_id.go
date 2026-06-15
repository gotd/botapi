package botapi

import (
	"encoding/base64"
	"strings"

	"github.com/gotd/td/bin"
	"github.com/gotd/td/tg"
)

// encodeInlineMessageID serializes an MTProto inline message id into the Bot API
// inline_message_id string. It is the base64url (unpadded) TL serialization of
// the inputBotInlineMessageID constructor, matching the official Bot API server
// (which delegates to TDLib's base64url_encode(serialize(...))).
func encodeInlineMessageID(id tg.InputBotInlineMessageIDClass) (string, error) {
	if id == nil {
		return "", nil
	}

	var buf bin.Buffer

	if err := id.Encode(&buf); err != nil {
		return "", &Error{Code: 500, Description: "Internal Server Error: " + err.Error()}
	}

	return base64.RawURLEncoding.EncodeToString(buf.Buf), nil
}

// decodeInlineMessageID parses a Bot API inline_message_id string back into the
// MTProto inline message id. Padding, if present, is tolerated.
func decodeInlineMessageID(s string) (tg.InputBotInlineMessageIDClass, error) {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimRight(s, "="))
	if err != nil {
		return nil, errInvalidInlineMessageID()
	}

	id, err := tg.DecodeInputBotInlineMessageID(&bin.Buffer{Buf: raw})
	if err != nil {
		return nil, errInvalidInlineMessageID()
	}

	return id, nil
}

// errInvalidInlineMessageID is returned for a malformed inline_message_id.
func errInvalidInlineMessageID() error {
	return &Error{Code: 400, Description: "Bad Request: invalid inline_message_id"}
}
