package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

var (
	file = flag.String("file", "2212137.shtml", "html file name")
)

func main() {
	flag.Parse()
	data, _ := ioutil.ReadFile(*file)
	html := string(data)
	title, content := Parse("", html)
	fmt.Println("title:\n", title, "\n=================\n\ncontent:\n", content)
}
