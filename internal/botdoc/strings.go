package botdoc

import "strings"

// IndexOneOf is like strings.Index, but for multiple substrings.
//
// It returns index like strings.Index and found substring.
func IndexOneOf(s string, subs ...string) (idx int, variant string) {
	for _, substr := range subs {
		idx := strings.Index(s, substr)
		if idx >= 0 {
			return idx, substr
		}
	}

	return -1, ""
}
