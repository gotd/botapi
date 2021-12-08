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
	charBoundRegex = regexp.MustCompile(`(\d+)-(\d+) characters`)
	intBoundRegex  = regexp.MustCompile(`Values between (\d+)-(\d+) are accepted`)
)

func matchBounds(matches [][]byte) (a, b int) {
	start, err := strconv.Atoi(string(matches[1]))
	if err != nil {
		return a, b
	}
	end, err := strconv.Atoi(string(matches[2]))
	if err != nil {
		return a, b
	}
	return start, end
}

func regexBounds(r *regexp.Regexp, s string) (a, b int) {
	matches := r.FindSubmatch([]byte(s))
	if len(matches) != 3 {
		return a, b
	}
	return matchBounds(matches)
}

func stringBounds(s string) bound {
	start, end := regexBounds(charBoundRegex, s)
	return bound{Min: int64(start), Max: uint64(end)}
}

func intBounds(s string) bound {
	start, end := regexBounds(intBoundRegex, s)
	return bound{Min: int64(start), Max: uint64(end)}
}
