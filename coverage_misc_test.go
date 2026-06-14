package botapi

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/gotd/td/fileid"
	"github.com/gotd/td/tg"
)

// TestSetPassportDataErrorsAllVariants drives SetPassportDataErrors with one of
// every PassportElementError variant, covering each variant's toTg plus the
// nil-skip path of the dispatch loop.
func TestSetPassportDataErrorsAllVariants(t *testing.T) {
	inv := newMockInvoker()
	inv.reply(tg.UsersSetSecureValueErrorsRequestTypeID, &tg.BoolTrue{})
	b := newMockBot(inv)

	h := base64.StdEncoding.EncodeToString([]byte("abc"))
	errs := []PassportElementError{
		nil, // skipped
		&PassportElementErrorDataField{Type: "personal_details", FieldName: "f", DataHash: h, Message: "m"},
		&PassportElementErrorFrontSide{Type: "passport", FileHash: h, Message: "m"},
		&PassportElementErrorReverseSide{Type: "driver_license", FileHash: h, Message: "m"},
		&PassportElementErrorSelfie{Type: "identity_card", FileHash: h, Message: "m"},
		&PassportElementErrorFile{Type: "utility_bill", FileHash: h, Message: "m"},
		&PassportElementErrorFiles{Type: "bank_statement", FileHashes: []string{h, h}, Message: "m"},
		&PassportElementErrorTranslationFile{Type: "rental_agreement", FileHash: h, Message: "m"},
		&PassportElementErrorTranslationFiles{Type: "passport_registration", FileHashes: []string{h}, Message: "m"},
		&PassportElementErrorUnspecified{Type: "email", ElementHash: h, Message: "m"},
	}
	if err := b.SetPassportDataErrors(context.Background(), 42, errs); err != nil {
		t.Fatalf("SetPassportDataErrors: %v", err)
	}
	if !inv.called(tg.UsersSetSecureValueErrorsRequestTypeID) {
		t.Fatal("secure value errors RPC not invoked")
	}
}

// TestPassportErrorBadHashPropagates ensures a bad hash in any variant aborts
// the whole call.
func TestPassportErrorBadHashPropagates(t *testing.T) {
	b := newMockBot(newMockInvoker())
	errs := []PassportElementError{
		&PassportElementErrorFiles{Type: "passport", FileHashes: []string{"!!!"}, Message: "m"},
	}
	if err := b.SetPassportDataErrors(context.Background(), 42, errs); err == nil {
		t.Fatal("expected error for invalid hash")
	}
}

// TestFileUniqueID covers every branch of the file_unique_id derivation.
func TestFileUniqueID(t *testing.T) {
	cases := []struct {
		name string
		f    fileid.FileID
	}{
		{"web", fileid.FileID{URL: "https://example.com/x"}},
		{"photo-volume", fileid.FileID{Type: fileid.Photo, PhotoSizeSource: fileid.PhotoSizeSource{VolumeID: 9, LocalID: 3}}},
		{"photo-novolume", fileid.FileID{Type: fileid.Photo, ID: 7}},
		{"document", fileid.FileID{Type: fileid.Document, ID: 5}},
		{"secure", fileid.FileID{Type: fileid.Secure, ID: 1}},
		{"encrypted", fileid.FileID{Type: fileid.Encrypted, ID: 2}},
		{"temp", fileid.FileID{Type: fileid.Temp, ID: 3}},
		{"profilephoto", fileid.FileID{Type: fileid.ProfilePhoto, ID: 4}},
	}
	for _, c := range cases {
		if got := fileUniqueID(c.f); got == "" {
			t.Errorf("%s: empty file_unique_id", c.name)
		}
	}
}

func TestGetFile(t *testing.T) {
	b := &Bot{}
	id := documentFileID(t, 0x1234)
	f, err := b.GetFile(context.Background(), id)
	if err != nil {
		t.Fatal(err)
	}
	if f.FileID != id || f.FileUniqueID == "" {
		t.Fatalf("file = %#v", f)
	}
	if _, err := b.GetFile(context.Background(), "not-a-file-id"); err == nil {
		t.Fatal("expected error for bad file_id")
	}
}

func TestDownloadFileBadID(t *testing.T) {
	b := &Bot{}
	if _, err := b.DownloadFile(context.Background(), "bad", nil); err == nil {
		t.Fatal("expected error for bad file_id (DownloadFile)")
	}
	if err := b.DownloadFileToPath(context.Background(), "bad", "/tmp/x"); err == nil {
		t.Fatal("expected error for bad file_id (DownloadFileToPath)")
	}
}
