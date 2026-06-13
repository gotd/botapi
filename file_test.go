package botapi

import (
	"encoding/base64"
	"encoding/binary"
	"testing"

	"github.com/gotd/td/fileid"
)

func TestFileUniqueID_Document(t *testing.T) {
	f := fileid.FileID{Type: fileid.Document, ID: 0x1122334455667788}
	got := fileUniqueID(f)

	// Expected: little-endian [uint32 type=2][int64 id], RLE-zero encoded, base64url.
	raw := make([]byte, 12)
	binary.LittleEndian.PutUint32(raw[0:], uniqueTypeDocument)
	binary.LittleEndian.PutUint64(raw[4:], 0x1122334455667788)
	want := base64.RawURLEncoding.EncodeToString(rleEncode(raw))

	if got != want {
		t.Fatalf("got %q, want %q", got, want)
	}
	if got == "" || got == "todo" {
		t.Fatalf("unique id not derived: %q", got)
	}
}

func TestFileUniqueID_Web(t *testing.T) {
	f := fileid.FileID{URL: "https://example.com/a.jpg"}
	got := fileUniqueID(f)
	if got == "" {
		t.Fatal("empty unique id for web file")
	}
	// Different URLs must yield different unique ids.
	other := fileUniqueID(fileid.FileID{URL: "https://example.com/b.jpg"})
	if got == other {
		t.Fatal("web unique ids collide across URLs")
	}
}

func TestFileUniqueID_Stable(t *testing.T) {
	f := fileid.FileID{Type: fileid.Video, ID: 42}
	if a, b := fileUniqueID(f), fileUniqueID(f); a != b {
		t.Fatalf("unstable unique id: %q vs %q", a, b)
	}
}

func TestGetFileInvalid(t *testing.T) {
	b := &Bot{}
	if _, err := b.GetFile(t.Context(), "not-a-valid-file-id"); err == nil {
		t.Fatal("expected error for invalid file_id")
	}
}
