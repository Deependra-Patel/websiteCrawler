package main

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"net/url"
)

func worker(links chan *url.URL, results chan *Page) {
	for link := range links {
		results <- GetSameDomainLinks(link)
	}
}

func main() {
	const siteMapFile = "./siteMap.json"
	const startingUrl = "https://github.com/"
	const maxUrlsToCrawl = 100
	const maxThreads = 20

	jobs := make(chan *url.URL)
	results := make(chan *Page)
	for i := 0; i < maxThreads; i++ {
		go worker(jobs, results)
	}

	startLink, err := url.Parse(startingUrl)
	check(err)
	toCrawl := mapset.NewSet(startLink)
	alreadyCrawled := mapset.NewSet()
	siteMap := make(map[string][]string)

	pendingResults := 0
	for alreadyCrawled.Cardinality() < maxUrlsToCrawl {
		if toCrawl.Cardinality() > 0 && pendingResults < maxThreads {
			linkToGet := toCrawl.Pop().(*url.URL)
			jobs <- linkToGet
			pendingResults++
			alreadyCrawled.Add(*linkToGet)
		} else if pendingResults == 0 {
			break
		} else {
			page := <-results
			pendingResults--
			links := page.sameDomainLinks
			linksStr := make([]string, len(links))
			for _, link := range links {
				if !alreadyCrawled.Contains(*link) {
					toCrawl.Add(link)
				}
				linksStr = append(linksStr, link.String())
			}
			siteMap[page.link.String()] = linksStr
		}
	}

	jsonSiteMap, err := json.Marshal(siteMap)
	check(err)
	err = ioutil.WriteFile(siteMapFile, jsonSiteMap, 0644)
	check(err)
}
