package nl

import (
	"errors"
	"regexp"
	"strings"
)

var (
	ErrToken = errors.New("invalid token")
)

func Tokenize(s string) ([]string, error) {
	return tokenize([]string{
		`[A-Za-z0-9'/-]+`,
		`,`,
		`.`,
	}, s)
}

func tokenize(patterns []string, s string) (tokens []string, err error) {
	for {
		var found bool
		s = strings.TrimLeft(s, " \n\t\r")
		for _, pat := range patterns {
			result := regexp.MustCompile("^" + pat).FindString(s)
			if result != "" {
				found = true
				s = strings.TrimPrefix(s, result)
				tokens = append(tokens, result)
				break
			}
		}
		if len(s) == 0 {
			break
		}
		if !found {
			err = ErrToken
			return
		}
	}
	return
}
