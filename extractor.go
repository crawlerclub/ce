package ce

import (
	"fmt"
	"html"
	"net/url"
	"strings"
	"time"

	"github.com/abadojack/whatlanggo"
	"github.com/crawlerclub/ce/opengraph"
	"github.com/crawlerclub/ce/twitter"
	"github.com/joeguo/tldextract"
	"github.com/liuzl/goutil"
	"github.com/liuzl/ip2loc"
	"github.com/liuzl/ip2tz"
)

type Doc struct {
	Url             string                 `json:"url"`
	From            string                 `json:"from"`
	SiteInfo        *tldextract.Result     `json:"site_info"`
	CanonicalUrl    string                 `json:"canonical_url"`
	Title           string                 `json:"title"`
	Text            string                 `json:"text"`
	Html            string                 `json:"html"`
	Language        string                 `json:"language"`
	Location        string                 `json:"location"`
	Favicon         string                 `json:"favicon"`
	Images          []string               `json:"images"`
	Tags            string                 `json:"tags"`
	Author          string                 `json:"author"`
	Published       string                 `json:"published"`
	PublishedParsed time.Time              `json:"published_parsed"`
	Debug           map[string]interface{} `json:"debug,omitempty"`
}

var tldExtractor, _ = tldextract.New("./tld.cache", false)

func Parse(rawurl, rawHtml string) *Doc {
	return ParsePro(rawurl, rawHtml, "", false)
}

func ParsePro(rawurl, rawHtml, ip string, debug bool) *Doc {
	doc := &Doc{Url: rawurl}
	if debug {
		doc.Debug = make(map[string]interface{})
	}

	pUrl, err := url.Parse(rawurl)
	if err != nil {
		if debug {
			doc.Debug["url_error"] = err.Error()
		}
	} else {
		doc.From = pUrl.Hostname()
		doc.Favicon = fmt.Sprintf("%s://%s/favicon.ico", pUrl.Scheme, pUrl.Host)
		doc.SiteInfo = tldExtractor.Extract(rawurl)
		if doc.SiteInfo.Root != "" && doc.SiteInfo.Tld != "" {
			doc.From = fmt.Sprintf("%s.%s", doc.SiteInfo.Root, doc.SiteInfo.Tld)
		}
	}

	favicon := getFavicon(rawHtml)
	if favicon != "" {
		absUrl, err := goutil.MakeAbsoluteURL(favicon, rawurl)
		if err == nil {
			doc.Favicon = absUrl
		}
	}

	loc, err := ip2loc.Find(ip)
	doc.Location = loc
	if err != nil {
		if debug {
			doc.Debug["ip2loc_error"] = err.Error()
		}
	}

	tz, err := ip2tz.CountryToTz(loc)
	if err != nil {
		if debug {
			doc.Debug["ip2tz_error"] = err.Error()
		}
		tz = "America/New_York" // use US eastern time by default
	}

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

	// process title begin
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
	doc.Title = html.UnescapeString(doc.Title)
	ret := ReTitleNoNoisy.FindAllStringSubmatch(doc.Title, -1)
	if len(ret) > 0 {
		doc.Title = strings.TrimSpace(ret[0][0])
	}
	// process title end

	if ogMeta != nil && ogMeta.Url != "" {
		doc.CanonicalUrl = ogMeta.Url
	}

	now := time.Now()
	var cDate time.Time
	contTime := getTime(raw, title)
	if contTime != "" {
		t := ParseTime(tz, contTime)
		if now.Sub(t).Seconds() > 61 {
			cDate = t
			doc.Published = contTime
			if debug {
				doc.Debug["content_date_str"] = contTime
			}
		}
	}

	tmp := now
	for _, metaTime := range metaTimes {
		t := ParseTime(tz, metaTime)
		if now.Sub(t).Seconds() < 61 {
			continue
		}
		if tmp.After(t) {
			doc.Published = metaTime
			if debug {
				doc.Debug["meta_data_str"] = metaTime
			}
			tmp = t
		}
	}
	if !tmp.Equal(now) {
		doc.PublishedParsed = tmp
		if debug {
			doc.Debug["pub_time_from"] = "meta"
		}
	}

	if doc.PublishedParsed.IsZero() && !cDate.IsZero() {
		doc.PublishedParsed = cDate
		if debug {
			doc.Debug["pub_time_from"] = "content"
		}
	}

	images := make(map[string]string)
	ret = ReImg.FindAllStringSubmatch(raw, -1)
	for _, r := range ret {
		if len(r) <= 0 {
			continue
		}
		md := goutil.MD5(r[0])
		img := ""
		src := ReImgSrc.FindAllStringSubmatch(r[0], -1)
		for _, s := range src {
			u := s[1]
			if strings.HasPrefix(u, "data:image") {
				img = u
				continue
			}

			if strings.HasPrefix(u, "//") {
				u = "http:" + u
			}
			absUrl, err := goutil.MakeAbsoluteURL(u, rawurl)
			if err == nil {
				img = absUrl
				break
			}
		}
		if _, ok := IgnoreImgs[img]; ok {
			continue
		}
		if img != "" {
			images[md] = img
			raw = strings.Replace(raw, r[0], md, -1)
		}
	}
	// get raw text
	text := html.UnescapeString(ReTag.ReplaceAllString(raw, ""))
	content := mainText(text)
	plainText := content
	lines := strings.Split(content, "\n")
	for i := range lines {
		lines[i] = "<p>" + html.EscapeString(lines[i]) + "</p>"
	}
	content = strings.Join(lines, "\n")

	for k, v := range images {
		if strings.Contains(content, k) {
			content = strings.Replace(content, k,
				fmt.Sprintf("<img src=\"%s\" />", v), -1)
			plainText = strings.Replace(plainText, k, "", -1)
			doc.Images = append(doc.Images, v)
		}
	}
	if len(doc.Images) == 0 && ogMeta != nil && len(ogMeta.Image) > 0 {
		doc.Images = append(doc.Images, ogMeta.Image[0].Url)
		content = fmt.Sprintf("<p><img src=\"%s\" /></p>\n%s",
			ogMeta.Image[0].Url, content)
	}

	doc.Html = content
	doc.Text = plainText

	if htmlMeta != nil {
		doc.Tags = htmlMeta.Tags
		if doc.Tags == "" && htmlMeta.Keywords != "" {
			doc.Tags = htmlMeta.Keywords
		}
		if doc.Text == "" && htmlMeta.Description != "" {
			doc.Text = htmlMeta.Description
		}
		doc.Author = htmlMeta.Author
	}

	if doc.Text != "" {
		doc.Language = whatlanggo.LangToString(whatlanggo.DetectLang(doc.Text))
	} else if doc.Title != "" {
		doc.Language = whatlanggo.LangToString(whatlanggo.DetectLang(doc.Title))
	}
	return doc
}
