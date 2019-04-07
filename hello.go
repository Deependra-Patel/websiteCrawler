package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func main() {
	link, _ := url.Parse("https://github.com")
	urls := GetExternalLinks(link)
	fmt.Print(urls)
}

func GetExternalLinks(link *url.URL) []*url.URL {
	log.Print("Getting ", link)
	resp, err := http.Get(link.String())
	if err != nil {
		log.Panicf("Call failed for link: %s with error: %s", link, err)
		return nil
	} else {
		buffer := make([]byte, 1024*1024)
		count, err := resp.Body.Read(buffer)
		if err != nil {
			log.Panicf("Reading response body for link %s failed with %s", link, err)
		}
		log.Printf("Number of bytes read %d", count)
		return GetUrls(link, string(buffer))
	}
}

func GetUrls(parse *url.URL, body string) []*url.URL {
	compile, err := regexp.Compile("href=\"[^\"]*\"")
	if err != nil {
		panic(err)
	}
	links := make([]*url.URL, 0)
	for _, str := range compile.FindAllString(body, -1) {
		link, err := parse.Parse(str[6 : len(str)-1])
		if err != nil {
			panic(err)
		}
		links = append(links, link)
	}
	return links
}
