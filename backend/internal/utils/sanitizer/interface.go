package sanitizer

import "regexp"

type Filter interface {
	Apply(input string) string
}

type DefaultFilter struct {
	patterns []*regexp.Regexp
}
