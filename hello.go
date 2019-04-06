package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	url := "https://github.com"
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Call failed for url: %s with error: %s", url, err)
	} else {
		buffer := make([]byte, 1024*1024)
		_, _ = resp.Body.Read(buffer)
		fmt.Print(string(buffer))
	}
	fmt.Print("Test")
}
