package ce

import (
	"fmt"
	"github.com/tkuchiki/parsetime"
	"html"
	"strings"
	"time"
)

var (
	TimeParser, _ = parsetime.NewParseTime()
)

type Doc struct {
	Url      string    `json:"url"`
	Title    string    `json:"title"`
	Text     string    `json:"text"`
	Html     string    `json:"html"`
	Images   []string  `json:"images"`
	Keywords string    `json:"keywords"`
	Tags     string    `json:"tags"`
	Author   string    `json:"author"`
	Date     time.Time `json:"date"`
}

func Parse(url, rawHtml string) *Doc {
	doc := &Doc{Url: url}

	meta := RawMeta(rawHtml)
	metaTitle, _, metaTimes := InfoFromMeta(meta)
	htmlMeta := GetMeta(meta)

	raw := clean(rawHtml) // get cleaned raw html
	title := getTitle(raw)
	// choose short title
	if metaTitle != "" && strings.HasPrefix(title, metaTitle) {
		title = metaTitle
	}

	doc.Title = title
	doc.Tags = htmlMeta.Tags
	doc.Keywords = htmlMeta.Keywords
	doc.Author = htmlMeta.Author

	now := time.Now()
	contTime := getTime(raw, title)
	fmt.Println("[" + contTime + "]")
	if contTime != "" {
		t, _ := TimeParser.Parse(contTime)
		if now.Sub(t).Seconds() > 61 {
			doc.Date = t
		}
	}
	if doc.Date.IsZero() {
		tmp := now
		for _, metaTime := range metaTimes {
			t, _ := TimeParser.Parse(metaTime)
			if now.Sub(t).Seconds() < 61 {
				continue
			}
			if tmp.After(t) {
				tmp = t
			}
		}
		if !tmp.Equal(now) {
			doc.Date = tmp
		}
	}

	images := make(map[string]string)
	ret := ReImg.FindAllStringSubmatch(raw, -1)
	for _, r := range ret {
		if len(r) <= 0 {
			continue
		}
		src := ReImgSrc.FindAllStringSubmatch(r[0], -1)
		if len(src) <= 0 {
			continue
		}
		md := MD5(r[0])
		images[md] = fmt.Sprintf("<img src=\"%s\" />", src[0][1])
		raw = strings.Replace(raw, r[0], md, -1)
	}
	// get raw text
	text := html.UnescapeString(ReTag.ReplaceAllString(raw, ""))
	content := mainText(text)
	for k, v := range images {
		content = strings.Replace(content, k, v, -1)
	}
	doc.Html = content
	return doc
}
