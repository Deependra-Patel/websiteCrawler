package main

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"io/ioutil"
)

func worker(links chan string, results chan *Page) {
	for link := range links {
		results <- GetSameDomainLinks(link)
	}
}

func main() {
	const siteMapFile = "./siteMap.json"
	const startingUrl = "http://github.com/"
	const maxUrlsToCrawl = 5000
	const maxThreads = 50

	jobs := make(chan string)
	results := make(chan *Page)
	for i := 0; i < maxThreads; i++ {
		go worker(jobs, results)
	}

	toCrawl := mapset.NewSet(startingUrl)
	alreadyCrawled := mapset.NewSet()
	siteMap := make(map[string][]string)

	pendingResults := 0
	for alreadyCrawled.Cardinality() < maxUrlsToCrawl {
		if toCrawl.Cardinality() > 0 && pendingResults < maxThreads {
			linkToGet := toCrawl.Pop().(string)
			jobs <- linkToGet
			pendingResults++
			alreadyCrawled.Add(linkToGet)
		} else if pendingResults == 0 {
			break
		} else {
			page := <-results
			pendingResults--
			links := page.sameDomainLinks
			for _, link := range links {
				if !alreadyCrawled.Contains(link) {
					toCrawl.Add(link)
				}
			}
			siteMap[page.link] = links
		}
	}

	jsonSiteMap, err := json.Marshal(siteMap)
	check(err)
	err = ioutil.WriteFile(siteMapFile, jsonSiteMap, 0644)
	check(err)
}
