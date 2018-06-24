package ce

import (
	"html"
	"strings"
)

func TextFromHTML(rawhtml string) string {
	lines := strings.Split(rawhtml, "\n")
	for i, _ := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	rawhtml = strings.Join(lines, "\n")
	for _, v := range ReIgnoreBlock {
		rawhtml = v.ReplaceAllString(rawhtml, "")
	}
	for k, v := range ReNewLineBlock {
		rawhtml = v.ReplaceAllString(rawhtml, "\n"+k)
	}
	rawhtml = ReMultiNewLine.ReplaceAllString(rawhtml, "\n")
	text := strings.TrimSpace(ReTag.ReplaceAllString(rawhtml, ""))
	lines = strings.Split(text, "\n")
	for i, _ := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return html.UnescapeString(strings.Join(lines, "\n"))
}
