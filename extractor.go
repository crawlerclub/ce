package ce

import (
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/crawlerclub/ce/opengraph"
	"github.com/crawlerclub/ce/twitter"
	"github.com/liuzl/ip2loc"
	"github.com/liuzl/ip2tz"
	"github.com/tkuchiki/parsetime"
	"html"
	"strings"
	"time"
)

type Doc struct {
	Url      string                 `json:"url"`
	Title    string                 `json:"title"`
	Text     string                 `json:"text"`
	Html     string                 `json:"html"`
	Lang     string                 `json:"lang"`
	Country  string                 `json:"country"`
	Images   []string               `json:"images"`
	Keywords string                 `json:"keywords"`
	Tags     string                 `json:"tags"`
	Author   string                 `json:"author"`
	PubDate  time.Time              `json:"pub_date"`
	Debug    map[string]interface{} `json:"debug"`
}

func Parse(url, rawHtml string) *Doc {
	return ParsePro(url, rawHtml, "", false)
}

func ParsePro(url, rawHtml, ip string, debug bool) *Doc {
	doc := &Doc{Url: url}
	if debug {
		doc.Debug = make(map[string]interface{})
	}
	loc, err := ip2loc.Find(ip)
	doc.Country = loc
	if err != nil {
		if debug {
			doc.Debug["error"] = err.Error()
		}
	}

	tz, err := ip2tz.CountryToTz(loc)
	if err != nil {
		if debug {
			doc.Debug["error"] = err.Error()
		}
		tz = "America/New_York" // use US eastern time by default
	}
	timeParser, _ := parsetime.NewParseTime(tz)

	meta := RawMeta(rawHtml)
	metaTitle, _, metaTimes := InfoFromMeta(meta)
	if debug {
		doc.Debug["meta_times"] = metaTimes
		doc.Debug["meta_title"] = metaTitle
	}
	htmlMeta := GetMeta(meta)
	if debug {
		doc.Debug["html_meta"] = htmlMeta
	}

	raw := clean(rawHtml) // get cleaned raw html
	title := getTitle(raw)
	if metaTitle != "" && title == "" {
		title = metaTitle
	}
	ogMeta := og.GetOgp(meta)
	twtrMeta := twitter.GetTwitterCard(meta)
	if debug {
		doc.Debug["og"] = ogMeta
		doc.Debug["twitter"] = twtrMeta
	}
	if ogMeta != nil && ogMeta.Title != "" {
		doc.Title = ogMeta.Title
	} else if twtrMeta != nil && twtrMeta.Title != "" {
		doc.Title = twtrMeta.Title
	} else {
		doc.Title = title
	}

	if htmlMeta != nil {
		doc.Tags = htmlMeta.Tags
		doc.Keywords = htmlMeta.Keywords
		doc.Author = htmlMeta.Author
	}

	now := time.Now()
	var cDate time.Time
	contTime := getTime(raw, title)
	if contTime != "" {
		t, _ := timeParser.Parse(contTime)
		if now.Sub(t).Seconds() > 61 {
			cDate = t
			if debug {
				doc.Debug["content_date_str"] = contTime
			}
		}
	}

	tmp := now
	for _, metaTime := range metaTimes {
		t, _ := timeParser.Parse(metaTime)
		if now.Sub(t).Seconds() < 61 {
			continue
		}
		if tmp.After(t) {
			if debug {
				doc.Debug["meta_data_str"] = metaTime
			}
			tmp = t
		}
	}
	if !tmp.Equal(now) {
		doc.PubDate = tmp
		if debug {
			doc.Debug["pub_time_from"] = "meta"
		}
	}

	if doc.PubDate.IsZero() && !cDate.IsZero() {
		doc.PubDate = cDate
		if debug {
			doc.Debug["pub_time_from"] = "content"
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
	if doc.Text != "" {
		doc.Lang = whatlanggo.LangToString(whatlanggo.DetectLang(doc.Text))
	} else if doc.Title != "" {
		doc.Lang = whatlanggo.LangToString(whatlanggo.DetectLang(doc.Title))
	}
	return doc
}
