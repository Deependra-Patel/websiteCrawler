package main

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func main() {
	siteMapFile := "./siteMap.json"
	startingUrl := "https://github.com/"
	maxUrlsToCrawl := 10

	startLink, err := url.Parse(startingUrl)
	check(err)
	toCrawl := mapset.NewSet(startLink)
	alreadyCrawled := mapset.NewSet()
	siteMap := make(map[string][]string)
	for i := 0; i < maxUrlsToCrawl; i++ {
		linkToGet := toCrawl.Pop().(*url.URL)
		links := GetSameDomainLinks(linkToGet)
		alreadyCrawled.Add(*linkToGet)

		linksStr := make([]string, len(links))
		for _, link := range links {
			if !alreadyCrawled.Contains(*link) {
				toCrawl.Add(link)
			}
			linksStr = append(linksStr, link.String())
		}
		siteMap[linkToGet.String()] = linksStr
	}

	jsonSiteMap, err := json.Marshal(siteMap)
	check(err)
	err = ioutil.WriteFile(siteMapFile, jsonSiteMap, 0644)
	check(err)
}

func GetSameDomainLinks(link *url.URL) []*url.URL {
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
		return FilterToSameDomain(link.Host, GetUrls(link, string(buffer)))
	}
}

func GetUrls(parse *url.URL, body string) []*url.URL {
	compile, err := regexp.Compile("href=\"[^\"]*\"")
	check(err)
	links := make([]*url.URL, 0)
	for _, str := range compile.FindAllString(body, -1) {
		link, err := parse.Parse(str[6 : len(str)-1])
		check(err)
		links = append(links, link)
	}
	return links
}

func FilterToSameDomain(host string, links []*url.URL) []*url.URL {
	sameDomainLinks := make([]*url.URL, 0)
	for _, link := range links {
		if link.Host == host {
			sameDomainLinks = append(sameDomainLinks, link)
		}
	}
	return sameDomainLinks
}
