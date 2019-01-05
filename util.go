package ce

import (
	"regexp"
	"time"
	"unicode"

	"github.com/tkuchiki/parsetime"
)

func FilterControlChar(in string) string {
	var ret []rune
	for _, r := range []rune(in) {
		if unicode.IsControl(r) && r != '\n' {
			continue
		}
		ret = append(ret, r)
	}
	return string(ret)
}

var dateth = regexp.MustCompile(`\d+(?:st|nd|rd|th)`)

func ParseTime(tz string, s string) time.Time {
	timeParser, _ := parsetime.NewParseTime(tz)
	dst := dateth.ReplaceAllStringFunc(s, func(in string) string {
		return in[:len(in)-2]
	})
	t, _ := timeParser.Parse(dst)
	return t
}
