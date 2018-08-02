# Multilingual Web Page Content Extractor

## Introduction
`ce` is a golang package for multilingual web page content extraction. It is used to extract the content of article type web pages, such as news, blog posts, etc.

## Basic usage
```go
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	"github.com/crawlerclub/ce"
	"github.com/crawlerclub/dl"
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
```

## Fields

`ce` can extract the following fields from raw web htmls:
* `title`: the title of article
* `text`: the main content of article in plain text
* `html`: the main content of article with basic html format, images included
* `publish_date`: the publish time of article
* `language`: the language of article
* `location`: the country code
* `author`: the author of artile
* `images`: the images used in the article
