package botdoc

import (
	"regexp"
	"strconv"
)

type bound struct {
	Min int64
	Max uint64
}

var (
	charBoundRegex = regexp.MustCompile(`(?P<start>\d+)-(?P<end>\d+) (characters|bytes)`)
	intBoundRegex  = regexp.MustCompile(`(?:between|;) (?P<start>\d+)(?:-|\sand\s)(?P<end>\d+)`)
)

func matchBounds(r *regexp.Regexp, matches []string) (a, b int) {
	start, err := strconv.Atoi(matches[r.SubexpIndex("start")])
	if err != nil {
		return a, b
	}
	end, err := strconv.Atoi(matches[r.SubexpIndex("end")])
	if err != nil {
		return a, b
	}
	return start, end
}

func regexBounds(r *regexp.Regexp, s string) (a, b int) {
	matches := r.FindStringSubmatch(s)
	if len(matches) < 3 {
		return a, b
	}
	return matchBounds(r, matches)
}

func stringBounds(s string) bound {
	start, end := regexBounds(charBoundRegex, s)
	return bound{Min: int64(start), Max: uint64(end)}
}

func intBounds(s string) bound {
	start, end := regexBounds(intBoundRegex, s)
	return bound{Min: int64(start), Max: uint64(end)}
}
