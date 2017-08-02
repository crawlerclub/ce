package ce

import (
	"regexp"
	"strings"
	"unicode"
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
	ReImgSrc       = regexp.MustCompile(`(?ims)<img.+?src=\s*?"(.+?)"|'(.+?)'.*?>`)
	ReTitle        = regexp.MustCompile(`(?ims)<title.*?>(.+?)</title>`)
	ReH            = regexp.MustCompile(`(?ims)<h\d+.*?>(.*?)</h\d+>`)
	ReHead         = regexp.MustCompile(`(?ims)<head.*?>(.*?)<\/head>`)

	MonthStr = `(?:(?:jan|feb|mar|apr|may|jun|jul|aug|sep|oct|nov|dec)[a-z]*)`
	ReDate   = regexp.MustCompile(`(?is)((?:` + MonthStr + `[\.,\-\s]*\d{1,2}(?:st|nd|rd|th)*[\.,\-\s]*(\d{4}))|` +
		`(?:\d{1,2}(?:st|nd|rd|th)*[\.,\-\s]*` + MonthStr + `[\.,\-\s]*(\d{4}))|` +
		MonthStr + `.\d{1,2}|` +
		`(?:(19|20)\d{2}[^0-9]\d{1,2}[^0-9]\d{1,2})|` +
		`(?:\d{1,2}[^0-9]\d{1,2}[^0-9](19|20)\d{2})|` +
		`(?:(\d{4}年){0,1}\d{1,2}月\d{1,2}日))`)

	ReTime = regexp.MustCompile(`(?is)((?:0?|[12])\d\s*:+\s*[0-5]\d(?:\s*:+\s*[0-5]\d)?(?:\s*[,:.]*\s*(?:am|pm))?|` +
		`(?:0?|[12])\d\s*[.\s]+\s*[0-5]\d(?:\s*[,:.]*\s*(?:am|pm))+)`)

	ReFavicon = regexp.MustCompile(`(?ims)<link rel="shortcut icon" href="(.+?)"/>`)
)

func clean(rawhtml string) string {
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
	return rawhtml
}

func getFavicon(rawhtml string) string {
	ret := ReFavicon.FindAllStringSubmatch(rawhtml, -1)
	if len(ret) > 0 {
		return ret[0][1]
	}
	return ""
}

func getTitle(rawhtml string) string {
	title := ""
	ret := ReTitle.FindAllStringSubmatch(rawhtml, -1)
	if len(ret) > 0 {
		title = ret[0][1]
	}
	h := ReH.FindAllStringSubmatch(rawhtml, -1)
	hTitle := ""
	for _, i := range h {
		text := strings.TrimSpace(ReTag.ReplaceAllString(i[1], ""))
		ratio := float32(len(text)) / float32(len(i[1]))
		//println(`"` + text + `"`)
		//println(ratio)
		if ratio < 0.5 {
			continue
		}
		if strings.HasPrefix(title, text) && len(text) > len(hTitle) {
			hTitle = text
		}
	}
	if len(hTitle) > 0 {
		title = hTitle
	}
	return strings.TrimSpace(title)
}

func getTime(text, title string) string {
	bodyText := ReHead.ReplaceAllString(text, "")
	titlePos := strings.Index(bodyText, title)
	if titlePos > 0 {
		bodyText = bodyText[titlePos:]
	}
	bodyText = ReTag.ReplaceAllString(bodyText, "")
	ret := ReDate.FindAllStringSubmatch(bodyText, -1)
	d := ""
	t := ""
	if len(ret) > 0 {
		d = ret[0][0]
		d = strings.Replace(d, `年`, `-`, -1)
		d = strings.Replace(d, `月`, `-`, -1)
		d = strings.Replace(d, `日`, ``, -1)
	}
	ret = ReTime.FindAllStringSubmatch(bodyText, -1)
	if len(ret) > 0 {
		t = ret[0][0]
	}
	return strings.TrimSpace(d + " " + t)
}

func mainText(text string) string {
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
						break
					}
				}
			}
		}
		if indexDist[i] > Threshold && !startFlag {
			for j := i + 1; j <= i+3 && j < len(indexDist); j++ {
				if indexDist[j] != 0 {
					startFlag = true
					start = i
					break
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
	return strings.TrimRightFunc(main, unicode.IsSpace)
}
