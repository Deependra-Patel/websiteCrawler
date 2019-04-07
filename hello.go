package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
)

func main() {
	urls := GetExternalLinks("https://github.com")
	fmt.Print(urls)
}

func GetExternalLinks(url string) []string {
	log.Print("Getting ", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Panicf("Call failed for url: %s with error: %s", url, err)
		return nil
	} else {
		buffer := make([]byte, 1024*1024)
		count, err := resp.Body.Read(buffer)
		if err != nil {
			log.Panicf("Reading response body for url %s failed with %s", url, err)
		}
		log.Printf("Number of bytes read %d", count)
		return GetUrls(string(buffer))
	}
}

func GetUrls(body string) []string {
	compile, err := regexp.Compile("href=\"[^\"]*\"")
	if err != nil {
		panic(err)
	}
	return compile.FindAllString(body, -1)
}
