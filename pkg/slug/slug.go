package slug

import (
	"regexp"
	"strings"
)

var (
	re1 = regexp.MustCompile(`(?m)[^a-zA-Z0-9ก-๙]+`)
	re2 = regexp.MustCompile(`(?m)\s\s+`)
	re3 = regexp.MustCompile(`(?m)^(-|\s)`)
)

// Make returns slug generated from provided string.
func Make(s string) string {
	return strings.ReplaceAll(re3.ReplaceAllString(re2.ReplaceAllString(re1.ReplaceAllString(strings.ToLower(s), " "), " "), ""), " ", "-")
}
