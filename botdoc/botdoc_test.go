package botdoc

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/require"
)

func assertNoUnknown(t testing.TB, def []Definition) {
	t.Helper()
	for _, d := range def {
		for _, f := range d.Fields {
			if f.Type.Kind == "" {
				t.Errorf("invalid type %s/%s: %s", d.Name, f.Name, f.Type.Name)
			}
		}
	}
}

func TestExtract(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("_testdata", "api.html"))
	require.NoError(t, err)

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	require.NoError(t, err)

	a := Extract(doc)
	assertNoUnknown(t, a.Methods)
	assertNoUnknown(t, a.Types)

	for _, dd := range a.Types {
		t.Log(dd.Name, dd.Description)
		for _, f := range dd.Fields {
			t.Logf(" %s %s (%s)", f.Name, f.Type, f.Description)
		}
	}
}
