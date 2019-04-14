package main

import (
	"reflect"
	"testing"
)

const startingUrl = "https://www.github.com"
const maxThreads = 100
const maxUrlsToCrawl = 10

func TestGetSiteMap_NoLinks(t *testing.T) {
	mockResponse := map[string]string{"https://www.github.com": "adsf"}
	actual := getSiteMap(&ClientMock{mockResponse}, maxThreads, startingUrl, maxUrlsToCrawl)
	expected := map[string][]string{startingUrl: {}}
	assertMapEquals(actual, expected, t)
}

func TestGetSiteMap_SingleLink(t *testing.T) {
	mockResponse := map[string]string{"https://www.github.com": "<a href=\"https://www.github.com/notification\"><a href=\"https://www.xyz.com\">",
		"https://www.github.com/notification": "Notification",
	}
	actual := getSiteMap(&ClientMock{mockResponse}, maxThreads, startingUrl, maxUrlsToCrawl)
	expected := map[string][]string{startingUrl: {"https://www.github.com/notification"},
		"https://www.github.com/notification": {},
	}
	assertMapEquals(actual, expected, t)
}

func TestGetSiteMap_SingleLink_ReferringItself(t *testing.T) {
	mockResponse := map[string]string{"https://www.github.com": "<a href=\"https://www.github.com\"><a href=\"https://www.xyz.com\">"}
	actual := getSiteMap(&ClientMock{mockResponse}, maxThreads, startingUrl, maxUrlsToCrawl)
	expected := map[string][]string{startingUrl: {startingUrl}}
	assertMapEquals(actual, expected, t)
}

func TestGetSiteMap_LimitMaxUrlsToCrawl(t *testing.T) {
	mockResponse := map[string]string{"https://www.github.com": "<a href=\"https://www.github.com/notification\"><a href=\"https://www.xyz.com\">",
		"https://www.github.com/notification": "<a href=\"https://www.github.com/status\">",
	}
	actual := getSiteMap(&ClientMock{mockResponse}, maxThreads, startingUrl, 2)
	expected := map[string][]string{startingUrl: {"https://www.github.com/notification"}, "https://www.github.com/notification": {"https://www.github.com/status"}}
	assertMapEquals(actual, expected, t)
}

func assertMapEquals(actual map[string][]string, expected map[string][]string, t *testing.T) {
	if !reflect.DeepEqual(actual, expected) {
		t.Error(actual)
	}
}
