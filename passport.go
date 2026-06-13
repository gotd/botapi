package botapi

import (
	"context"
	"encoding/base64"

	"github.com/gotd/td/tg"
)

// PassportElementError is a sealed union describing one error in a Telegram
// Passport element. The user will not be able to resubmit the element until the
// error is resolved.
//
// Concrete variants mirror the Bot API: *PassportElementErrorDataField,
// *PassportElementErrorFrontSide, *PassportElementErrorReverseSide,
// *PassportElementErrorSelfie, *PassportElementErrorFile,
// *PassportElementErrorFiles, *PassportElementErrorTranslationFile,
// *PassportElementErrorTranslationFiles, *PassportElementErrorUnspecified.
type PassportElementError interface {
	isPassportElementError()
	toTg() (tg.SecureValueErrorClass, error)
}

// PassportElementErrorDataField is an error in a data field.
type PassportElementErrorDataField struct {
	Type      string `json:"type"`
	FieldName string `json:"field_name"`
	DataHash  string `json:"data_hash"`
	Message   string `json:"message"`
}

// PassportElementErrorFrontSide is an error in the document's front side.
type PassportElementErrorFrontSide struct {
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorReverseSide is an error in the document's reverse side.
type PassportElementErrorReverseSide struct {
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorSelfie is an error in the selfie with the document.
type PassportElementErrorSelfie struct {
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorFile is an error in a document scan.
type PassportElementErrorFile struct {
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorFiles is an error in a list of document scans.
type PassportElementErrorFiles struct {
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

// PassportElementErrorTranslationFile is an error in a translation scan.
type PassportElementErrorTranslationFile struct {
	Type     string `json:"type"`
	FileHash string `json:"file_hash"`
	Message  string `json:"message"`
}

// PassportElementErrorTranslationFiles is an error in a list of translation
// scans.
type PassportElementErrorTranslationFiles struct {
	Type       string   `json:"type"`
	FileHashes []string `json:"file_hashes"`
	Message    string   `json:"message"`
}

// PassportElementErrorUnspecified is an error in an unspecified place.
type PassportElementErrorUnspecified struct {
	Type        string `json:"type"`
	ElementHash string `json:"element_hash"`
	Message     string `json:"message"`
}

func (*PassportElementErrorDataField) isPassportElementError()        {}
func (*PassportElementErrorFrontSide) isPassportElementError()        {}
func (*PassportElementErrorReverseSide) isPassportElementError()      {}
func (*PassportElementErrorSelfie) isPassportElementError()           {}
func (*PassportElementErrorFile) isPassportElementError()             {}
func (*PassportElementErrorFiles) isPassportElementError()            {}
func (*PassportElementErrorTranslationFile) isPassportElementError()  {}
func (*PassportElementErrorTranslationFiles) isPassportElementError() {}
func (*PassportElementErrorUnspecified) isPassportElementError()      {}

func (e *PassportElementErrorDataField) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.DataHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorData{Type: secureValueType(e.Type), DataHash: h, Field: e.FieldName, Text: e.Message}, nil
}

func (e *PassportElementErrorFrontSide) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.FileHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorFrontSide{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorReverseSide) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.FileHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorReverseSide{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorSelfie) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.FileHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorSelfie{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorFile) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.FileHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorFile{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorFiles) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHashes(e.FileHashes)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorFiles{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorTranslationFile) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.FileHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorTranslationFile{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorTranslationFiles) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHashes(e.FileHashes)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueErrorTranslationFiles{Type: secureValueType(e.Type), FileHash: h, Text: e.Message}, nil
}

func (e *PassportElementErrorUnspecified) toTg() (tg.SecureValueErrorClass, error) {
	h, err := decodeHash(e.ElementHash)
	if err != nil {
		return nil, err
	}
	return &tg.SecureValueError{Type: secureValueType(e.Type), Hash: h, Text: e.Message}, nil
}

// SetPassportDataErrors reports errors in the Telegram Passport data the user
// submitted to the bot, so the user can fix and resubmit them.
func (b *Bot) SetPassportDataErrors(ctx context.Context, userID int64, errs []PassportElementError) error {
	user, err := b.resolveInputUser(ctx, userID)
	if err != nil {
		return err
	}

	tgErrors := make([]tg.SecureValueErrorClass, 0, len(errs))
	for _, e := range errs {
		if e == nil {
			continue
		}
		converted, err := e.toTg()
		if err != nil {
			return err
		}
		tgErrors = append(tgErrors, converted)
	}

	if _, err := b.raw.UsersSetSecureValueErrors(ctx, &tg.UsersSetSecureValueErrorsRequest{
		ID:     user,
		Errors: tgErrors,
	}); err != nil {
		return asAPIError(err)
	}
	return nil
}

// decodeHash decodes a base64 Passport data hash into raw bytes.
func decodeHash(s string) ([]byte, error) {
	h, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, &Error{Code: 400, Description: "Bad Request: invalid passport data hash"}
	}
	return h, nil
}

func decodeHashes(hashes []string) ([][]byte, error) {
	out := make([][]byte, 0, len(hashes))
	for _, s := range hashes {
		h, err := decodeHash(s)
		if err != nil {
			return nil, err
		}
		out = append(out, h)
	}
	return out, nil
}

// secureValueType maps a Bot API Passport element type to the MTProto secure
// value type. An unknown type falls back to personal details.
func secureValueType(t string) tg.SecureValueTypeClass {
	switch t {
	case "personal_details":
		return &tg.SecureValueTypePersonalDetails{}
	case "passport":
		return &tg.SecureValueTypePassport{}
	case "driver_license":
		return &tg.SecureValueTypeDriverLicense{}
	case "identity_card":
		return &tg.SecureValueTypeIdentityCard{}
	case "internal_passport":
		return &tg.SecureValueTypeInternalPassport{}
	case "address":
		return &tg.SecureValueTypeAddress{}
	case "utility_bill":
		return &tg.SecureValueTypeUtilityBill{}
	case "bank_statement":
		return &tg.SecureValueTypeBankStatement{}
	case "rental_agreement":
		return &tg.SecureValueTypeRentalAgreement{}
	case "passport_registration":
		return &tg.SecureValueTypePassportRegistration{}
	case "temporary_registration":
		return &tg.SecureValueTypeTemporaryRegistration{}
	case "phone_number":
		return &tg.SecureValueTypePhone{}
	case "email":
		return &tg.SecureValueTypeEmail{}
	default:
		return &tg.SecureValueTypePersonalDetails{}
	}
}
