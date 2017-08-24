package ce

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestTitle(t *testing.T) {
	files := []string{
		"./test_data/cbsnews.html",
		"./test_data/dailycaller.html",
		"./test_data/huanqiu.html",
		"./test_data/sina.html",
		"./test_data/weiyangx.html",
	}
	for _, file := range files {
		bytes, _ := ioutil.ReadFile(file)
		doc := ParsePro("", string(bytes), "", true)
		j, _ := json.Marshal(doc.Title)
		t.Log(string(j))
	}
}
