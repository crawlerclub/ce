package extractors

import (
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
			m[strings.ToLower(kv[1])] = kv[2]
		}
		list = append(list, m)
	}
	return list
}
