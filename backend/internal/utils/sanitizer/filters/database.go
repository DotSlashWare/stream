package filters

import "regexp"

type DatabaseFilter struct {
	patterns []*regexp.Regexp
}

func NewDatabaseFilter() *DatabaseFilter {
	return &DatabaseFilter{
		patterns: []*regexp.Regexp{
			regexp.MustCompile(`postgresql://[^:]+:[^@]+@`),
			regexp.MustCompile(`mongodb://[^:]+:[^@]+@`),
			regexp.MustCompile(`mysql://[^:]+:[^@]+@`),
			regexp.MustCompile(`redis://[^:]*:[^@]*@`),
			regexp.MustCompile(`sqlite://[^?]*\?.*password=[^&\s]+`),
			regexp.MustCompile(`mssql://[^:]+:[^@]+@`),
		},
	}
}

func (f *DatabaseFilter) Apply(input string) string {
	for _, pattern := range f.patterns {
		input = pattern.ReplaceAllString(input, "********@")
	}
	return input
}
