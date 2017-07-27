package extractors

import (
	"regexp"
)

var (
	ReMeta = regexp.MustCompile(`(?ims)<meta.*?>`)
	ReKV   = regexp.MustCompile(`(?ims)([^\s]+?)\s*?=\s*?"(.+?)"|'(.+?)'`)
)

func RawMeta(raw string) []map[string]string {
	var list []map[string]string
	metas := ReMeta.FindAllStringSubmatch(raw, -1)
	for i := range metas {
		m := make(map[string]string)
		kvs := ReKV.FindAllStringSubmatch(metas[i][0], -1)
		for _, kv := range kvs {
			m[kv[1]] = kv[2]
		}
		list = append(list, m)
	}
	return list
}
