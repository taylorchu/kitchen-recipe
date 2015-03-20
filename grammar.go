package main

import (
	"regexp"
	"strings"
)

var (
	Grammar1 = map[string]string{
		"PNP":       "Prep NP",
		"NP":        "(?:Adj|Unit)* Noun+",
		"NP_LIST":   "NP (?:(?:Comma|Conj)+ NP)*",
		"VERB_LIST": "Verb (?:(?:Comma|Conj)+ Verb)*",

		"MODIFIER": "Adv? Prep? (Unit*) (PNP?)",
		"ACTION":   "(Noun?) (VERB_LIST) (Adv*) (Prep? NP_LIST|Adj)? MODIFIER",

		"ACTION_LIST":            "ACTION (?:Comma Conj? (Prep?) ACTION)? (?:Comma Conj? (Prep?) ACTION)?",
		"PNP_ACTION":             "(PNP) Comma ACTION",
		"SIMPLE_ACTION":          "(Conj) ACTION Comma ACTION",
		"ACTION_ACTION_MODIFIER": "ACTION Comma ACTION Comma MODIFIER",

		"SENT": "(?:ACTION_LIST|PNP_ACTION|SIMPLE_ACTION|ACTION_ACTION_MODIFIER) Period",
	}
)

func Expand(grammar map[string]string, s string) string {
	for {
		next := regexp.MustCompile("[A-Z_]{2,}").ReplaceAllStringFunc(s, func(from string) string {
			return "(?:" + grammar[from] + ")"
		})
		if next == s {
			break
		}
		s = next
	}
	return strings.Replace(s, " ", "", -1)
}
