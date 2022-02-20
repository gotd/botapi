package botdoc

import (
	"regexp"
	"strconv"
)

type lengthBound struct {
	Min uint64
	Max uint64
}

type intBound struct {
	Min int64
	Max int64
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

func stringBounds(s string) lengthBound {
	start, end := regexBounds(charBoundRegex, s)
	return lengthBound{Min: uint64(start), Max: uint64(end)}
}

func intBounds(s string) intBound {
	start, end := regexBounds(intBoundRegex, s)
	return intBound{Min: int64(start), Max: int64(end)}
}
