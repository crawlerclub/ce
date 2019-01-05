package ce

import (
	"html"
	"regexp"
	"strings"
)

var (
	ReMeta = regexp.MustCompile(`(?ims)<meta.*?>`)
	ReKV   = regexp.MustCompile(`(?ims)([^\s]+?)\s*?=\s*?"(.+?)"|'(.+?)'`)
)

type Meta struct {
	Keywords    string `json:"keywords"`
	Tags        string `json:"tags"`
	Description string `json:"description"`
	Author      string `json:"author"`
}

const (
	metaKeywords    = "keywords"
	metaTags        = "tags"
	metaDescription = "description"
	metaAuthor      = "author"
)

func RawMeta(raw string) []map[string]string {
	var list []map[string]string
	metas := ReMeta.FindAllStringSubmatch(raw, -1)
	for i := range metas {
		m := make(map[string]string)
		kvs := ReKV.FindAllStringSubmatch(metas[i][0], -1)
		for _, kv := range kvs {
			m[strings.ToLower(kv[1])] = strings.TrimSpace(html.UnescapeString(kv[2]))
		}
		list = append(list, m)
	}
	return list
}

func GetMeta(meta []map[string]string) *Meta {
	var obj *Meta
	for _, m := range meta {
		name, has := m["name"]
		if !has {
			continue
		}
		content, has := m["content"]
		if !has {
			continue
		}
		switch name {
		case metaKeywords:
			if obj == nil {
				obj = new(Meta)
			}
			obj.Keywords = content
		case metaTags:
			if obj == nil {
				obj = new(Meta)
			}
			obj.Tags = content
		case metaDescription:
			if obj == nil {
				obj = new(Meta)
			}
			obj.Description = content
		case metaAuthor:
			if obj == nil {
				obj = new(Meta)
			}
			obj.Author = content
		}
	}
	return obj
}

func InfoFromMeta(meta []map[string]string) (string, string, []string) {
	title := ""
	desc := ""
	var times []string
	for _, m := range meta {
		content, has := m["content"]
		if !has || content == "" {
			continue
		}
		name, has := m["name"]
		if !has {
			name, has = m["property"]
			if !has {
				continue
			}
		}
		switch {
		case strings.Contains(name, "title"):
			if title == "" || len(content) < len(title) {
				title = content
			}
		case strings.Contains(name, "desc"):
			if desc == "" || len(content) > len(desc) {
				desc = content
			}
		case strings.Contains(name, "date") ||
			strings.Contains(name, "time") ||
			strings.Contains(name, "_at"):
			if ReDate.MatchString(content) || ReTime.MatchString(content) {
				times = append(times, content)
			}
		}
	}
	return title, desc, times
}
