package main

import (
	"encoding/json"
	"github.com/deckarep/golang-set"
	"io/ioutil"
	"log"
	"time"
)

func worker(links chan string, results chan *Page) {
	for link := range links {
		results <- GetSameDomainLinks(link)
	}
}

func main() {
	const siteMapFile = "./siteMap.json"
	const startingUrl = "https://www.facebook.com/"
	const maxUrlsToCrawl = 3000
	const maxThreads = 100

	startTime := time.Now().Unix()
	jobs := make(chan string)
	results := make(chan *Page)
	for i := 0; i < maxThreads; i++ {
		go worker(jobs, results)
	}

	toCrawl := mapset.NewSet(startingUrl)
	alreadyCrawled := mapset.NewSet()
	siteMap := make(map[string][]string)

	pendingResults := 0
	for {
		if alreadyCrawled.Cardinality() < maxUrlsToCrawl && toCrawl.Cardinality() > 0 && pendingResults < maxThreads {
			linkToGet := toCrawl.Pop().(string)
			jobs <- linkToGet
			pendingResults++
			alreadyCrawled.Add(linkToGet)
		} else if pendingResults == 0 {
			close(jobs)
			close(results)
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
	log.Println("Number of pages crawled:", alreadyCrawled.Cardinality())
	log.Println("Number of pages left in queue:", toCrawl.Cardinality())
	log.Println("Time taken in seconds: ", time.Now().UTC().Unix()-startTime)

	jsonSiteMap, err := json.Marshal(siteMap)
	check(err)
	err = ioutil.WriteFile(siteMapFile, jsonSiteMap, 0644)
	check(err)
}
