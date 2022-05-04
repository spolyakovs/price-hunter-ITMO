package api

import (
	"io/ioutil"
	"log"
	"net/http"
)

// TODO: get steam key from config
// TODO: probably implement common interface for APIs (baseURL, key, ...)
// TODO: write steamAPI connect
func get_test() {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/posts")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	resp.Body.Close()

	sb := string(body)
	log.Printf(sb)
}
