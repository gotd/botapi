package botapi

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/PuerkitoBio/goquery"

	"github.com/gotd/botapi/internal/botdoc"
)

// coveredByOtherMeans lists published methods the library satisfies through a
// different API shape than a same-named method, mapped to where.
var coveredByOtherMeans = map[string]string{
	"getMe": "Bot.Self()",
}

// deferredMethods lists published methods that are planned but not yet
// implemented. Each is an acknowledged gap tracked in docs/roadmap.md; the
// conformance test allows them so it can still catch *unacknowledged* drift.
var deferredMethods = map[string]string{
	"answerPreCheckoutQuery": "payments — deferred until payment updates land",
	"answerShippingQuery":    "payments — deferred until payment updates land",
	"sendInvoice":            "payments — deferred",
	"setPassportDataErrors":  "Telegram Passport — deferred",
}

// notApplicableMethods lists published methods that do not apply to the
// MTProto-native model and are intentionally not implemented.
var notApplicableMethods = map[string]string{
	"getUpdates":     "MTProto-native: updates arrive on the persistent connection (Decision #2)",
	"setWebhook":     "no webhook surface (Decision #2)",
	"deleteWebhook":  "no webhook surface (Decision #2)",
	"getWebhookInfo": "no webhook surface (Decision #2)",
	"logOut":         "HTTP Bot API server lifecycle method; not applicable over MTProto",
	"close":          "HTTP Bot API server lifecycle method; not applicable over MTProto",
}

// TestMethodConformance asserts that every method published in the Bot API docs
// is either implemented on *Bot, covered by other means, or explicitly
// acknowledged as deferred / not-applicable. It is a drift guard: when Telegram
// ships a new method, this test fails until the method is implemented or
// categorized. It also fails if an allowlist entry stops being a published
// method (so the lists cannot rot).
func TestMethodConformance(t *testing.T) {
	api := loadAPI(t)

	implemented := botMethodSet()
	published := map[string]struct{}{}

	var uncategorized []string
	for _, def := range api.Methods {
		name := def.Name
		published[name] = struct{}{}

		switch {
		case implemented[goName(name)]:
		case coveredByOtherMeans[name] != "":
		case deferredMethods[name] != "":
		case notApplicableMethods[name] != "":
		default:
			uncategorized = append(uncategorized, name)
		}
	}

	if len(uncategorized) > 0 {
		t.Errorf("API drift: %d published method(s) are neither implemented nor categorized: %v\n"+
			"Implement them, or add to deferredMethods/notApplicableMethods in conformance_test.go.",
			len(uncategorized), uncategorized)
	}

	// Guard against stale allowlists: every acknowledged name must still be a
	// published method.
	for _, lists := range []map[string]string{coveredByOtherMeans, deferredMethods, notApplicableMethods} {
		for name := range lists {
			if _, ok := published[name]; !ok {
				t.Errorf("stale allowlist entry %q is no longer a published method; remove it", name)
			}
		}
	}
}

// loadAPI extracts the Bot API surface from the committed documentation
// snapshot.
func loadAPI(t *testing.T) botdoc.API {
	t.Helper()
	path := filepath.Join("internal", "botdoc", "_testdata", "api.html")
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open api docs: %v", err)
	}
	defer func() { _ = f.Close() }()
	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatalf("parse api docs: %v", err)
	}
	return botdoc.Extract(doc)
}

// botMethodSet returns the set of exported method names on *Bot.
func botMethodSet() map[string]bool {
	t := reflect.TypeFor[*Bot]()
	set := make(map[string]bool, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		set[t.Method(i).Name] = true
	}
	return set
}

// goName maps a Bot API method name (lowerCamelCase) to the Go method name that
// would implement it (UpperCamelCase): only the first letter changes.
func goName(apiName string) string {
	if apiName == "" {
		return ""
	}
	b := []byte(apiName)
	if b[0] >= 'a' && b[0] <= 'z' {
		b[0] -= 'a' - 'A'
	}
	return string(b)
}
