package main

import (
	"path/filepath"
	"testing"

	"github.com/gotd/botapi/internal/botdoc"
)

const testdata = "../../internal/botdoc/_testdata/api.html"

func TestLoadAndExtract(t *testing.T) {
	doc, err := loadDocument("", filepath.FromSlash(testdata))
	if err != nil {
		t.Fatal(err)
	}
	api := botdoc.Extract(doc)
	if api.Version == "" {
		t.Fatal("expected a version")
	}
	if len(api.Types) == 0 || len(api.Methods) == 0 {
		t.Fatalf("expected types and methods, got %d/%d", len(api.Types), len(api.Methods))
	}
}

func TestFindDefinition(t *testing.T) {
	doc, err := loadDocument("", filepath.FromSlash(testdata))
	if err != nil {
		t.Fatal(err)
	}
	api := botdoc.Extract(doc)

	// Case-insensitive match against a method.
	if _, ok := findDefinition(api.Methods, "sendmessage"); !ok {
		t.Fatal("sendMessage method should be found case-insensitively")
	}
	// And a type.
	if _, ok := findDefinition(api.Types, "Message"); !ok {
		t.Fatal("Message type should be found")
	}
	if _, ok := findDefinition(api.Types, "NoSuchThing"); ok {
		t.Fatal("unknown name should not be found")
	}
}
