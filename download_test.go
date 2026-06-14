package botapi

import (
	"bytes"
	"context"
	"testing"
)

func TestDownloadFileErrors(t *testing.T) {
	b := newMockBot(newMockInvoker())
	if _, err := b.DownloadFile(context.Background(), "not-a-file-id", nil); err == nil {
		t.Fatal("expected error for invalid file_id")
	}
	if err := b.DownloadFileToPath(context.Background(), "not-a-file-id", "/tmp/x"); err == nil {
		t.Fatal("expected error for invalid file_id")
	}
}

func TestCountWriter(t *testing.T) {
	var buf bytes.Buffer
	cw := &countWriter{w: &buf}
	n, err := cw.Write([]byte("hello"))
	if err != nil || n != 5 || cw.n != 5 {
		t.Fatalf("write: n=%d cw.n=%d err=%v", n, cw.n, err)
	}
	_, _ = cw.Write([]byte("!"))
	if cw.n != 6 || buf.String() != "hello!" {
		t.Fatalf("count = %d buf = %q", cw.n, buf.String())
	}
}
