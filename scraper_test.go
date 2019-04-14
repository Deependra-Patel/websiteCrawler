package main

import (
	"net/url"
	"testing"
)

func TestGetUrls_absolute(t *testing.T) {
	pageUrl, _ := url.Parse("https://www.github.com/Deependra-Patel/")
	actual := GetUrls(pageUrl, "<html><a href=\"https://www.deependra.com/\"><div><a href=\"http://www.xyz.com/\"><div></html>")
	var expected [2]*url.URL
	expected[0], _ = url.Parse("https://www.deependra.com/")
	expected[1], _ = url.Parse("http://www.xyz.com/")
	if len(actual) != 2 || *actual[0] != *expected[0] || *actual[1] != *expected[1] {
		t.Error(actual, expected)
	}
}

func TestGetUrls_relativeRoot(t *testing.T) {
	pageUrl, _ := url.Parse("https://www.github.com/Deependra-Patel/")
	actual := GetUrls(pageUrl, "<html><a href=\"/notifications\"></html>")

	expected, _ := url.Parse("https://www.github.com/notifications")
	if len(actual) != 1 || *actual[0] != *expected {
		t.Error(actual, expected)
	}
}

func TestGetUrls_relativeSibling(t *testing.T) {
	pageUrl, _ := url.Parse("https://www.github.com/Deependra-Patel/")
	actual := GetUrls(pageUrl, "<html><a href=\"./notifications\"></html>")

	expected, _ := url.Parse("https://www.github.com/Deependra-Patel/notifications")
	if len(actual) != 1 || *actual[0] != *expected {
		t.Error(actual, expected)
	}
}

func TestFilterToSameDomain(t *testing.T) {
	host := "www.github.com"
	url1, _ := url.Parse("https://www.github.com/Deependra-Patel")
	url2, _ := url.Parse("https://www.xyz.com/a")
	actual := FilterToSameDomain(host, []*url.URL{url1, url2})

	expected := "https://www.github.com/Deependra-Patel"
	if len(actual) != 1 || actual[0] != expected {
		t.Error(actual, expected)
	}
}

func TestScraper_GetSameDomainLinks(t *testing.T) {
	mockResponse := make(map[string]string, 0)
	pageUrl := "https://www.github.com"
	mockResponse[pageUrl] = "<a href=\"https://www.github.com/notification\"><a href=\"https://www.xyz.com\">"
	clientMock := ClientMock{mockResponse}

	scraper := Scraper{&clientMock}
	actual := scraper.GetSameDomainLinks(pageUrl)
	if actual.link != pageUrl {
		t.Error(actual, pageUrl)
	}
	if len(actual.sameDomainLinks) != 1 || actual.sameDomainLinks[0] != "https://www.github.com/notification" {
		t.Error(actual, mockResponse)
	}
}
