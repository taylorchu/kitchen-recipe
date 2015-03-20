package nl

import "strings"

var (
	pluralSuffix = []string{
		"sses",
		"zzes",
		"ses",
		"zes",
		"xes",
		"shes",
		"ches",
		"ies",
		"oes",
		"s",
	}
	singularSuffix = []string{
		"s",
		"z",
		"s",
		"z",
		"x",
		"sh",
		"ch",
		"y",
		"o",
		"",
	}
)

func Singular(s string) string {
	for i, suffix := range pluralSuffix {
		if strings.HasSuffix(s, suffix) {
			return strings.TrimSuffix(s, suffix) + singularSuffix[i]
		}
	}
	return s
}
