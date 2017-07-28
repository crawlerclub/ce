package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/crawlerclub/ce"
	"github.com/crawlerclub/x/downloader"
	"github.com/crawlerclub/x/types"
)

var (
	url = flag.String("url",
		"http://china.huanqiu.com/article/2017-07/11034896.html",
		"news url")
)

func main() {
	flag.Parse()
	req := &types.HttpRequest{Url: *url, Method: "GET", UseProxy: false, Platform: "pc"}
	res := downloader.Download(req)
	if res.Error != nil {
		fmt.Println(res.Error)
		return
	}
	//ce.Parse("", res.Text)
	doc := ce.Parse(*url, res.Text)
	j, _ := json.Marshal(doc)
	fmt.Println(string(j))
	//fmt.Println("title:\n", title, "\n=================\n\ncontent:\n", content)
}
