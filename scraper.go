package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

type Page struct {
	link            *url.URL
	sameDomainLinks []*url.URL
}

func New(link *url.URL, sameDomainLinks []*url.URL) *Page {
	page := new(Page)
	page.link = link
	page.sameDomainLinks = sameDomainLinks
	return page
}

func GetSameDomainLinks(link *url.URL) *Page {
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
		return New(link, FilterToSameDomain(link.Host, GetUrls(link, string(buffer))))
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
