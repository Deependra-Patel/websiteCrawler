package main

import (
	"github.com/deckarep/golang-set"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func main() {
	domain := "https://github.com/"
	maxUrlsToCrawl := 10

	link, _ := url.Parse(domain)
	host := link.Host
	toCrawl := mapset.NewSet(link)
	alreadyCrawled := mapset.NewSet()
	for i := 0; i < maxUrlsToCrawl; i++ {
		linkToGet := toCrawl.Pop().(*url.URL)
		links := GetExternalLinks(linkToGet)
		alreadyCrawled.Add(*linkToGet)
		sameDomainLinks := filterToSameDomain(host, links)
		for _, link := range sameDomainLinks {
			if !alreadyCrawled.Contains(*link) {
				toCrawl.Add(link)
			}
		}
	}
}

func filterToSameDomain(host string, links []*url.URL) []*url.URL {
	sameDomainLinks := make([]*url.URL, 0)
	for _, link := range links {
		if link.Host == host {
			sameDomainLinks = append(sameDomainLinks, link)
		}
	}
	return sameDomainLinks
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
		if err != nil && err != io.EOF {
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
