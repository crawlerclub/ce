package ce

import (
	"encoding/json"
	"github.com/crawlerclub/ce/opengraph"
	"github.com/crawlerclub/ce/twitter"
	"io/ioutil"
	"testing"
)

func TestMeta(t *testing.T) {
	files := []string{
		"./test_data/cbsnews.html",
		"./test_data/dailycaller.html",
		"./test_data/huanqiu.html",
		"./test_data/sina.html",
	}
	for _, file := range files {
		bytes, _ := ioutil.ReadFile(file)
		ret := RawMeta(string(bytes))
		data, _ := json.Marshal(ret)
		if false {
			t.Log(string(data))
		}
		ogp := og.GetOgp(ret)
		j, _ := json.Marshal(ogp)
		t.Log("og: ", string(j))

		tw := twitter.GetTwitterCard(ret)
		j, _ = json.Marshal(tw)
		t.Log("twitter: ", string(j))

		mt := GetMeta(ret)
		j, _ = json.Marshal(mt)
		t.Log("meta: ", string(j))

		InfoFromMeta(ret)
	}
}
