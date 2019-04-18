package main

import (
	mapset "github.com/deckarep/golang-set"
	"log"
)

func worker(client HttpClient, links chan string, results chan *Page) {
	scraper := Scraper{client}
	for link := range links {
		results <- scraper.GetSameDomainLinks(link)
	}
}

func getSiteMap(client HttpClient, maxThreads int, startingUrl string, maxUrlsToCrawl int) map[string][]string {
	jobs := make(chan string)
	results := make(chan *Page)
	for i := 0; i < maxThreads; i++ {
		go worker(client, jobs, results)
	}
	toCrawl := mapset.NewThreadUnsafeSet()
	toCrawl.Add(startingUrl)
	alreadyCrawled := mapset.NewThreadUnsafeSet()
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
	return siteMap
}
