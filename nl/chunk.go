package nl

import (
	"errors"
	"regexp"
	"strings"
	"unicode"
)

var (
	ErrChunk = errors.New("invalid chunk pattern")
)

func allUpper(s string) bool {
	for _, r := range s {
		if unicode.IsLower(r) {
			return false
		}
	}
	return true
}

func Chunk(pat, s string) (chunks [][]string, err error) {
	tokens, err := Tokenize(s)
	if err != nil {
		return
	}

	// generate all possible assignments
	var tagged []string
	for _, token := range tokens {
		token = strings.Replace(token, "-", "_", -1)
		if !allUpper(token) {
			token = Singular(strings.ToLower(token))
		}
		pos := TagPOS(token)
		if len(pos) == 0 {
			continue
		}
		var next []string
		for _, p := range pos {
			attached := "(" + token + "/" + p + ")"
			if len(tagged) == 0 {
				next = append(next, attached)
			} else {
				for _, t := range tagged {
					next = append(next, t+attached)
				}
			}
		}
		tagged = next
	}

	// modify pattern to match original word
	pat = regexp.MustCompile("[A-Z][a-z]*").ReplaceAllStringFunc(pat, func(from string) string {
		return `(?:\([^)]+/` + from + `\))`
	})

	// try all assignments on pattern
	for _, t := range tagged {
		match := regexp.MustCompile("^" + pat + "$").FindStringSubmatch(t)
		if match == nil {
			continue
		}
		for _, submatch := range match[1:] {
			if submatch == "" {
				// skip empty submatch
				continue
			}
			// further chunk down a submatch for Comma, Conj, Peroid
			var chunk []string
			for _, token := range regexp.MustCompile(`\([^)]+\)`).FindAllString(submatch, -1) {
				token = strings.Trim(token, "()")
				switch strings.SplitN(token, "/", 2)[1] {
				case Conj:
					fallthrough
				case Comma:
					fallthrough
				case Period:
					if chunk != nil {
						chunks = append(chunks, chunk)
						chunk = nil
					}
				default:
					chunk = append(chunk, token)
				}
			}
			if chunk != nil {
				chunks = append(chunks, chunk)
			}
		}
		return
	}
	err = ErrChunk
	return
}
