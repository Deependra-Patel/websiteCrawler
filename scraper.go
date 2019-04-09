package main

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

var hrefMatcher, _ = regexp.Compile("href=\"[^\"]*\"")

type Page struct {
	link            string
	sameDomainLinks []string
}

func New(link string, sameDomainLinks []string) *Page {
	page := new(Page)
	page.link = link
	page.sameDomainLinks = sameDomainLinks
	return page
}

func GetSameDomainLinks(link string) *Page {
	log.Print("Getting ", link)
	resp, err := http.Get(link)
	if err != nil {
		log.Panicf("Call failed for link: %s with error: %s", link, err)
		return nil
	} else {
		buffer := make([]byte, 1024*1024)
		_, err := resp.Body.Read(buffer)
		if err != nil && err != io.EOF {
			log.Panicf("Reading response body for link %s failed with %s", link, err)
		}
		parsedLink, err := url.Parse(link)
		check(err)
		return New(link, FilterToSameDomain(parsedLink.Host, GetUrls(parsedLink, string(buffer))))
	}
}

func GetUrls(pageUrl *url.URL, body string) []*url.URL {
	links := make([]*url.URL, 0)
	for _, str := range hrefMatcher.FindAllString(body, -1) {
		link, err := pageUrl.Parse(str[6 : len(str)-1])
		check(err)
		links = append(links, link)
	}
	return links
}

func FilterToSameDomain(host string, links []*url.URL) []string {
	sameDomainLinks := make([]string, 0)
	for _, link := range links {
		if link.Host == host {
			sameDomainLinks = append(sameDomainLinks, link.String())
		}
	}
	return sameDomainLinks
}
