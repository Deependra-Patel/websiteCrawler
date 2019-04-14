package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	const siteMapFile = "./siteMap.json"
	const startingUrl = "https://www.facebook.com/"
	const maxUrlsToCrawl = 3000
	const maxThreads = 100

	startTime := time.Now().Unix()
	siteMap := getSiteMap(http.DefaultClient, maxThreads, startingUrl, maxUrlsToCrawl)
	log.Println("Time taken in seconds: ", time.Now().UTC().Unix()-startTime)

	jsonSiteMap, err := json.Marshal(siteMap)
	check(err)
	err = ioutil.WriteFile(siteMapFile, jsonSiteMap, 0644)
	check(err)
}
