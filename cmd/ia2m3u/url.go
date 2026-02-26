package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

// func getUrlJSON(client *http.Client, url string, cachePrefix *string, useCache bool, items *scrapeItems, cursor string, cache *Cache) error {
func getUrlJSON(client *http.Client, url string, useCache bool, items interface{}, cursor string, cache *Cache) error {

	var body []byte

	if useCache {
		if body = cache.GetKey(url); body != nil {

		} else {
			//res, err := http.Get(url)

			req, err := http.NewRequest(http.MethodGet, url, nil)
			if err != nil {
				fmt.Printf("client: could not create request: %s\n", err)
				log.Fatal()
			}

			res, err := client.Do(req)
			if err != nil {
				fmt.Printf("client: error making http request: %s\n", err)
				log.Fatal()
			}

			if res.StatusCode != 200 {
				body, err = io.ReadAll(res.Body)
				if err == nil {
					log.Println("Error. Response body:")
					log.Println("--------------------------------------------------------------")
					log.Println(string(body))
					log.Println("--------------------------------------------------------------")
				}
				return fmt.Errorf("Failing http code %d (!200)", res.StatusCode)
			}
			if err != nil {
				log.Println("Status code", res.StatusCode)
				log.Println(url)
				log.Println(err)
				log.Fatal(err)
			}

			body, err = io.ReadAll(res.Body)
			if err != nil {
				log.Println(err)
				return err
			}
			cache.AddToCache(url, body)
			//time.Sleep(400 * time.Millisecond)
		}
	}

	dec := json.NewDecoder(bytes.NewBuffer(body))
	err := dec.Decode(items)

	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	return err

}

func newClient() *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        10,               // Maximum idle connections
		MaxIdleConnsPerHost: 10,               // Maximum idle connections per host
		IdleConnTimeout:     90 * time.Second, // Idle connection timeout
		DisableCompression:  false,            // Enable compression
		DisableKeepAlives:   false,            // Enable keep-alives
	}

	return &http.Client{
		Transport: transport,
		Timeout:   60 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Custom redirect handling
			if len(via) >= 10 {
				return errors.New("stopped after 10 redirects")
			}
			return nil
		},
	}
}
