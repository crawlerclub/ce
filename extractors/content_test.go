package extractors

import (
	"io/ioutil"
	"testing"
)

func TestContent(t *testing.T) {
	files := []string{
		"./test_data/cbsnews.html",
		"./test_data/dailycaller.html",
		"./test_data/huanqiu.html",
		"./test_data/sina.html",
	}
	for _, file := range files {
		bytes, _ := ioutil.ReadFile(file)
		title, content := Parse("", string(bytes))
		t.Log(title)
		t.Log(content)
	}
}