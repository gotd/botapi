package botapi

import (
	"encoding/base64"
	"testing"

	"github.com/gotd/td/tg"
)

func TestPassportErrorToTg(t *testing.T) {
	hash := base64.StdEncoding.EncodeToString([]byte("abc"))

	df := &PassportElementErrorDataField{Type: "passport", FieldName: "first_name", DataHash: hash, Message: "blurry"}
	got, err := df.toTg()
	if err != nil {
		t.Fatal(err)
	}
	data, ok := got.(*tg.SecureValueErrorData)
	if !ok || data.Field != "first_name" || data.Text != "blurry" || string(data.DataHash) != "abc" {
		t.Fatalf("data field: %#v", got)
	}
	if _, ok := data.Type.(*tg.SecureValueTypePassport); !ok {
		t.Fatalf("type: %T", data.Type)
	}

	files := &PassportElementErrorFiles{Type: "utility_bill", FileHashes: []string{hash, hash}, Message: "x"}
	fg, err := files.toTg()
	if err != nil {
		t.Fatal(err)
	}
	if f, ok := fg.(*tg.SecureValueErrorFiles); !ok || len(f.FileHash) != 2 {
		t.Fatalf("files: %#v", fg)
	}
}

func TestPassportErrorInvalidHash(t *testing.T) {
	df := &PassportElementErrorFrontSide{Type: "driver_license", FileHash: "!!!not-base64!!!", Message: "x"}
	if _, err := df.toTg(); err == nil {
		t.Fatal("expected error for invalid base64 hash")
	}
}

func TestSecureValueTypeMapping(t *testing.T) {
	cases := map[string]any{
		"driver_license":   &tg.SecureValueTypeDriverLicense{},
		"email":            &tg.SecureValueTypeEmail{},
		"phone_number":     &tg.SecureValueTypePhone{},
		"unknown_whatever": &tg.SecureValueTypePersonalDetails{}, // fallback
	}
	for in, want := range cases {
		got := secureValueType(in)
		if got.TypeID() != want.(tg.SecureValueTypeClass).TypeID() {
			t.Fatalf("%q: got %T, want %T", in, got, want)
		}
	}
}
