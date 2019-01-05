package ce

import (
	"unicode"
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
