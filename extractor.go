package ce

import (
	"fmt"
	"github.com/crawlerclub/ce/opengraph"
	"github.com/crawlerclub/ce/twitter"
	"github.com/tkuchiki/parsetime"
	"html"
	"strings"
	"time"
)

var (
	TimeParser, _ = parsetime.NewParseTime()
)

type Doc struct {
	Url      string      `json:"url"`
	Title    string      `json:"title"`
	Text     string      `json:"text"`
	Html     string      `json:"html"`
	Images   []string    `json:"images"`
	OgImages interface{} `json:"og_images"`
	Keywords string      `json:"keywords"`
	Tags     string      `json:"tags"`
	Author   string      `json:"author"`
	Date     time.Time   `json:"date"`
}

func Parse(url, rawHtml string) *Doc {
	doc := &Doc{Url: url}

	meta := RawMeta(rawHtml)
	metaTitle, _, metaTimes := InfoFromMeta(meta)
	htmlMeta := GetMeta(meta)

	raw := clean(rawHtml) // get cleaned raw html
	title := getTitle(raw)
	if metaTitle != "" && title == "" {
		title = metaTitle
	}
	ogMeta := og.GetOgp(meta)
	twtrMeta := twitter.GetTwitterCard(meta)
	if ogMeta != nil && ogMeta.Title != "" {
		doc.Title = ogMeta.Title
	} else if twtrMeta != nil && twtrMeta.Title != "" {
		doc.Title = twtrMeta.Title
	} else {
		doc.Title = title
	}

	if ogMeta != nil && len(ogMeta.Image) > 0 {
		doc.OgImages = ogMeta.Image
	}
	if htmlMeta != nil {
		doc.Tags = htmlMeta.Tags
		doc.Keywords = htmlMeta.Keywords
		doc.Author = htmlMeta.Author
	}

	now := time.Now()
	contTime := getTime(raw, title)
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
	plainText := content
	lines := strings.Split(content, "\n")
	for i, _ := range lines {
		lines[i] = "<p>" + html.EscapeString(lines[i]) + "</p>"
	}
	content = strings.Join(lines, "\n")

	for k, v := range images {
		if strings.Contains(content, k) {
			content = strings.Replace(content, k, v, -1)
			plainText = strings.Replace(plainText, k, "", -1)
			doc.Images = append(doc.Images, v)
		}
	}
	doc.Html = content
	doc.Text = plainText
	return doc
}
