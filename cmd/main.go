package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"crawler.club/ce"
	"crawler.club/dl"
)

var (
	url = flag.String("url",
		"http://china.huanqiu.com/article/2017-07/11034896.html",
		"news url")
	debug = flag.Bool("debug", false, "debug mode")
)

func main() {
	flag.Parse()
	res := dl.DownloadUrl(*url)
	if res.Error != nil {
		fmt.Println(res.Error)
		return
	}

	items := strings.Split(res.RemoteAddr, ":")
	ip := ""
	if len(items) > 0 {
		ip = items[0]
	}
	doc := ce.ParsePro(*url, res.Text, ip, *debug)
	j, _ := json.Marshal(doc)
	fmt.Println(string(j))
}
