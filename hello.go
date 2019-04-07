package main

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"net/url"
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

	ch := make(chan []*url.URL)
	for i := 0; i < maxUrlsToCrawl; i++ {
		linkToGet := toCrawl.Pop().(*url.URL)
		go GetSameDomainLinks(linkToGet, ch)
		alreadyCrawled.Add(*linkToGet)
		links := <-ch

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
