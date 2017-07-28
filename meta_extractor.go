package ce

import (
	"fmt"
	"github.com/tkuchiki/parsetime"
	"html"
	"regexp"
	"strings"
)

var (
	ReMeta = regexp.MustCompile(`(?ims)<meta.*?>`)
	ReKV   = regexp.MustCompile(`(?ims)([^\s]+?)\s*?=\s*?"(.+?)"|'(.+?)'`)

	TimeParser, _ = parsetime.NewParseTime()
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

func RawMeta(raw string) []map[string]string {
	var list []map[string]string
	metas := ReMeta.FindAllStringSubmatch(raw, -1)
	for i := range metas {
		m := make(map[string]string)
		kvs := ReKV.FindAllStringSubmatch(metas[i][0], -1)
		for _, kv := range kvs {
			m[strings.ToLower(kv[1])] = html.UnescapeString(kv[2])
		}
		list = append(list, m)
	}
	return list
}

func FromMeta(meta []map[string]string) {
	fmt.Println("FromMeta")
	for _, m := range meta {
		content, has := m["content"]
		if !has {
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
			fmt.Println(name, content)
		case strings.Contains(name, "desc"):
			fmt.Println(name, content)
		case strings.Contains(name, "date") ||
			strings.Contains(name, "time") ||
			strings.Contains(name, "_at"):
			fmt.Println(name, content)
			t, _ := TimeParser.Parse(content)
			fmt.Println(t)
		}
	}
	fmt.Println("\n\n")
}
