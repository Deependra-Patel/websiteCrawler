package main

import (
	"io/ioutil"
	"net/http"
	"strings"
)

type ClientMock struct {
	mockResponse map[string]string
}

func (c *ClientMock) Get(url string) (*http.Response, error) {
	response := c.mockResponse[url]
	if response == "" {
		panic("Url not mocked " + url)
	}
	reader := strings.NewReader(response)
	return &http.Response{Body: ioutil.NopCloser(reader)}, nil
}
