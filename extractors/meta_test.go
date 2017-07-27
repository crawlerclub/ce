package extractors

import (
	"encoding/json"
	"io/ioutil"
	"testing"
)

func TestMeta(t *testing.T) {
	files := []string{"./test_data/cbsnews.html", "./test_data/dailycaller.html", "./test_data/huanqiu.html"}
	for _, file := range files {
		bytes, _ := ioutil.ReadFile(file)
		ret := Meta(string(bytes))
		data, _ := json.Marshal(ret)
		t.Log(string(data))
	}
}
