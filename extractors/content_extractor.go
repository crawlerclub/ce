package extractors

import (
	"crypto/md5"
	"encoding/hex"
	"html"
	"regexp"
	"strings"
)

const (
	BlocksWidth = 3
	Threshold   = 600 // in bytes
)

var (
	ReIgnoreBlock = map[string]*regexp.Regexp{
		"doctype": regexp.MustCompile(`(?ims)<!DOCTYPE.*?>`),           // raw doctype
		"comment": regexp.MustCompile(`(?ims)<!--.*?-->`),              // raw comment
		"script":  regexp.MustCompile(`(?ims)<script.*?>.*?</script>`), // javascript
		"style":   regexp.MustCompile(`(?ims)<style.*?>.*?</style>`),   // css
		"link":    regexp.MustCompile(`(?ims)<link.*?>`),               // css
	}
	ReNewLineBlock = map[string]*regexp.Regexp{
		"<div>": regexp.MustCompile(`(?is)<div.*?>`),
		"<p>":   regexp.MustCompile(`(?is)<p.*?>`),
		"<br>":  regexp.MustCompile(`(?is)<br.*?>`),
		"<hr>":  regexp.MustCompile(`(?is)<hr.*?>`),
		"<li>":  regexp.MustCompile(`(?is)<li.*?>`),
	}
	ReMultiNewLine = regexp.MustCompile(`(?m)\n+`)
	ReSpaces       = regexp.MustCompile(`(?m)\s+`)
	ReTag          = regexp.MustCompile(`(?ims)<.*?>`)
	ReImg          = regexp.MustCompile(`(?ims)<img.*?>`)
	ReImgSrc       = regexp.MustCompile(`(?ims)<img.+?src=('|")(.+?)('|").*?>`)
	ReTitle        = regexp.MustCompile(`(?ims)<title.*?>(.+?)</title>`)
	ReH            = regexp.MustCompile(`(?ims)<h\d+.*?>(.*?)</h\d+>`)
	ReHead         = regexp.MustCompile(`(?ims)<head.*?>(.*?)<\/head>`)

	MonthStr   = `(?:(?:jan|feb|mar|apr|may|jun|aug|sep|oct|nov|dec)[a-z]*)`
	ReDateTime = regexp.MustCompile(`(?is)((?:` + MonthStr + `[\.,\-\s]*\d{1,2}(?:st|nd|rd|th)*[\.,\-\s]*(\d{4}))|` +
		`(?:\d{1,2}(?:st|nd|rd|th)*[\.,\-\s]*` + MonthStr + `[\.,\-\s]*(\d{4}))|` +
		`(?:(\d{4}-)\d{1,2}-\d{1,2})|` +
		`(?:(\d{1,2}-)\d{1,2}-\d{4})|` +
		`(?:(\d{4}年){0,1}\d{1,2}月\d{1,2}日))`)

	ReTime = regexp.MustCompile(`(?is)((?:0?|[12])\d\s*:+\s*[0-5]\d(?:\s*:+\s*[0-5]\d)?(?:\s*[,:.]*\s*(?:am|pm))?|` +
		`(?:0?|[12])\d\s*[.\s]+\s*[0-5]\d(?:\s*[,:.]*\s*(?:am|pm))+)`)

	ReContinuousA = regexp.MustCompile(`(?is)</a><a`)

	NavSpliters = []string{`|`, `┊`, `-`}
)

func MD5(text string) string {
	h := md5.New()
	h.Write([]byte(text))
	return hex.EncodeToString(h.Sum(nil))
}

func isUsefulLine(line string) bool {
	for _, sep := range NavSpliters {
		if len(strings.Split(line, sep)) >= 5 {
			return false
		}
	}
	return true
}

func getTitle(raw string) string {
	title := ""
	ret := ReTitle.FindAllStringSubmatch(raw, -1)
	if len(ret) > 0 {
		title = ret[0][1]
	}
	return strings.TrimSpace(title)
}

func clean(raw string) string {
	lines := strings.Split(raw, "\n")
	for i, _ := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	raw = strings.Join(lines, "\n")
	//raw = ReContinuousA.ReplaceAllString(raw, "</a> <a")
	for _, v := range ReIgnoreBlock {
		raw = v.ReplaceAllString(raw, "")
	}
	for k, v := range ReNewLineBlock {
		raw = v.ReplaceAllString(raw, "\n"+k)
	}
	raw = ReMultiNewLine.ReplaceAllString(raw, "\n")
	return raw
}

func getMain(text string) string {
	lines := strings.Split(text, "\n")
	var indexDist []int
	size := len(lines)
	for i := 0; i < size-BlocksWidth+1; i++ {
		bytesNum := 0
		for j := i; j < i+BlocksWidth; j++ {
			noSpaces := ReSpaces.ReplaceAllString(lines[j], "")
			bytesNum += len(noSpaces)
		}
		indexDist = append(indexDist, bytesNum)
	}
	main := ""
	start := -1
	end := -1
	startFlag := false
	endFlag := false
	firstMatch := true
	for i := 0; i < len(indexDist)-1; i++ {
		if firstMatch && !startFlag {
			if indexDist[i] > Threshold/2 {
				for j := i + 1; j <= i+2 && j < len(indexDist); j++ {
					if indexDist[j] != 0 {
						firstMatch = false
						startFlag = true
						start = i
						continue
					}
				}
			}
		}
		if indexDist[i] > Threshold && !startFlag {
			for j := i + 1; j <= i+3 && j < len(indexDist); j++ {
				if indexDist[j] != 0 {
					startFlag = true
					start = i
					continue
				}
			}
		}
		if startFlag {
			if indexDist[i] == 0 || indexDist[i+1] == 0 {
				endFlag = true
				end = i
			}
		}
		if endFlag {
			tmp := ""
			for j := start; j <= end; j++ {
				tmp += lines[j] + "\n"
			}
			main += tmp
			startFlag = false
			endFlag = false
		}
	}
	return strings.TrimSpace(main)
}

func Parse(url, raw string) (string, string) {
	raw = clean(raw)
	title := getTitle(raw)

	images := make(map[string]string)
	ret := ReImg.FindAllStringSubmatch(raw, -1)
	for _, r := range ret {
		if len(r) <= 0 {
			continue
		}
		md := MD5(r[0])
		images[md] = r[0]
		raw = strings.Replace(raw, r[0], md, -1)
	}
	text := html.UnescapeString(ReTag.ReplaceAllString(raw, ""))
	content := getMain(text)
	for k, v := range images {
		content = strings.Replace(content, k, v, -1)
	}
	return title, content
}
