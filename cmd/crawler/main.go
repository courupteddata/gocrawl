package main

import (
	"fmt"
	"gocrawl"
	"net/url"
)

func main() {
	crawler := gocrawl.NewClient()
	rawUrl := "https://google.com"
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		fmt.Printf("Failed to parse %v, err=%v", rawUrl, err)
		return
	}
	crawler.Crawl(*parsedUrl)
}
